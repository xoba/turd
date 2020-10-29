// Package poset is for computing poset things like meet and join
// (even though posets don't necessarily have unique meets etc.)
package poset

import (
	"container/list"
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

// a graph of nodes where every two has a unique meet (semi-lattice)
type Poset struct {
	m map[string]*Node
}

type Node struct {
	id, label, group, shape string
	Type                    string
	Children                Nodeset
	Descendents             Nodeset
}

func (n Node) ID() string {
	return n.id
}
func (n Node) Label() string {
	return n.label
}
func (n Node) Group() string {
	return n.group
}
func (n Node) Shape() string {
	switch n.Type {
	case "top", "bottom":
		return "Mcircle"
	case "merge":
		return "ellipse"
	case "block":
		return "box"
	default:
		panic("illegal: " + n.Type)
	}
	return n.shape
}

type Nodeset map[string]struct{}

func (n Nodeset) Sorted() (out []string) {
	for k := range n {
		out = append(out, k)
	}
	sort.Sort(sort.StringSlice(out))
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

// BUG: go run . -m poset -s 8896833665541862532 ... 1.2 should be meet, not 0.2

// TODO: also test cases of verified meets
// for instance: go run . -m lattice -s 617624903177646721
// is meet non-unique among 1.2 and 2.2??? reverse sorting changes answer, a bad sign!
// i think the answer is that choosing either 1.2 or 2.2 as meet brings together the same data;
// so in some sense meet should be union of 1.2 and 2.2? the clincher is that 1.2 and 2.2 are not
// ordered with respect to one another.
// also check that meet is idempotent, commutative, and associative
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
	const (
		chains = 5
		length = 7
	)
	chain := Generate(rand.New(rand.NewSource(int64(seed))), chains, length)
	a, b := fmt.Sprintf("%d.%d", 0, length-1), fmt.Sprintf("%d.%d", 1, length-1)
	meet := chain.Meet(a, b)

	names := make(map[string]string)

	display := func() error {
		if err := chain.ToGraphViz("g.svg", map[string]string{
			a:      "yellow",
			b:      "yellow",
			"join": "cyan",
			meet:   "red",
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

	if false {
		BreadthFirstSearch(chain, a, func(id string) bool {
			names[id] = fmt.Sprintf("%d", len(names))
			return false
		})
		return display()
	}

	join := NewNode("join", "", "", "top")
	chain.m[join.id] = join
	join.Children.Add(a)
	join.Children.Add(b)
	join.Descendents = chain.CalcDescendents(join.id)

	return display()
}

func (p Poset) Children(id string) []string {
	return p.m[id].Children.Sorted()
}

type HasChildren interface {
	// GetChildren returns the sorted children of given node
	Children(string) []string
}

func BreadthFirstSearch(p HasChildren, root string, stop func(string) bool) string {
	q := newQueue()
	enqueue := func(id string) bool {
		if !q.enqueue(id) {
			return false
		}
		return stop(id)
	}

	if enqueue(root) {
		return root
	}
	for !q.empty() {
		for _, c := range p.Children(q.dequeue()) {
			if enqueue(c) {
				return c
			}
		}
	}
	return ""
}

type queue struct {
	list       *list.List
	discovered map[string]bool
}

func newQueue() *queue {
	return &queue{
		list:       list.New(),
		discovered: make(map[string]bool),
	}
}

func (q *queue) empty() bool {
	return q.list.Len() == 0
}

// enqueue if it hasn't been enqueued before, returns whether queued
func (q *queue) enqueue(id string) bool {
	if q.discovered[id] {
		return false
	}
	q.discovered[id] = true
	q.list.PushBack(id)
	return true
}

func (q *queue) dequeue() string {
	e := q.list.Front()
	q.list.Remove(e)
	return e.Value.(string)
}

// perhaps open up in a browser, highlighting specific nodes with colors
func (l Poset) ToGraphViz(svg string, colors map[string]string) error {
	buf, err := gviz.Compile(l, svg, colors)
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

func (l Poset) Nodes() (out []gviz.Node) {
	for _, n := range l.m {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID() < out[j].ID()
	})
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

func (l Poset) Edges() (out []gviz.Edge) {
	for from, parent := range l.m {
		for to := range parent.Children {
			out = append(out, edge{from: from, to: to})
		}
	}
	return
}

// returns meet of two nodes, if any, or one if many
func (l Poset) Meet(a, b string) string {
	desc := l.m[b].Descendents
	return BreadthFirstSearch(l, a, func(id string) bool {
		if desc.Has(id) {
			return true
		}
		return false
	})
}

func NewNode(id, label, group, t string) *Node {
	if id == "" {
		panic("needs id")
	}
	if label == "" {
		label = id
	}
	return &Node{
		id:          id,
		label:       label,
		group:       group,
		Type:        t,
		Children:    make(Nodeset),
		Descendents: make(Nodeset),
	}
}

func Generate(r *rand.Rand, chains, length int) Poset {
	out := Poset{
		m: make(map[string]*Node),
	}
	add := func(n *Node) {
		out.m[n.id] = n
	}
	newNode := func(id, label, group, t string) *Node {
		if n, ok := out.m[id]; ok {
			return n
		}
		return NewNode(id, label, group, t)
	}
	genesis := newNode("g", "genesis", "", "bottom")
	add(genesis)
	var last map[int]string
	for j := 0; j < length; j++ {
		m := make(map[int]string)
		for i := 0; i < chains; i++ {

			// TODO: if only one child, skip the merge node
			var children []string
			if last == nil {
				children = append(children, genesis.id)
			} else {
				children = append(children, last[i])
				if k := r.Intn(chains); k != i {
					children = append(children, last[k])
				}
			}
			sort.Strings(children)

			var merge *Node
			if len(children) > 1 {
				id := "[" + strings.Join(children, ",") + "]"
				merge = newNode(id, "m", "", "merge")
				for _, c := range children {
					merge.Children.Add(c)
				}
			} else {
				merge = out.m[children[0]]
			}
			add(merge)
			chain := newNode(
				fmt.Sprintf("%d.%d", i, j), // id
				fmt.Sprintf("%d", j),       // label
				fmt.Sprintf("%c", 'a'+i),   // group
				"block",
			)
			chain.Children.Add(merge.id)
			add(chain)

			m[i] = chain.id

		}
		last = m
	}
	for _, n := range out.m {
		n.Descendents = out.CalcDescendents(n.id)
	}
	return out
}

func (l Poset) CalcDescendents(id string) Nodeset {
	out := make(Nodeset)
	for c := range l.m[id].Children {
		out.Add(c)
		out.Merge(l.CalcDescendents(c))
	}
	return out
}
