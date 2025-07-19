package loader

import (
	"fmt"
	"log"

	"golang.org/x/tools/go/packages"
)

func LoadPackages(patterns ...string) []*packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatalf("failed to load packages: %v", err)
	}

	if len(pkgs) == 0 {
		log.Println("Warning: no packages matched the given pattern(s)")
	}

	for _, pkg := range pkgs {
		fmt.Printf("Loaded package: %s, GoFiles: %d\n", pkg.Name, len(pkg.GoFiles))
	}

	return pkgs
}
