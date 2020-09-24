package gossip

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
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
			if err := n.Gossip(seeds...); err != nil {
				log.Fatal(err)
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
	lock        sync.Locker
	n           tnet.Network
	key         *tnet.PrivateKey
	seen        map[string]bool // message id's already seen
	connections map[string]tnet.Conn
	records     map[string]NodeRecord // map by address, or rather public key hash?
}

func NewNode(port int) (*Node, error) {

	node := Node{
		lock:    new(sync.Mutex),
		seen:    make(map[string]bool),
		records: make(map[string]NodeRecord),
	}
	key, err := tnet.NewKey()
	if err != nil {
		return nil, err
	}
	node.key = key
	n, err := tnet.NewTCPLocalhost(port)
	if err != nil {
		return nil, err
	}
	node.n = n
	return &node, nil
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

func (node *Node) Gossip(seeds ...string) error {
	go func() {
		if err := node.engageWithNodes(seeds...); err != nil {
			log.Fatal(err)
		}
	}()
	ln, err := node.n.Listen()
	if err != nil {
		return err
	}
	for {
		c, err := ln.Accept(node.key)
		if err != nil {
			return err
		}
		go node.handleConn(c)
	}
}

func (node *Node) process(m *Message) {
	node.lock.Lock()
	defer node.lock.Unlock()
	if node.seen[m.ID] {
		return
	}
	node.seen[m.ID] = true
	for _, r := range m.Records {
		// don't check last times seen yet
		node.records[r.Address] = r
	}
}

func (node *Node) engageWithNodes(seeds ...string) error {
	if len(seeds) == 0 {
		return nil
	}
	time.Sleep(time.Second)
	c, err := node.n.Dial(node.key, tnet.Node{Address: seeds[0]})
	if err != nil {
		return err
	}
	m := Message{
		ID:   uuid.New().String(),
		Type: "nodes",
	}
	if err := send(c, m); err != nil {
		return err
	}
	if err := send(c, m); err != nil {
		return err
	}
	return nil
}

func (node *Node) handleConn(c tnet.Conn) {
	fmt.Printf("handling %s\n", c.Remote())
	if addr := c.Remote().Address; node.alreadyConnected(addr) {
		fmt.Printf("already connected to %s\n", addr)
		return
	}
	go func() {
		for {
			m, err := receive(c)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("got %s from %s\n", m, c.Remote())
			node.process(m)
		}
	}()
}

func (node *Node) alreadyConnected(remote string) bool {
	node.lock.Lock()
	defer node.lock.Unlock()
	_, ok := node.connections[remote]
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
