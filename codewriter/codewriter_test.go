package codewriter

import (
	"path"
	"testing"
)

func TestSetFileName(t *testing.T) {
	filenames := []string{
		"filename.asm",
		path.Join("path", "to", "file.asm"),
	}
	for i, filename := range filenames {
		c := New()
		c.SetFileName(filename)
		if c.filename != filename {
			t.Errorf("#%d: got: %v wanted: %v", i, c.filename, filename)
		}
	}
}
