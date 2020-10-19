// Package lattice is for computing lattice things like meet and join
package lattice

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/xoba/turd/cnfg"
)

func Run(cnfg.Config) error {
	chain := Generate(3)
	return chain.ToGraphViz(nil)
}

// perhaps open up in a browser, highlighting specific nodes with colors
func (l Lattice) ToGraphViz(colors map[string]string) error {
	f, err := os.Create("lattice.gv")
	if err != nil {
		return err
	}
	id := func(n string) string {
		return n
	}
	fmt.Fprintf(f, `digraph {
`)
	for _, n := range l.Nodes {
		fmt.Fprintf(f, "%s [ label=%q ];\n", id(n.ID), n.ID)
	}
	for _, n := range l.Nodes {
		for _, c := range n.Children {
			fmt.Fprintf(f, "%s -> %s [ label=%q ];\n", id(n.ID), id(c), "")
		}
	}
	fmt.Fprintf(f, "}\n")
	f.Close()
	return graphviz("dot", "lattice.gv", "lattice.svg", "svg")
}

type Node struct {
	ID       string
	Children []string
	Time     time.Time // to help order, assuming timestamps approx. correct
}

// a graph of nodes where every two has a unique meet (semi-lattice)
type Lattice struct {
	Nodes []Node
}

// returns meet of two nodes, if any
func (l Lattice) Meet(a, b string) string {
	panic("")
}

func Generate(chains int) (out Lattice) {
	add := func(n Node) {
		out.Nodes = append(out.Nodes, n)
	}
	genesis := Node{
		ID:   "g",
		Time: time.Now(),
	}
	add(genesis)
	last := genesis
	for i := 0; i < 5; i++ {
		chain := Node{
			ID:       fmt.Sprintf("%d", i),
			Children: []string{last.ID},
		}
		add(chain)
		last = chain
	}
	return
}

func graphviz(graphvizCommand, gv, out, format string) error {
	if err := exec.Command(graphvizCommand, "-v", "-o", out, fmt.Sprintf("-T%s", format), gv).Run(); err != nil {
		return fmt.Errorf("can't run graphviz (%s): %v", graphvizCommand, err)
	}
	return nil
}
