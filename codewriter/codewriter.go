package codewriter

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sato11/the-hack-vm-translator/parser"
)

// CodeWriter translates VM commands into Hack assembly code.
type CodeWriter struct {
	filename string
	eqIndex  int
	gtIndex  int
	ltIndex  int
	writer   *bytes.Buffer
}

// New opens file in write mode to write translations into.
func New() *CodeWriter {
	var buffer bytes.Buffer
	return &CodeWriter{
		"",
		0,
		0,
		0,
		&buffer,
	}
}

// SetFileName informs the code writer that the translation is started.
func (c *CodeWriter) SetFileName(filename string) {
	c.filename = filename
}

func binaryCommandOperator(command string) string {
	switch command {
	case "add":
		return "+"
	case "sub":
		return "-"
	case "and":
		return "&"
	case "or":
		return "|"
	default:
		panic(fmt.Errorf("%s is not a valid binary command", command))
	}
}

func unaryCommandOperator(command string) string {
	switch command {
	case "neg":
		return "-"
	case "not":
		return "!"
	default:
		panic(fmt.Errorf("%s is not a valid unary command", command))
	}
}

func (c *CodeWriter) getIndex(command string) int {
	switch command {
	case "eq":
		return c.eqIndex
	case "gt":
		return c.gtIndex
	case "lt":
		return c.ltIndex
	default:
		panic(fmt.Errorf("%s is not a valid command to get index for", command))
	}
}

func (c *CodeWriter) incrementIndex(command string) {
	switch command {
	case "eq":
		c.eqIndex++
	case "gt":
		c.gtIndex++
	case "lt":
		c.ltIndex++
	}
}

// WriteArithmetic writes the assembly code that is the translation of the given arithmetic command.
func (c *CodeWriter) WriteArithmetic(command string) {
	code := ""

	switch command {
	case "add", "sub", "and", "or":
		code = "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@SP\n" +
			"M=M-1\n" +
			"A=M\n"

		// invert order for subtraction
		if command == "sub" {
			code += "D=-D\n" +
				"M=-M\n"
		}

		code += fmt.Sprintf("M=D%sM\n", binaryCommandOperator(command)) +
			"@SP\n" +
			"M=M+1\n"

	case "neg", "not":
		code = "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			fmt.Sprintf("M=%sM\n", unaryCommandOperator(command)) +
			"@SP\n" +
			"M=M+1\n"

	case "eq", "lt", "gt":
		upperCommand := strings.ToUpper(command)
		labelIndex := c.getIndex(command)
		code =
			fmt.Sprintf("@CHECK%s%d\n", upperCommand, labelIndex) +
				"0;JMP\n" +
				fmt.Sprintf("(IS%s%d)\n", upperCommand, labelIndex) +
				"@SP\n" +
				"A=M\n" +
				"M=-1\n" +
				fmt.Sprintf("@%sEND%d\n", upperCommand, labelIndex) +
				"0;JMP\n" +
				fmt.Sprintf("(CHECK%s%d)\n", upperCommand, labelIndex) +
				"@SP\n" +
				"M=M-1\n" +
				"A=M\n" +
				"D=M\n" +
				"@SP\n" +
				"M=M-1\n" +
				"A=M\n" +
				"D=D-M\n" +
				"D=-D\n" +
				fmt.Sprintf("@IS%s%d\n", upperCommand, labelIndex) +
				fmt.Sprintf("D;J%s\n", upperCommand) +
				"@SP\n" +
				"A=M\n" +
				"M=0\n" +
				fmt.Sprintf("(%sEND%d)\n", upperCommand, labelIndex) +
				"@SP\n" +
				"M=M+1\n"

		c.incrementIndex(command)
	}

	c.writer.WriteString(code)
}

// WritePushPop writes the assembly code that is the translation of the given command,
// where command is either PushCommand or PopCommand.
func (c *CodeWriter) WritePushPop(command parser.CommandTypes, segment string, index int) {
	code := ""

	switch segment {
	case "constant":
		code += fmt.Sprintf("@%d\n", index) +
			"D=A\n"
	}

	switch command {
	case parser.PushCommand:
		code += "@SP\n" +
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
