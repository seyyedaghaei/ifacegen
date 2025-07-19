package loader

import (
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

	var pkgsWithFile []*packages.Package

	for _, pkg := range pkgs {
		if len(pkg.GoFiles) > 0 {
			pkgsWithFile = append(pkgsWithFile, pkg)
		}
	}

	return pkgsWithFile
}
