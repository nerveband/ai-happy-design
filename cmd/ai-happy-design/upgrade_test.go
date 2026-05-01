package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpgradeTargetPathsCanonicalizesAhdAlias(t *testing.T) {
	target, alias := upgradeTargetPaths(filepath.Join(string(filepath.Separator), "tmp", "bin", "ahd-figma"))
	if target != filepath.Join(string(filepath.Separator), "tmp", "bin", "ai-happy-design") {
		t.Fatalf("unexpected target: %s", target)
	}
	if alias != filepath.Join(string(filepath.Separator), "tmp", "bin", "ahd-figma") {
		t.Fatalf("unexpected alias: %s", alias)
	}
}

func TestUpgradeTargetPathsMirrorsCanonicalToAlias(t *testing.T) {
	target, alias := upgradeTargetPaths(filepath.Join(string(filepath.Separator), "tmp", "bin", "ai-happy-design"))
	if target != filepath.Join(string(filepath.Separator), "tmp", "bin", "ai-happy-design") {
		t.Fatalf("unexpected target: %s", target)
	}
	if alias != filepath.Join(string(filepath.Separator), "tmp", "bin", "ahd-figma") {
		t.Fatalf("unexpected alias: %s", alias)
	}
}

func TestCopyExecutablePreservesExecutableMode(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "ai-happy-design")
	dst := filepath.Join(dir, "ahd-figma")
	if err := os.WriteFile(src, []byte("binary"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := copyExecutable(src, dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "binary" {
		t.Fatalf("unexpected copy contents: %q", data)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0755 {
		t.Fatalf("unexpected mode: %v", info.Mode().Perm())
	}
}
