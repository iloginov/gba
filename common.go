package gba

import (
	"bytes"
	"fmt"
	"io"
)

type ModuleDesc struct {
	Name string
	Size int64
	Deps []string
}

// DotConfig contains attributes about how a graph should be
// constructed and how it should look.
type DotConfig struct {
	Title     string   // The title of the DOT graph
	LegendURL string   // The URL to link to from the legend.
	Labels    []string // The labels for the DOT's legend

	FormatValue func(int64) string // A formatting function for values
	Total       int64              // The total weight of the graph, used to compute percentages
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

func (b *builder) addNode(id string, m ModuleDesc) {
	fmt.Fprintf(b, "%s [shape=box label=\"%s\\nSize=%s\"]\n", id, m.Name, ByteCountBinary(m.Size))
}

func (b *builder) addEdge(from string, to string) {
	fmt.Fprintf(b, "%s -> %s\n", from, to)
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func MakeDotFile(modules map[string]ModuleDesc) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	builder := &builder{buffer, &DotConfig{}}

	builder.start()

	for k, v := range modules {
		builder.addNode(k, v)
	}

	for k, m := range modules {
		for _, e := range m.Deps {
			builder.addEdge(k, e)
		}
	}

	builder.finish()

	return buffer.Bytes(), nil
}
