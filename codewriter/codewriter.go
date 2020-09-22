package codewriter

import (
	"bytes"
)

// CodeWriter translates VM commands into Hack assembly code.
type CodeWriter struct {
	filename string
	writer   *bytes.Buffer
}

// New opens file in write mode to write translations into.
func New() *CodeWriter {
	var buffer bytes.Buffer
	return &CodeWriter{
		"",
		&buffer,
	}
}

// SetFileName informs the code writer that the translation is started.
func (c *CodeWriter) SetFileName(filename string) {
	c.filename = filename
}
