// Package lattice is for computing lattice things like meet and join
package lattice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/gviz"
)

type Node struct {
	ID          string
	Group       string // e.g., like a chain identity
	Children    Nodeset
	Descendants Nodeset
}

type Nodeset map[string]struct{}

func (n Nodeset) Sorted() (out []string) {
	for k := range n {
		out = append(out, k)
	}
	sort.Strings(out)
	return
}

func (n Nodeset) Add(id string) {
	n[id] = struct{}{}
}

func (n Nodeset) Has(id string) bool {
	_, ok := n[id]
	return ok
}

func (n Nodeset) Remove(id string) {
	delete(n, id)
}

func (n Nodeset) Merge(o Nodeset) {
	for k := range o {
		n[k] = struct{}{}
	}
}

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
	names := make(map[string]string)
	chain.BreadthFirst(a, func(id string) {
		names[id] = fmt.Sprintf("%d", len(names))
	})
	meet := "" //chain.Meet(a, b)
	if err := chain.ToGraphViz("g.svg", map[string]string{
		a:    "yellow",
		b:    "yellow",
		meet: "red",
	}, names); err != nil {
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

func (l Lattice) BreadthFirst(root string, f func(string)) {
	q := queue{}
	discovered := make(map[string]bool)
	node := func(id string) *Node {
		return l.m[id]
	}
	enqueue := func(id string) {
		if discovered[id] {
			return
		}
		discovered[id] = true
		q.enqueue(id)
		fmt.Printf("%q; %q\n", id, q.slice)
		f(id)
	}
	dequeue := func() string {
		return q.dequeue()
	}
	enqueue(root)
	for !q.empty() {
		for _, c := range node(dequeue()).Children.Sorted() {
			enqueue(c)
		}
	}
}

type queue struct {
	slice []string
}

func (q *queue) empty() bool {
	return len(q.slice) == 0
}

func (q *queue) enqueue(x string) {
	q.slice = append(q.slice, x)
}

func (q *queue) dequeue() string {
	x := q.slice[0]
	q.slice = q.slice[1:]
	return x
}

// perhaps open up in a browser, highlighting specific nodes with colors
func (l Lattice) ToGraphViz(svg string, names, colors map[string]string) error {
	buf, err := gviz.Compile(l, names, colors)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("g.gv", buf, os.ModePerm); err != nil {
		return err
	}
	return gviz.Dot("g.gv", svg)
}

func (n Node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

// a graph of nodes where every two has a unique meet (semi-lattice)
type Lattice struct {
	m map[string]*Node
}

func (l Lattice) Nodes() (out []string) {
	for k := range l.m {
		out = append(out, k)
	}
	sort.Strings(out)
	return
}

type edge struct {
	from, to string
}

func (e edge) From() string {
	return e.from
}
func (e edge) To() string {
	return e.to
}

func (l Lattice) Edges() (out []gviz.Edge) {
	for from, parent := range l.m {
		for to := range parent.Children {
			out = append(out, edge{from: from, to: to})
		}
	}
	return
}

func (l Lattice) Children(a string) map[string]*Node {
	out := make(map[string]*Node)
	node := func(id string) *Node {
		return l.m[id]
	}
	for _, c := range node(a).Children.Sorted() {
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
		panic("can't sort nodes")
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
		m: make(map[string]*Node),
	}
	add := func(n *Node) {
		out.m[n.ID] = n
	}
	newNode := func(name string) *Node {
		if n, ok := out.m[name]; ok {
			return n
		}
		return &Node{
			ID:          name,
			Children:    make(Nodeset),
			Descendants: make(Nodeset),
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
				for _, c := range children {
					merge.Children.Add(c)
				}
			} else {
				merge = out.m[children[0]]
			}
			add(merge)

			chain := newNode(fmt.Sprintf("%d.%d", i, j))
			chain.Children.Add(merge.ID)
			add(chain)

			m[i] = chain.ID

		}
		last = m
	}
	return out
}
