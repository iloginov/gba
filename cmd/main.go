package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/iloginov/gba"
)

var modules map[string]gba.ModuleDesc

func scanSubDir(dir string, sub string) {

	files, err := ioutil.ReadDir(path.Join(dir, sub))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		matched, err := path.Match("*.a", f.Name())
		if err != nil {
			log.Error(err)
			continue
		}
		if matched {
			var mod gba.ModuleDesc
			var ok bool
			if mod, ok = modules[sub]; !ok {
				mod = gba.ModuleDesc{
					Deps: []string{},
				}
			}
			mod.Size = f.Size()
			modules[sub] = mod
			log.Infof("Size of module %s is %d", sub, f.Size())
		}

		if f.Name() == "importcfg" || f.Name() == "importcfg.link" {
			file, err := os.Open(path.Join(dir, sub, f.Name()))
			if err != nil {
				log.Error(err)
				continue
			}

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			var mod gba.ModuleDesc
			var ok bool
			if mod, ok = modules[sub]; !ok {
				mod = gba.ModuleDesc{
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
					buildDir := strings.TrimPrefix(strings.TrimSuffix(parts[1], "/_pkg_.a"), dir)

					if buildDir != sub {
						mod.Deps = append(mod.Deps, buildDir)

						var mod1 gba.ModuleDesc
						var ok1 bool
						if mod1, ok1 = modules[buildDir]; !ok1 {
							mod1 = gba.ModuleDesc{
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
}

func main() {
	dir := "/tmp/go-build656671572/"

	modules = make(map[string]gba.ModuleDesc)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			scanSubDir(dir, f.Name())
		}
	}

	for k, v := range modules {
		log.Info(k, " -> ", v)
	}

	dot, err := gba.MakeDotFile(modules)
	if err != nil {
		log.Error(err)
		return
	}

	ioutil.WriteFile("graph.dot", dot, 0644)

	log.Info("Finished")
}
