package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/iloginov/gba/internal/gba"
)

var tree = flag.Bool("tree", false, "print dependency tree")
var level = flag.Int("level", 3, "max tree level to print (-1 to out whole tree)")
var dot = flag.Bool("dot", false, "create dependency graph in Graphviz dot format")
var file = flag.String("file", "graph.dot", "name of the output .dot file (default 'graph.dot'")

func main() {
	var err error

	flag.Parse()

	if !*tree && !*dot {
		fmt.Println("You should choose one of the output options")
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Println("You should give package name")
		return
	}
	pkg := flag.Arg(0)

	workDir, err := gba.BuildPackage(pkg)
	if err != nil {
		log.Fatal(err.Error())
	}

	graph, err := gba.BuildModuleGraph(pkg, workDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *dot {
		dot, err := graph.MakeDotFile()
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		ioutil.WriteFile(*file, dot, 0644)
	}

	if *tree {
		tree := graph.PrintTree(*level)
		fmt.Println(tree)
	}
}
