package generator_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/seyyedaghaei/ifacegen/internal/generator"
	"github.com/seyyedaghaei/ifacegen/internal/loader"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// internal/generator -> internal -> repo root
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}

func goldenPath(t *testing.T, filename string) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "testdata", "golden", filename)
}

func loadAndGenerate(t *testing.T, pkgPattern string, pkgName string, namePattern string, matches []string) []byte {
	t.Helper()

	root := repoRoot(t)
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	pkgs := loader.LoadPackages(pkgPattern)
	for _, pkg := range pkgs {
		if pkg.Name == pkgName {
			out, err := generator.Generate(pkg, namePattern, matches)
			if err != nil {
				t.Fatalf("generator.Generate: %v", err)
			}
			return out
		}
	}

	t.Fatalf("package not found: pattern=%q name=%q", pkgPattern, pkgName)
	return nil
}

func compareGolden(t *testing.T, filename string, got []byte) {
	t.Helper()

	wantFile := goldenPath(t, filename)
	gotStr := strings.TrimRight(string(got), "\n")

	update := os.Getenv("UPDATE_GOLDEN") == "1"
	if wantBytes, err := os.ReadFile(wantFile); err == nil {
		wantStr := strings.TrimRight(string(wantBytes), "\n")
		if gotStr != wantStr {
			t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", filename, gotStr, wantStr)
		}
		return
	} else if !update {
		t.Fatalf("missing golden file %s (set UPDATE_GOLDEN=1 to create it): %v", filename, err)
	}

	if err := os.MkdirAll(filepath.Dir(wantFile), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(wantFile, []byte(gotStr+"\n"), 0o644); err != nil {
		t.Fatalf("write golden: %v", err)
	}
}

func TestGenerator_Golden(t *testing.T) {
	tests := []struct {
		name        string
		pkgPattern  string
		pkgName     string
		namePattern string
		matches     []string
		goldenFile  string
	}{
		{
			name:        "match *Service includes skip/generate/method-skip",
			pkgPattern:  "./testdata/matchskip",
			pkgName:     "matchskip",
			namePattern: "I{}",
			matches:     []string{"*Service"},
			goldenFile:  "matchskip_match_service.go",
		},
		{
			name:        "match *Service,*Repository includes repository",
			pkgPattern:  "./testdata/matchskip",
			pkgName:     "matchskip",
			namePattern: "I{}",
			matches:     []string{"*Service", "*Repository"},
			goldenFile:  "matchskip_match_service_repository.go",
		},
		{
			name:        "import aliasing for same package names",
			pkgPattern:  "./testdata/alias/consumer",
			pkgName:     "consumer",
			namePattern: "I{}",
			matches:     []string{"*Service"},
			goldenFile:  "alias_import_aliasing.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loadAndGenerate(t, tt.pkgPattern, tt.pkgName, tt.namePattern, tt.matches)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}
