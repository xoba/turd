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
	Nodes() []string
	Edges() []Edge
}

type Edge interface {
	From() string
	To() string
}

// Compile creates a gv output
func Compile(g Graph, colors, names map[string]string) ([]byte, error) {
	init := func(m map[string]string) map[string]string {
		if m == nil {
			return make(map[string]string)
		}
		return m
	}
	names = init(names)
	colors = init(colors)
	if colors == nil {
		colors = make(map[string]string)
	}
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
		name := names[i]
		if name == "" {
			name = i
		}
		fmt.Fprintf(f, "%s [ label=%q; fillcolor=%s style=filled ];\n", id(i), name, c)
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
