package loader

import (
	"fmt"
	"log"

	"golang.org/x/tools/go/packages"
)

// LoadPackages loads Go packages from the given patterns.
// It returns only packages that contain Go files and logs any loading errors.
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

		for _, e := range pkg.Errors {
			fmt.Println("Load error:", e)
		}
	}

	return pkgsWithFile
}
