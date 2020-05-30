package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/iloginov/gba/internal/gba"
)

func main() {
	var err error

	pkg := "gandalf/cmd/gandalf"
	//workDir := "/tmp/go-build656671572/"

	workDir, err := gba.BuildPackage(pkg)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Build directory: %s", workDir)

	graph, err := gba.BuildModuleGraph(workDir)
	if err != nil {
		log.Fatal(err)
	}

	// dot, err := graph.MakeDotFile()
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }

	// ioutil.WriteFile("graph.dot", dot, 0644)

	tree := graph.PrintTree(3)
	log.Info(tree)
}
