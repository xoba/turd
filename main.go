package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/xoba/turd/taws"
	"github.com/xoba/turd/tnet"
	"golang.org/x/net/websocket"
)

type Config struct {
	Mode          string
	AWSProfile    string
	Port          int
	PublicKeyFile string
}

const SeedNode = "http://localhost:8080"

func main() {
	var c Config
	flag.StringVar(&c.Mode, "m", "node", "mode to run")
	flag.StringVar(&c.PublicKeyFile, "key", "pub.dat", "public key file")
	flag.StringVar(&c.AWSProfile, "aws", "", "aws profile")
	flag.IntVar(&c.Port, "p", 8080, "http port to run on")
	flag.Parse()
	if err := c.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (c Config) Run() error {
	modes := map[string]func() error{
		"node":    c.RunNode,
		"launch":  c.LaunchNode,
		"listen":  c.Listen,
		"connect": c.Connect,
	}
	handler, ok := modes[c.Mode]
	if !ok {
		var list []string
		for k := range modes {
			list = append(list, k)
		}
		sort.Strings(list)
		return fmt.Errorf("unrecognized mode %q; should be one of %q",
			c.Mode,
			strings.Join(list, ", "),
		)
	}
	return handler()
}

// connect to a network listener
func (config Config) Connect() error {
	n, err := tnet.NewNetwork(nil, 8081)
	if err != nil {
		return err
	}
	var p *tnet.PublicKey
	if f := config.PublicKeyFile; f != "" {
		var key tnet.PublicKey
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		if err := key.UnmarshalBinary(buf); err != nil {
			return err
		}
		p = &key
	}
	c, err := n.Dial(tnet.Node{Address: "localhost:8080", PublicKey: p})
	if err != nil {
		return err
	}
	fmt.Printf("remote: %v\n", c.Remote())
	for {
		buf, err := c.Receive()
		if err != nil {
			return err
		}
		fmt.Printf("received %q\n", string(buf))
		if err := c.Send([]byte(fmt.Sprintf("got %q", string(buf)))); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

// play with network listeners
func (config Config) Listen() error {
	key, err := tnet.NewKey()
	if err != nil {
		return err
	}
	if f := config.PublicKeyFile; f != "" {
		buf, err := key.Public().MarshalBinary()
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(f, buf, os.ModePerm); err != nil {
			return err
		}
	}
	n, err := tnet.NewNetwork(key, 8080)
	if err != nil {
		return err
	}
	ln, err := n.Listen()
	if err != nil {
		return err
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("remote: %v\n", c.Remote())
		go func() {
			if err := handleConnection(c); err != nil {
				log.Printf("oops: %v", err)
			}
		}()
	}
	return nil
}

func handleConnection(c tnet.Conn) error {
	var i int
	for {
		i++
		if err := c.Send([]byte(fmt.Sprintf("packet %d", i))); err != nil {
			return err
		}
		buf, err := c.Receive()
		if err != nil {
			return err
		}
		fmt.Printf("received %q\n", string(buf))
	}
}

type NodeID struct {
	RemoteAddr string // verified "network" address
	PublicKey  []byte
}

type LiveNode struct {
	NodeID
	LastSeen time.Time
}

type Handler struct {
	lock  sync.Locker
	nodes map[string]*LiveNode
}

func NewHandler() *Handler {
	return &Handler{
		lock:  new(sync.Mutex),
		nodes: make(map[string]*LiveNode),
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("%s %q", r.Method, r.RequestURI)
	switch r.URL.Path {
	case "/":
	case "/t":
		websocket.Handler(h.ServeWebsocket).ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

type NodeMessage struct {
	Type    string
	Payload []byte
}

func (h Handler) ServeWebsocket(ws *websocket.Conn) {
	for {
		var m NodeMessage
		if err := websocket.JSON.Receive(ws, &m); err != nil {
			log.Printf("error serving %s: %v", ws.RemoteAddr(), err)
			break
		}
		switch m.Type {
		case "register":

		}
	}
}

func (c Config) RunNode() error {
	h := NewHandler()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: h,
	}
	return s.ListenAndServe()
}

func (c Config) LaunchNode() error {
	s, err := taws.NewSessionFromProfile(c.AWSProfile)
	if err != nil {
		return err
	}
	r, err := ec2.New(s).DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return err
	}
	fmt.Println(r)
	return fmt.Errorf("LaunchNode unimplemented")
}
