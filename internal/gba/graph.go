package gba

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	ID       string
	Name     string
	Size     int64
	Deps     []*Node
	DepsSize int64

	HasParent bool
}

type ModuleGraph struct {
	Root  *Node
	Nodes map[string]*Node
}

type dirDesc struct {
	Name string
	Size int64
	Deps []string
}

func getDeps(n *Node) []*Node {
	if n.Deps == nil {
		return nil
	}

	m := make(map[string]*Node)
	for _, v := range n.Deps {
		_, ok := m[v.ID]
		if !ok {
			m[v.ID] = v
		}
	}

	for _, v := range n.Deps {
		deps := getDeps(v)
		for _, v1 := range deps {
			_, ok := m[v1.ID]
			if !ok {
				m[v1.ID] = v1
			}
		}
	}

	res := []*Node{}
	for _, v := range m {
		res = append(res, v)
	}

	return res
}

func calcDepsCost(n *Node) {
	deps := getDeps(n)

	if deps == nil {
		n.DepsSize = 0
		return
	}

	var cost int64
	for _, v := range deps {
		cost = cost + v.Size
	}

	n.DepsSize = cost

	if n.Deps != nil {
		for _, v := range n.Deps {
			calcDepsCost(v)
		}
	}
}

func scanSubDir(dir string, sub string, modules map[string]dirDesc) (map[string]dirDesc, error) {
	files, err := ioutil.ReadDir(path.Join(dir, sub))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		matched, err := path.Match("*.a", f.Name())
		if err != nil {
			log.Error(err)
			continue
		}
		if matched {
			var mod dirDesc
			var ok bool
			if mod, ok = modules[sub]; !ok {
				mod = dirDesc{
					Deps: []string{},
				}
			}
			mod.Size = f.Size()
			modules[sub] = mod
		}

		if f.Name() == "importcfg" || f.Name() == "importcfg.link" {
			file, err := os.Open(path.Join(dir, sub, f.Name()))
			if err != nil {
				log.Error(err)
				continue
			}

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			var mod dirDesc
			var ok bool
			if mod, ok = modules[sub]; !ok {
				mod = dirDesc{
					Deps: []string{},
				}
			}

			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "#") {
					continue
				}

				if strings.HasPrefix(line, "packagefile ") {
					parts := strings.Split(strings.Split(line, " ")[1], "=")
					name := parts[0]
					buildDir := strings.TrimPrefix(strings.TrimPrefix(strings.TrimSuffix(parts[1], "/_pkg_.a"), dir), "/")

					if buildDir != sub {
						mod.Deps = append(mod.Deps, buildDir)

						var mod1 dirDesc
						var ok1 bool
						if mod1, ok1 = modules[buildDir]; !ok1 {
							mod1 = dirDesc{
								Deps: []string{},
							}
						}
						mod1.Name = name
						modules[buildDir] = mod1
					} else {
						mod.Name = name
					}
				}
			}
			file.Close()

			modules[sub] = mod
		}
	}

	return modules, nil
}

func BuildModuleGraph(dir string) (*ModuleGraph, error) {
	mods := make(map[string]dirDesc)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			mods, err = scanSubDir(dir, f.Name(), mods)
			if err != nil {
				return nil, err
			}
		}
	}

	// Make graph
	graph := &ModuleGraph{
		Nodes: make(map[string]*Node),
	}
	for k, v := range mods {
		graph.Nodes[k] = &Node{
			ID:       k,
			Name:     v.Name,
			Size:     v.Size,
			Deps:     []*Node{},
			DepsSize: -1,
		}
	}
	for k, v := range mods {
		n := graph.Nodes[k]
		for _, v1 := range v.Deps {
			d := graph.Nodes[v1]
			d.HasParent = true
			n.Deps = append(n.Deps, d)
		}
	}
	for _, v := range graph.Nodes {
		if !v.HasParent {
			if graph.Root == nil {
				graph.Root = v
			} else {
				return nil, fmt.Errorf("Resulting graph has several roots")
			}
		}
	}

	// Calc dependencies cost
	calcDepsCost(graph.Root)
	log.Infof("Root deps cost: %s", ByteCountBinary(graph.Root.DepsSize))

	var totalCost int64
	for _, v := range graph.Nodes {
		totalCost = totalCost + v.Size
	}
	log.Infof("Total cost: %s", ByteCountBinary(totalCost))

	return graph, nil
}

func (g *ModuleGraph) MakeDotFile() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	builder := &builder{buffer, &DotConfig{}}

	builder.start()

	for _, v := range g.Nodes {
		builder.addNode(v.ID, v)
	}

	for _, v := range g.Nodes {
		for _, v1 := range v.Deps {
			builder.addEdge(v.ID, v1.ID)
		}
	}

	builder.finish()

	return buffer.Bytes(), nil
}

func (n *Node) stringHelper(prefix string, level int, maxLevel int, buf *bytes.Buffer) {
	if level > maxLevel {
		return
	}
	buf.WriteString(prefix)
	if level > 0 {
		buf.WriteString("├")
		buf.WriteString(strings.Repeat("─", (level*4)-2))
		buf.WriteString(" ")
	}
	buf.WriteString(n.Name + " " + ByteCountBinary(n.Size) + "\n")
	level++

	for _, ch := range n.Deps {
		ch.stringHelper(prefix, level, maxLevel, buf)
	}
}

func (g *ModuleGraph) PrintTree(level int) string {
	b := bytes.Buffer{}

	b.WriteString("\n")
	g.Root.stringHelper("", 0, level, &b)

	return b.String()
}
