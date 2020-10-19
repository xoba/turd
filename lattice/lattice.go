// Package lattice is for computing lattice things like meet and join
package lattice

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/xoba/turd/cnfg"
)

func Run(cnfg.Config) error {
	chain := Generate(3, 9)
	a, b := "2.8", "1.8"
	meet := chain.Meet(a, b)
	return chain.ToGraphViz(map[string]string{
		a:    "yellow",
		b:    "yellow",
		meet: "red",
	})
}

// perhaps open up in a browser, highlighting specific nodes with colors
func (l Lattice) ToGraphViz(colors map[string]string) error {
	id := func(name string) string {
		h := md5.New()
		h.Write([]byte(name))
		return fmt.Sprintf("N%x", h.Sum(nil))[:8]
	}
	f, err := os.Create("g.gv")
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "digraph {\n")
	for _, n := range l.Nodes {
		c := colors[n.ID]
		if c == "" {
			c = "white"
		}
		fmt.Fprintf(f, "%s [ label=%q; fillcolor=%s style=filled ];\n", id(n.ID), n.ID, c)
	}
	for _, n := range l.Nodes {
		for _, c := range n.Children {
			fmt.Fprintf(f, "%s -> %s [ label=%q ];\n", id(n.ID), id(c), "")
		}
	}
	fmt.Fprintf(f, "}\n")
	f.Close()
	return graphviz("dot", "g.gv", "g.svg", "svg")
}

type Node struct {
	ID       string
	Children []string
	Time     time.Time // to help order, assuming timestamps approx. correct
}

// a graph of nodes where every two has a unique meet (semi-lattice)
type Lattice struct {
	Nodes map[string]*Node
}

func (l Lattice) Children(a string) map[string]*Node {
	out := make(map[string]*Node)
	node := func(id string) *Node {
		return l.Nodes[id]
	}
	for _, c := range node(a).Children {
		out[c] = node(c)
		for k, v := range l.Children(c) {
			out[k] = v
		}
	}
	return out
}

// returns meet of two nodes, if any
func (l Lattice) Meet(a, b string) string {
	intersection := make(map[string]*Node)
	ac := l.Children(a)
	for k := range l.Children(b) {
		if v, ok := ac[k]; ok {
			intersection[k] = v
		}
	}
	var list []*Node
	for _, n := range intersection {
		list = append(list, n)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Time.After(list[j].Time)
	})
	return list[0].ID
}

func Generate(chains, length int) Lattice {
	out := Lattice{
		Nodes: make(map[string]*Node),
	}
	add := func(n *Node) {
		out.Nodes[n.ID] = n
		time.Sleep(time.Millisecond)
	}
	newNode := func(name string) *Node {
		return &Node{
			ID:   name,
			Time: time.Now(),
		}
	}
	genesis := newNode("g")
	add(genesis)
	var lastm map[int]string
	for j := 0; j < length; j++ {
		m := make(map[int]string)
		for i := 0; i < chains; i++ {
			chain := newNode(fmt.Sprintf("%d.%d", i, j))
			if lastm == nil {
				chain.Children = append(chain.Children, genesis.ID)
			} else {
				chain.Children = append(chain.Children, lastm[i])
				if r := rand.Intn(3); r != i {
					chain.Children = append(chain.Children, lastm[r])
				}
			}
			m[i] = chain.ID
			add(chain)
		}
		lastm = m
	}
	return out
}

func graphviz(graphvizCommand, gv, out, format string) error {
	if err := exec.Command(graphvizCommand, "-v", "-o", out, fmt.Sprintf("-T%s", format), gv).Run(); err != nil {
		return fmt.Errorf("can't run graphviz (%s): %v", graphvizCommand, err)
	}
	return nil
}
