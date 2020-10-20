package gviz

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os/exec"
	"path/filepath"
)

type Graph interface {
	Nodes() []string
	Edges() []Edge
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
	var ids []string
	for _, k := range g.Nodes() {
		ids = append(ids, k)
	}
	for _, i := range ids {
		c := colors[i]
		if c == "" {
			c = "white"
		}
		fmt.Fprintf(f, "%s [ label=%q; fillcolor=%s style=filled ];\n", id(i), i, c)
	}
	for _, e := range g.Edges() {
		fmt.Fprintf(f, "%s -> %s;\n", id(e.From()), id(e.To()))
	}
	fmt.Fprintf(f, "}\n")
	return f.Bytes(), nil
}

func Dot(gv, out string) error {
	switch ext := filepath.Ext(out); ext {
	case ".svg":
		return Graphviz("dot", gv, out, "svg")
	default:
		return fmt.Errorf("unhandled extension: %q", ext)
	}
}

func Graphviz(graphvizCommand, gv, out, format string) error {
	if err := exec.Command(graphvizCommand, "-v", "-o", out, fmt.Sprintf("-T%s", format), gv).Run(); err != nil {
		return fmt.Errorf("can't run graphviz (%s): %v", graphvizCommand, err)
	}
	return nil
}
