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
	filename     string
	functionName string
	eqIndex      int
	gtIndex      int
	ltIndex      int
	writer       *bytes.Buffer
}

// New opens file in write mode to write translations into.
func New() *CodeWriter {
	var buffer bytes.Buffer
	return &CodeWriter{
		"",
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

// SetFunctionName informs the code writer that the translation is started.
func (c *CodeWriter) SetFunctionName(functionName string) {
	c.functionName = functionName
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

func (c CodeWriter) getIndex(command string) int {
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

func (c CodeWriter) handlePushCommand(segment string, index int) string {
	switch segment {
	case "constant":
		return fmt.Sprintf("@%d\n", index) +
			"D=A\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

	case "local":
		code := "@LCL\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

		return code

	case "argument":
		code := "@ARG\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

		return code

	case "this":
		code := "@THIS\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

		return code

	case "that":
		code := "@THAT\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

		return code

	case "temp":
		code := "@R5\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

		return code

	case "pointer":
		var code string

		if index == 0 {
			code = "@THIS\n"
		} else {
			code = "@THAT\n"
		}

		code += "D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

		return code

	case "static":
		return fmt.Sprintf("@%s.%d\n", c.filename, index) +
			"D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n"

	default:
		return ""
	}
}

func (c CodeWriter) handlePopCommand(segment string, index int) string {
	switch segment {
	case "local":
		code := "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@LCL\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "M=D\n"

		return code
	case "argument":
		code := "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@ARG\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "M=D\n"

		return code

	case "this":
		code := "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@THIS\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "M=D\n"

		return code

	case "that":
		code := "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@THAT\n" +
			"A=M\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "M=D\n"

		return code

	case "temp":
		code := "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			"@R5\n"

		for i := 0; i < index; i++ {
			code += "A=A+1\n"
		}

		code += "M=D\n"

		return code

	case "pointer":
		code := "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n"

		if index == 0 {
			code += "@THIS\n"
		} else {
			code += "@THAT\n"
		}

		code += "M=D\n"

		return code

	case "static":
		return "@SP\n" +
			"M=M-1\n" +
			"A=M\n" +
			"D=M\n" +
			fmt.Sprintf("@%s.%d\n", c.filename, index) +
			"M=D\n"

	default:
		return ""
	}
}

// WritePushPop writes the assembly code that is the translation of the given command,
// where command is either PushCommand or PopCommand.
func (c *CodeWriter) WritePushPop(command parser.CommandTypes, segment string, index int) {
	var code string

	switch command {
	case parser.PushCommand:
		code = c.handlePushCommand(segment, index)
	case parser.PopCommand:
		code = c.handlePopCommand(segment, index)
	default:
		panic(errors.New("codewriter.WritePushPop only accepts PushCommand and PopCommand"))
	}

	c.writer.WriteString(code)
}

// WriteLabel writes assembly code that effects the label command.
func (c *CodeWriter) WriteLabel(label string) {
	code := fmt.Sprintf("(%s$%s)\n", c.functionName, label)

	c.writer.WriteString(code)
}

// WriteGoto writes assembly code that effects the goto command.
func (c *CodeWriter) WriteGoto(label string) {
	code := fmt.Sprintf("@%s$%s\n", c.functionName, label) +
		"0;JMP\n"

	c.writer.WriteString(code)
}

// WriteIf writes assembly code that effects the if-goto command.
func (c *CodeWriter) WriteIf(label string) {
	code := "@SP\n" +
		"M=M-1\n" +
		"A=M\n" +
		"D=M\n" +
		fmt.Sprintf("@%s$%s\n", c.functionName, label) +
		"D;JNE\n"

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
