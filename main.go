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
	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/gossip"
	"github.com/xoba/turd/infl"
	"github.com/xoba/turd/taws"
	"github.com/xoba/turd/tnet"
	"golang.org/x/net/websocket"
)

func main() {
	var c cnfg.Config
	flag.StringVar(&c.Mode, "m", "node", "mode to run")
	flag.StringVar(&c.PublicKeyFile, "pub", "pub.dat", "public key file")
	flag.StringVar(&c.PrivateKeyFile, "priv", "priv.dat", "private key file")
	flag.StringVar(&c.AWSProfile, "aws", "", "aws profile")
	flag.IntVar(&c.Port, "p", 8080, "http port to run on")
	flag.Parse()
	if err := Run(c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(c cnfg.Config) error {
	modes := map[string]func(cnfg.Config) error{
		"keys":      tnet.SharedKey,
		"gossip":    gossip.Run,
		"hnode":     RunHTMLNode,
		"launch":    LaunchNode,
		"listen":    Listen,
		"connect":   Connect,
		"inflation": infl.Run,
		"infbig":    infl.RunBig,
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
	return handler(c)
}

// connect to a network listener
func Connect(config cnfg.Config) error {
	n, err := tnet.NewTCPLocalhost(8081)
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
	key, err := tnet.NewKey()
	if err != nil {
		return err
	}
	c, err := n.Dial(key, tnet.Node{Address: "localhost:8080", PublicKey: p})
	if err != nil {
		return err
	}
	defer c.Close()
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
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func CachedKey(f string) (*tnet.PrivateKey, error) {
	if f == "" {
		return tnet.NewKey()
	}
	if _, err := os.Stat(f); err != nil {
		key, err := tnet.NewKey()
		if err != nil {
			return nil, err
		}
		buf, err := key.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(f, buf, os.ModePerm); err != nil {
			return nil, err
		}
	}
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var key tnet.PrivateKey
	if err := key.UnmarshalBinary(buf); err != nil {
		return nil, err
	}
	return &key, nil
}

// play with network listeners
func Listen(config cnfg.Config) error {
	key, err := CachedKey(config.PrivateKeyFile)
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

	n, err := tnet.NewTCPLocalhost(8080)
	if err != nil {
		return err
	}
	ln, err := n.Listen()
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		c, err := ln.Accept(key)
		if err != nil {
			log.Printf("oops: %v", err)
			continue
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
	defer c.Close()
	var i int
	for {
		i++
		if err := c.Send([]byte(fmt.Sprintf("packet %d at %v", i, time.Now()))); err != nil {
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

func RunHTMLNode(c cnfg.Config) error {
	h := NewHandler()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: h,
	}
	return s.ListenAndServe()
}

func LaunchNode(c cnfg.Config) error {
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
