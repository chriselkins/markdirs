package markdirs

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMarkDirs_CreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested dirs
	os.MkdirAll(filepath.Join(tmpDir, "a", "b"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "c"), 0755)

	// Mark all dirs with a .marker file
	content := []byte("hello")
	err := MarkDirs(tmpDir, ".marker", content, false, true, false, 0644)
	if err != nil {
		t.Fatalf("MarkDirs failed: %v", err)
	}

	// Check file exists in every dir
	for _, sub := range []string{"", "a", "a/b", "c"} {
		path := filepath.Join(tmpDir, sub, ".marker")
		data, err := os.ReadFile(path)

		if err != nil {
			t.Errorf("Missing marker in %s: %v", sub, err)
		} else if string(data) != "hello" {
			t.Errorf("Wrong content in %s: got %q", sub, string(data))
		}
	}
}

func TestMarkDirs_Overwrite(t *testing.T) {
	tmpDir := t.TempDir()

	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	path := filepath.Join(tmpDir, "sub", "file")
	os.WriteFile(path, []byte("old"), 0644)

	err := MarkDirs(tmpDir, "file", []byte("new"), true, true, false, 0644)
	if err != nil {
		t.Fatalf("MarkDirs failed: %v", err)
	}

	data, _ := os.ReadFile(path)

	if string(data) != "new" {
		t.Errorf("Overwrite failed: got %q", string(data))
	}
}

func TestMarkDirs_NoOverwrite(t *testing.T) {
	tmpDir := t.TempDir()

	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	path := filepath.Join(tmpDir, "sub", "file")
	os.WriteFile(path, []byte("old"), 0644)

	err := MarkDirs(tmpDir, "file", []byte("new"), false, true, false, 0644)
	if err != nil {
		t.Fatalf("MarkDirs failed: %v", err)
	}

	data, _ := os.ReadFile(path)

	if string(data) != "old" {
		t.Errorf("Should not have overwritten: got %q", string(data))
	}
}

func TestMarkDirs_Permissions(t *testing.T) {
	tmpDir := t.TempDir()

	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)

	// Set mode to 0600 for test
	perm := os.FileMode(0600)

	err := MarkDirs(tmpDir, "privatefile", []byte("secret"), false, true, false, perm)
	if err != nil {
		t.Fatalf("MarkDirs failed: %v", err)
	}

	info, err := os.Stat(filepath.Join(tmpDir, "privatefile"))
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if info.Mode().Perm() != perm {
		t.Errorf("File permission: got %v, want %v", info.Mode().Perm(), perm)
	}
}

func TestMarkDirsFromReaderAt_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	os.MkdirAll(filepath.Join(tmpDir, "d1"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "d2"), 0755)

	content := []byte("streamtest")
	reader := bytes.NewReader(content)

	err := MarkDirsFromReaderAt(tmpDir, ".streamed", reader, int64(len(content)), false, true, false, 0644)
	if err != nil {
		t.Fatalf("MarkDirsFromReaderAt failed: %v", err)
	}

	// Check files in all dirs
	for _, sub := range []string{"", "d1", "d2"} {
		path := filepath.Join(tmpDir, sub, ".streamed")

		data, err := os.ReadFile(path)

		if err != nil {
			t.Errorf("Missing .streamed in %s: %v", sub, err)
		} else if string(data) != "streamtest" {
			t.Errorf("Wrong .streamed content in %s: got %q", sub, string(data))
		}
	}
}
