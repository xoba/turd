package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/xoba/turd/taws"
	"golang.org/x/net/websocket"
)

type Config struct {
	Mode       string
	AWSProfile string
	Port       int
}

func main() {
	var c Config
	flag.StringVar(&c.Mode, "m", "node", "mode to run")
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
		"node":   c.RunNode,
		"launch": c.LaunchNode,
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

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
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

func (h Handler) ServeWebsocket(ws *websocket.Conn) {
	for {

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
