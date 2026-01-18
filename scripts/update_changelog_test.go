package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// scripts/ -> repo root
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out))
}

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func copyFile(t *testing.T, from, to string, mode os.FileMode) {
	t.Helper()
	data, err := os.ReadFile(from)
	if err != nil {
		t.Fatalf("read %s: %v", from, err)
	}
	if err := os.WriteFile(to, data, mode); err != nil {
		t.Fatalf("write %s: %v", to, err)
	}
}

func TestUpdateChangelog_DetectBreakingFromBody(t *testing.T) {
	root := repoRoot(t)
	scriptSrc := filepath.Join(root, "scripts", "update_changelog.sh")

	tmp := t.TempDir()
	runGit(t, tmp, "init", "-b", "main")
	runGit(t, tmp, "config", "user.name", "test")
	runGit(t, tmp, "config", "user.email", "test@example.com")

	writeFile(t, filepath.Join(tmp, "README.md"), "base\n")
	runGit(t, tmp, "add", ".")
	runGit(t, tmp, "commit", "-m", "feat: base")
	runGit(t, tmp, "tag", "v1.0.0")

	// This commit's subject is "fix: ..." but BREAKING marker exists in the body.
	writeFile(t, filepath.Join(tmp, "CHANGELOG.md"), "placeholder\n")
	runGit(t, tmp, "add", ".")
	runGit(t, tmp, "commit", "-m", "fix: breaking in body", "-m", "BREAKING CHANGE: remove something")
	runGit(t, tmp, "tag", "v1.1.0")

	// Copy the script into the temp repo so it uses the temp repo root.
	scriptsDir := filepath.Join(tmp, "scripts")
	if err := os.MkdirAll(scriptsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	scriptDst := filepath.Join(scriptsDir, "update_changelog.sh")
	copyFile(t, scriptSrc, scriptDst, 0o755)

	cmd := exec.Command("bash", scriptDst, "v1.1.0")
	cmd.Dir = tmp
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run update_changelog.sh failed: %v\n%s", err, string(out))
	}

	got, err := os.ReadFile(filepath.Join(tmp, "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("read changelog: %v", err)
	}

	changelog := string(got)
	if !strings.Contains(changelog, "### Breaking Changes") {
		t.Fatalf("expected Breaking Changes section, got:\n%s", changelog)
	}
	if !strings.Contains(changelog, "- fix: breaking in body") {
		t.Fatalf("expected breaking bullet for commit subject, got:\n%s", changelog)
	}
}
