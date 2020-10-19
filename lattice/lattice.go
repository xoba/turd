// Package lattice is for computing lattice things like meet and join
package lattice

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/cnfg"
)

// TODO: need a random seed for reproducibility
// TODO: also test cases of verified meets
func Run(c cnfg.Config) error {
	var seed int
	if c.Seed == 0 {
		seed = int(rand.Uint64())
		if seed < 0 {
			seed = -seed
		}
		fmt.Printf("seed = %d\n", seed)
	} else {
		seed = c.Seed
	}
	chain := Generate(rand.New(rand.NewSource(int64(seed))), 3, 5)
	a, b := "1.4", "2.4"
	meet := chain.Meet(a, b)
	if err := chain.ToGraphViz("g.svg", map[string]string{
		a:    "yellow",
		b:    "yellow",
		meet: "red",
	}); err != nil {
		return err
	}
	f, err := os.Create("g.html")
	if err != nil {
		return err
	}
	fmt.Fprintf(f, `<!DOCTYPE html>
<img src='g.svg'>
`)
	f.Close()
	return open.Run("g.html")
}

// perhaps open up in a browser, highlighting specific nodes with colors
func (l Lattice) ToGraphViz(svg string, colors map[string]string) error {
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
	var ids []string
	for k := range l.Nodes {
		ids = append(ids, k)
	}
	for _, i := range ids {
		c := colors[i]
		if c == "" {
			c = "white"
		}
		fmt.Fprintf(f, "%s [ label=%q; fillcolor=%s style=filled ];\n", id(i), i, c)
	}
	for _, i := range ids {
		n := l.Nodes[i]
		for _, c := range n.Children {
			fmt.Fprintf(f, "%s -> %s [ label=%q ];\n", id(n.ID), id(c), "")
		}
	}
	fmt.Fprintf(f, "}\n")
	f.Close()
	return graphviz("dot", "g.gv", svg, "svg")
}

type Node struct {
	ID       string
	Time     time.Time // to help order, assuming timestamps approx. correct
	Children []string
}

func (n Node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
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
	var ids []*Node
	for _, x := range list {
		ids = append(ids, x)
	}
	fmt.Printf("intersection:\n")
	for _, n := range ids {
		fmt.Println(n)
	}
	return list[0].ID
}

func Generate(r *rand.Rand, chains, length int) Lattice {
	out := Lattice{
		Nodes: make(map[string]*Node),
	}
	add := func(n *Node) {
		out.Nodes[n.ID] = n
	}
	newNode := func(name string) *Node {
		if n, ok := out.Nodes[name]; ok {
			return n
		}
		time.Sleep(10 * time.Millisecond)
		return &Node{
			ID:   name,
			Time: time.Now().UTC(),
		}
	}
	genesis := newNode("genesis")
	add(genesis)
	var last map[int]string
	for j := 0; j < length; j++ {
		m := make(map[int]string)
		for i := 0; i < chains; i++ {

			// TODO: if only one child, skip the merge node
			var children []string
			if last == nil {
				children = append(children, genesis.ID)
			} else {
				children = append(children, last[i])
				if r := r.Intn(3); r != i {
					children = append(children, last[r])
				}
			}
			sort.Strings(children)

			var merge *Node
			if len(children) > 1 {
				merge = newNode("[" + strings.Join(children, ",") + "]")
				merge.Children = children
			} else {
				merge = out.Nodes[children[0]]
			}
			add(merge)

			chain := newNode(fmt.Sprintf("%d.%d", i, j))
			chain.Children = append(chain.Children, merge.ID)
			add(chain)

			m[i] = chain.ID

		}
		last = m
	}
	return out
}

func graphviz(graphvizCommand, gv, out, format string) error {
	if err := exec.Command(graphvizCommand, "-v", "-o", out, fmt.Sprintf("-T%s", format), gv).Run(); err != nil {
		return fmt.Errorf("can't run graphviz (%s): %v", graphvizCommand, err)
	}
	return nil
}
