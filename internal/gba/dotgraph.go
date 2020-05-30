package gba

import (
	"fmt"
	"io"
)

type DotConfig struct {
	Title string
}

type builder struct {
	io.Writer
	config *DotConfig
}

func (b *builder) start() {
	graphname := "unnamed"
	if b.config.Title != "" {
		graphname = b.config.Title
	}
	fmt.Fprintln(b, `digraph "`+graphname+`" {`)
	fmt.Fprintln(b, `node [style=filled fillcolor="#f8f8f8"]`)
}

func (b *builder) finish() {
	fmt.Fprintln(b, "}")
}

func (b *builder) addNode(id string, n *Node) {
	fmt.Fprintf(b, "%s [shape=box label=\"%s\\nSize=%s\\nDeps size=%s\"]\n", id, n.Name, ByteCountBinary(n.Size), ByteCountBinary(n.DepsSize))
}

func (b *builder) addEdge(from string, to string) {
	fmt.Fprintf(b, "%s -> %s\n", from, to)
}
