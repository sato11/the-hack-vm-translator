package codewriter

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/sato11/the-hack-vm-translator/parser"
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

// WriteArithmetic writes the assembly code that is the translation of the given arithmetic command.
func (c *CodeWriter) WriteArithmetic(command string) {
	code := ""

	switch command {
	case "add":
		code = "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"M=D+M\n" +
			"@SP\n" +
			"M=M+1"
	}

	c.writer.WriteString(code)
}

// WritePushPop writes the assembly code that is the translation of the given command,
// where command is either PushCommand or PopCommand.
func (c *CodeWriter) WritePushPop(command parser.CommandTypes, segment string, index int) {
	code := ""

	switch command {
	case parser.PushCommand:
		code = fmt.Sprintf("@%d\n", index) +
			"D=A\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"
	default:
		panic(errors.New("codewriter.WritePushPop only accepts PushCommand and PopCommand"))
	}

	c.writer.WriteString(code)
}

// Save writes the output to file.
func (c *CodeWriter) Save() {
	f, err := os.Create(c.filename)
	if err != nil {
		panic(err)
	}

	f.Write(c.writer.Bytes())
}
