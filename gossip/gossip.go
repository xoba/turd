package gossip

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/tnet"
)

const Seed = "localhost:8080"

func Run(c cnfg.Config) error {
	n, err := NewNode(c.Port, Seed)
	if err != nil {
		return err
	}
	return n.Gossip()
}

type serverRecord struct {
	LastSeen time.Time
	tnet.Node
}

type server struct {
	sync.Locker
	tnet.Network

	closed      bool
	key         *tnet.PrivateKey
	seeds       []string
	seen        map[string]bool          // message id's already seen
	connections map[string]connRecord    // map by key string
	records     map[string]*serverRecord // map by key string
}

type Message struct {
	ID      string
	Type    string
	Records []*serverRecord `json:",omitempty"`
}

type connRecord struct {
	tnet.Conn
	send chan *Message
}

func (n *server) Closed() bool {
	n.Lock()
	defer n.Unlock()
	return n.closed
}

func (n *server) Close() error {
	n.Lock()
	defer n.Unlock()
	n.closed = true
	return nil
}

func NewNode(port int, seeds ...string) (*server, error) {
	key, err := tnet.NewKey()
	if err != nil {
		return nil, err
	}
	n, err := tnet.NewTCPLocalhost(port)
	if err != nil {
		return nil, err
	}
	return &server{
		Locker:      new(sync.Mutex),
		Network:     n,
		key:         key,
		seeds:       seeds,
		seen:        make(map[string]bool),
		connections: make(map[string]connRecord),
		records:     make(map[string]*serverRecord),
	}, nil
}

func (m Message) String() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

func (node *server) Gossip() error {
	go node.engageWithNodes()
	ln, err := node.Listen()
	if err != nil {
		return err
	}
	for {
		if node.Closed() {
			return nil
		}
		c, err := ln.Accept(node.key)
		if err != nil {
			log.Printf("error accepting: %v", err)
			continue
		}
		go node.handleConn(c)
	}
}

func (node *server) process(m *Message) error {
	log.Printf("process(%v)", m)
	node.Lock()
	defer node.Unlock()
	if node.seen[m.ID] {
		return nil
	}
	node.seen[m.ID] = true
	for _, r := range m.Records {
		key := r.PublicKey.String()
		r0, ok := node.records[key]
		if ok {
			if r.LastSeen.After(r0.LastSeen) {
				node.records[key].LastSeen = r.LastSeen
			}
		} else {
			node.records[key] = r
		}
	}
	var list []connRecord
	func() {
		node.Lock()
		defer node.Unlock()
		for _, v := range node.connections {
			list = append(list, v)
		}
	}()
	for _, s := range list {
		s.send <- m
	}
	return nil
}

func (node *server) makeNewConnections() error {
	var nodes []tnet.Node
	func() {
		node.Lock()
		defer node.Unlock()
		for k, v := range node.records {
			if _, ok := node.connections[k]; !ok {
				nodes = append(nodes, v.Node)
			}
		}
	}()
	if len(nodes) == 0 {
		for _, s := range node.seeds {
			if s != node.Addr() {
				nodes = append(nodes, tnet.Node{Address: s})
			}
		}
	}
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	var bad []tnet.Node
	for _, n := range nodes {
		if key := n.PublicKey; key != nil {
			if node.alreadyConnected(key) {
				continue
			}
		}
		c, err := node.Dial(node.key, n)
		if err != nil {
			bad = append(bad, n)
			continue
		}
		go node.handleConn(c)
	}
	func() {
		node.Lock()
		defer node.Unlock()
		for _, n := range bad {
			fmt.Println(n) // TODO: public key is nil for seed!
			delete(node.records, n.PublicKey.String())
		}
	}()

	return nil
}

func (node *server) engageWithNodes() {
	for {
		if node.Closed() {
			return
		}
		time.Sleep(time.Second)
		if err := node.makeNewConnections(); err != nil {
			log.Printf("error making new connections: %v", err)
		}
	}
}

func (node *server) addConn(c connRecord) {
	node.Lock()
	defer node.Unlock()
	key := c.Remote().PublicKey.String()
	node.connections[key] = c
	node.records[key] = &serverRecord{
		LastSeen: time.Now(),
		Node:     c.Remote(),
	}
}

func (node *server) removeConn(c connRecord) {
	node.Lock()
	defer node.Unlock()
	defer c.Close()
	delete(node.connections, c.Remote().PublicKey.String())
}

func (node *server) handleConn(c tnet.Conn) error {
	if h := c.Remote().PublicKey; node.alreadyConnected(h) {
		return nil
	}
	log.Printf("handling %s", c.Remote())
	cr := connRecord{
		Conn: c,
		send: make(chan *Message),
	}
	defer close(cr.send)
	node.addConn(cr)
	defer node.removeConn(cr)
	go func() {
		defer node.removeConn(cr)
		for {
			m, err := receive(c)
			if err != nil {
				log.Printf("can't receive from %v: %v", c.Remote(), err)
				return
			}
			if err := node.process(m); err != nil {
				log.Printf("can't send to %v: %v", c.Remote(), err)
				return
			}
		}
	}()
	for m := range cr.send {
		send(cr.Conn, m)
	}
	return nil
}

func (node *server) alreadyConnected(remote *tnet.PublicKey) bool {
	node.Lock()
	defer node.Unlock()
	_, ok := node.connections[remote.String()]
	return ok
}

func send(c tnet.Conn, m *Message) error {
	fmt.Printf("send %s\n", m)
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return c.Send(buf)
}

func receive(c tnet.Conn) (*Message, error) {
	buf, err := c.Receive()
	if err != nil {
		return nil, err
	}
	var m Message
	if err := json.Unmarshal(buf, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
