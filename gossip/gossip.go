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

func Run(c cnfg.Config) error {
	const seedPort = 8080
	seedAddr := fmt.Sprintf("localhost:%d", seedPort)

	run := func(port int, seeds ...string) error {
		n, err := NewNode(port)
		if err != nil {
			return err
		}
		go func() {
			if err := n.Gossip(); err != nil {
				log.Printf("node on port %d failed: %v", port, err)
			}
		}()
		return nil
	}

	if err := run(seedPort); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		if err := run(i+1+seedPort, seedAddr); err != nil {
			log.Fatal(err)
		}
	}

	<-make(chan bool)

	return nil
}

type NodeRecord struct {
	LastSeen time.Time
	tnet.Node
}

type Node struct {
	sync.Locker
	tnet.Network
	closed      bool
	key         *tnet.PrivateKey
	seeds       []string
	seen        map[string]bool       // message id's already seen
	connections map[string]tnet.Conn  // map by hash of public key
	records     map[string]NodeRecord // map by public key hash?
}

func (n *Node) Closed() bool {
	n.Lock()
	defer n.Unlock()
	return n.closed
}

func (n *Node) Close() error {
	n.Lock()
	defer n.Unlock()
	n.closed = true
	return nil
}

func NewNode(port int, seeds ...string) (*Node, error) {
	key, err := tnet.NewKey()
	if err != nil {
		return nil, err
	}
	n, err := tnet.NewTCPLocalhost(port)
	if err != nil {
		return nil, err
	}
	return &Node{
		Locker:      new(sync.Mutex),
		Network:     n,
		key:         key,
		seeds:       seeds,
		seen:        make(map[string]bool),
		connections: make(map[string]tnet.Conn),
		records:     make(map[string]NodeRecord),
	}, nil
}

type Message struct {
	ID      string
	Type    string
	Records []NodeRecord `json:",omitempty"`
}

func (m Message) String() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

func (node *Node) Gossip() error {
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

func (node *Node) process(m *Message) {
	node.Lock()
	defer node.Unlock()
	if node.seen[m.ID] {
		return
	}
	node.seen[m.ID] = true
	for _, r := range m.Records {
		// don't check last times seen yet
		node.records[r.Address] = r
	}
}

func (node *Node) makeNewConnections() error {
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
			nodes = append(nodes, tnet.Node{Address: s})
		}
	}
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	for _, n := range nodes {
		if key := n.PublicKey; key != nil {
			if node.alreadyConnected(key) {
				continue
			}
		}
		fmt.Printf("going to connect with %s\n", n)
		c, err := node.Dial(node.key, n)
		if err != nil {
			return err
		}
		go node.handleConn(c)
	}
	return nil
}

func (node *Node) engageWithNodes() {
	for {
		if node.Closed() {
			return
		}
		time.Sleep(time.Second)
		fmt.Printf("%s engaging\n", node.key.Public())
		if err := node.makeNewConnections(); err != nil {
			log.Printf("error making new connections: %v", err)
		}
	}
}

func (node *Node) addConn(c tnet.Conn) {
	node.Lock()
	defer node.Unlock()
	node.connections[c.Remote().PublicKey.String()] = c
}

func (node *Node) removeConn(c tnet.Conn) {
	node.Lock()
	defer node.Unlock()
	defer c.Close()
	delete(node.connections, c.Remote().PublicKey.String())
}

// TODO: need to think more about structure here....
func (node *Node) handleConn(c tnet.Conn) {
	fmt.Printf("handling %s\n", c.Remote())
	defer c.Close()
	if h := c.Remote().PublicKey; node.alreadyConnected(h) {
		fmt.Printf("already connected to %s\n", h)
		return
	}
	node.addConn(c)
	go func() {
		defer node.removeConn(c)
		for {
			m, err := receive(c)
			if err != nil {
				log.Printf("problem with %v: %v", c.Remote(), err)
				return
			}
			fmt.Printf("got %s from %s\n", m, c.Remote())
			node.process(m)
		}
	}()
}

func (node *Node) alreadyConnected(remote *tnet.PublicKey) bool {
	node.Lock()
	defer node.Unlock()
	_, ok := node.connections[remote.String()]
	return ok
}

func send(c tnet.Conn, m Message) error {
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
