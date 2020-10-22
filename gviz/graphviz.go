package gviz

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Graph interface {
	Nodes() []Node
	Edges() []Edge
}

type Node interface {
	ID() string
	Label() string
	Group() string
}

type Edge interface {
	From() string
	To() string
}

// Compile creates a gv output
func Compile(g Graph, colors map[string]string) ([]byte, error) {
	id := func(name string) string {
		h := md5.New()
		h.Write([]byte(name))
		return fmt.Sprintf("N%x", h.Sum(nil))
	}
	f := new(bytes.Buffer)
	fmt.Fprintf(f, "digraph {\n")
	//fmt.Fprintf(f, "rankdir = LR\n")
	groups := make(map[string][]Node)
	for _, n := range g.Nodes() {
		groups[n.Group()] = append(groups[n.Group()], n)
	}
	var x int
	for _, v := range groups {
		//fmt.Fprintf(f, "subgraph cluster_%d {\n", x)
		x++
		for _, n := range v {
			c := colors[n.ID()]
			if c == "" {
				c = "white"
			}
			label := n.Label()
			if g := n.Group(); g != "" {
				label = g + "/" + label
			}
			fmt.Fprintf(f, "%s [ label=%q; fillcolor=%s style=filled ];\n", id(n.ID()), label, c)
		}
		//fmt.Fprintf(f, "}\n")
	}
	for _, e := range g.Edges() {
		fmt.Fprintf(f, "%s -> %s;\n", id(e.From()), id(e.To()))
	}
	fmt.Fprintf(f, "}\n")
	return f.Bytes(), nil
}

func Dot(in, out string) error {
	switch ext := filepath.Ext(out); ext {
	case ".svg":
		return graphviz("dot", in, out, "svg")
	default:
		return fmt.Errorf("unhandled extension: %q", ext)
	}
}

func graphviz(graphvizCommand, in, out, format string) error {
	cmd := exec.Command(graphvizCommand, "-v", "-o", out, fmt.Sprintf("-T%s", format), in)
	if false {
		fmt.Printf("cmd = %q\n", cmd.Args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("can't run graphviz (%s): %v", graphvizCommand, err)
	}
	return nil
}
