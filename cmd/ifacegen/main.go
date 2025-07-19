package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/seyyedaghaei/ifacegen/internal/generator"
	"github.com/seyyedaghaei/ifacegen/internal/loader"
)

func main() {
	matchFlag := flag.String("match", "", "Comma-separated glob patterns to match struct names (e.g. *Service)")
	outputFlag := flag.String("output", "iface_gen.go", "Name of the generated interface file")
	nameFlag := flag.String("name", "I{}", "Naming pattern for interfaces. '{}' replaced by struct name")
	helpFlag := flag.Bool("help", false, "Show usage information")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage:
  ifacegen [flags] <package pattern>

Example:
  ifacegen -match=*Service,*Repository ./...

Flags:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	pkgs := loader.LoadPackages(flag.Args()...)
	matches := splitMatches(*matchFlag)

	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup

	for _, pkg := range pkgs {
		wg.Add(1)
		pkg := pkg
		go func(pkg *packages.Package) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			code, err := generator.Generate(pkg, *nameFlag, matches)
			if err != nil {
				log.Printf("Error generating for package %s: %v", pkg.PkgPath, err)
				return
			}
			if code == nil {
				return
			}
			outFile := filepath.Join(filepath.Dir(pkg.GoFiles[0]), *outputFlag)
			if err := writeIfChanged(outFile, code); err != nil {
				log.Printf("Error writing file %s: %v", outFile, err)
				return
			}
			log.Printf("Generated interfaces for package %s", pkg.PkgPath)
		}(pkg)
	}
	wg.Wait()
}

func splitMatches(s string) []string {
	var res []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			res = append(res, part)
		}
	}
	return res
}

func writeIfChanged(path string, content []byte) error {
	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}
	return os.WriteFile(path, content, 0644)
}
