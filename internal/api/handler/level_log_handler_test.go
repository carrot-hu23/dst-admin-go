package handler

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileLogReaderRangeReadsBoundedLineWindow(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server_log.txt")
	if err := os.WriteFile(path, []byte("l1\nl2\nl3\nl4\nl5\n"), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := NewFileLogReader().Range(path, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 5 || page.Start != 2 || page.End != 3 {
		t.Fatalf("unexpected metadata: %+v", page)
	}
	if !reflect.DeepEqual(page.Lines, []string{"l2", "l3"}) {
		t.Fatalf("unexpected lines: %#v", page.Lines)
	}
}

func TestFileLogReaderTailReturnsLastWindowWithMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server_log.txt")
	if err := os.WriteFile(path, []byte("l1\nl2\nl3\nl4\nl5\n"), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := NewFileLogReader().Tail(path, 2)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 5 || page.Start != 4 || page.End != 5 {
		t.Fatalf("unexpected metadata: %+v", page)
	}
	if !reflect.DeepEqual(page.Lines, []string{"l4", "l5"}) {
		t.Fatalf("unexpected lines: %#v", page.Lines)
	}
}
