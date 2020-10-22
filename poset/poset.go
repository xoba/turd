// Package poset is for computing poset things like meet and join
// (even though posets don't necessarily have unique meets etc.)
package poset

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

// a graph of nodes where every two has a unique meet (semi-lattice)
type Poset struct {
	m map[string]*Node
}

type Node struct {
	ID          string
	Group       string // e.g., like a chain identity
	Children    Nodeset
	Descendents Nodeset
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
	chain := Generate(rand.New(rand.NewSource(int64(seed))), 3, 5)
	a, b := "1.4", "2.4"

	names := make(map[string]string)
	chain.BreadthFirstSearch(a, func(id string) bool {
		if false {
			names[id] = fmt.Sprintf("%d", len(names))
		}
		return false
	})

	join := NewNode("join")
	chain.m[join.ID] = join
	join.Children.Add(a)
	join.Children.Add(b)
	join.Descendents = chain.CalcDescendents(join.ID)

	meet := chain.Meet(a, b)
	if err := chain.ToGraphViz("g.svg", map[string]string{
		a:      "yellow",
		b:      "yellow",
		"join": "cyan",
		meet:   "red",
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

// returns id of node that passed function returns true on
func (l Poset) BreadthFirstSearch(root string, f func(string) bool) string {
	q := queue{}
	discovered := make(map[string]bool)
	node := func(id string) *Node {
		return l.m[id]
	}
	enqueue := func(id string) bool {
		if discovered[id] {
			return true
		}
		discovered[id] = true
		q.enqueue(id)
		return f(id)
	}
	dequeue := func() string {
		return q.dequeue()
	}
	if enqueue(root) {
		return root
	}
	for !q.empty() {
		for _, c := range node(dequeue()).Children.Sorted() {
			if enqueue(c) {
				return c
			}
		}
	}
	return ""
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
func (l Poset) ToGraphViz(svg string, names, colors map[string]string) error {
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

func (l Poset) Nodes() (out []string) {
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

func (l Poset) Edges() (out []gviz.Edge) {
	for from, parent := range l.m {
		for to := range parent.Children {
			out = append(out, edge{from: from, to: to})
		}
	}
	return
}

func (l Poset) Children(a string) map[string]*Node {
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
func (l Poset) Meet(a, b string) string {
	bn := l.m[b]
	return l.BreadthFirstSearch(a, func(id string) bool {
		if bn.Descendents.Has(id) {
			return true
		}
		return false
	})
}

func NewNode(id string) *Node {
	return &Node{
		ID:          id,
		Children:    make(Nodeset),
		Descendents: make(Nodeset),
	}
}

func Generate(r *rand.Rand, chains, length int) Poset {
	out := Poset{
		m: make(map[string]*Node),
	}
	add := func(n *Node) {
		out.m[n.ID] = n
	}
	newNode := func(name string) *Node {
		if n, ok := out.m[name]; ok {
			return n
		}
		return NewNode(name)
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
	for _, n := range out.m {
		n.Descendents = out.CalcDescendents(n.ID)
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
