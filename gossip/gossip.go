package gossip

import (
	"fmt"
	"log"
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
	connections map[string]tnet.Conn
	records     map[string]NodeRecord // map by address
}

func NewNode(port int) (*Node, error) {
	node := Node{
		lock: new(sync.Mutex),
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
	Type    string
	Records []NodeRecord
}

func (node *Node) Gossip(seeds ...string) error {
	go func() {
		time.Sleep(time.Second)
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
		go handleConn(c)
	}
}

func (node *Node) engageWithNodes(seeds ...string) error {
	return nil
}

func handleConn(c tnet.Conn) {
	fmt.Printf("handling %s\n", c.Remote())
}
