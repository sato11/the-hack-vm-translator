package codewriter

import (
	"path"
	"strings"
	"testing"

	"github.com/sato11/the-hack-vm-translator/parser"
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

func TestSetFunctionName(t *testing.T) {
	functions := []string{
		"power",
		"sqrt",
	}
	for i, functionName := range functions {
		c := New()
		c.SetFunctionName(functionName)
		if c.functionName != functionName {
			t.Errorf("#%d: got: %v wanted: %v", i, c.functionName, functionName)
		}
	}
}

type writeArithmeticTest struct {
	command string
	out     string
}

func TestWriteArithmetic(t *testing.T) {
	tests := []writeArithmeticTest{
		{"add", "@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nM=D+M\n@SP\nM=M+1\n"},
		{"sub", "@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nD=-D\nM=-M\nM=D-M\n@SP\nM=M+1\n"},
		{"and", "@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nM=D&M\n@SP\nM=M+1\n"},
		{"or", "@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nM=D|M\n@SP\nM=M+1\n"},
		{"neg", "@SP\nM=M-1\nA=M\nM=-M\n@SP\nM=M+1\n"},
		{"not", "@SP\nM=M-1\nA=M\nM=!M\n@SP\nM=M+1\n"},
		{"eq", "@CHECKEQ0\n0;JMP\n(ISEQ0)\n@SP\nA=M\nM=-1\n@EQEND0\n0;JMP\n(CHECKEQ0)\n@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nD=D-M\nD=-D\n@ISEQ0\nD;JEQ\n@SP\nA=M\nM=0\n(EQEND0)\n@SP\nM=M+1\n"},
		{"gt", "@CHECKGT0\n0;JMP\n(ISGT0)\n@SP\nA=M\nM=-1\n@GTEND0\n0;JMP\n(CHECKGT0)\n@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nD=D-M\nD=-D\n@ISGT0\nD;JGT\n@SP\nA=M\nM=0\n(GTEND0)\n@SP\nM=M+1\n"},
		{"lt", "@CHECKLT0\n0;JMP\n(ISLT0)\n@SP\nA=M\nM=-1\n@LTEND0\n0;JMP\n(CHECKLT0)\n@SP\nM=M-1\nA=M\nD=M\n@SP\nM=M-1\nA=M\nD=D-M\nD=-D\n@ISLT0\nD;JLT\n@SP\nA=M\nM=0\n(LTEND0)\n@SP\nM=M+1\n"},
	}

	for i, test := range tests {
		c := New()
		c.WriteArithmetic(test.command)
		if c.writer.String() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, c.writer.String(), test.out)
		}
	}
}

type writePushPopTest struct {
	commandType parser.CommandTypes
	segment     string
	index       int
	out         string
}

func TestWritePushPop(t *testing.T) {
	tests := []writePushPopTest{
		{parser.PushCommand, "constant", 420, "@420\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "local", 1, "@LCL\nA=M\nA=A+1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "argument", 1, "@ARG\nA=M\nA=A+1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "this", 1, "@THIS\nA=M\nA=A+1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "that", 1, "@THAT\nA=M\nA=A+1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "temp", 1, "@R5\nA=A+1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "pointer", 0, "@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "pointer", 1, "@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PushCommand, "static", 1, "@filename.1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{parser.PopCommand, "local", 1, "@SP\nM=M-1\nA=M\nD=M\n@LCL\nA=M\nA=A+1\nM=D\n"},
		{parser.PopCommand, "argument", 1, "@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nA=A+1\nM=D\n"},
		{parser.PopCommand, "this", 1, "@SP\nM=M-1\nA=M\nD=M\n@THIS\nA=M\nA=A+1\nM=D\n"},
		{parser.PopCommand, "that", 1, "@SP\nM=M-1\nA=M\nD=M\n@THAT\nA=M\nA=A+1\nM=D\n"},
		{parser.PopCommand, "temp", 1, "@SP\nM=M-1\nA=M\nD=M\n@R5\nA=A+1\nM=D\n"},
		{parser.PopCommand, "pointer", 0, "@SP\nM=M-1\nA=M\nD=M\n@THIS\nM=D\n"},
		{parser.PopCommand, "pointer", 1, "@SP\nM=M-1\nA=M\nD=M\n@THAT\nM=D\n"},
		{parser.PopCommand, "static", 1, "@SP\nM=M-1\nA=M\nD=M\n@filename.1\nM=D\n"},
	}

	for i, test := range tests {
		c := New()
		c.SetFileName("filename")
		c.WritePushPop(test.commandType, test.segment, test.index)
		if c.writer.String() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, c.writer.String(), test.out)
		}
	}
}

type labelTest struct {
	functionName string
	label        string
	out          string
}

func TestWriteLabel(t *testing.T) {
	tests := []labelTest{
		{"power", "LOOP", "(power$LOOP)\n"},
		{"sqrt", "LABEL", "(sqrt$LABEL)\n"},
	}

	for i, test := range tests {
		c := New()
		c.SetFunctionName(test.functionName)
		c.WriteLabel(test.label)
		if c.writer.String() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, c.writer.String(), test.out)
		}
	}
}

func TestWriteGoto(t *testing.T) {
	tests := []labelTest{
		{"power", "LOOP", "@power$LOOP\n0;JMP\n"},
		{"sqrt", "LABEL", "@sqrt$LABEL\n0;JMP\n"},
	}

	for i, test := range tests {
		c := New()
		c.SetFunctionName(test.functionName)
		c.WriteGoto(test.label)
		if c.writer.String() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, c.writer.String(), test.out)
		}
	}
}

func TestWriteIf(t *testing.T) {
	tests := []labelTest{
		{"power", "LOOP", "@SP\nM=M-1\nA=M\nD=M\n@power$LOOP\nD;JNE\n"},
		{"sqrt", "LABEL", "@SP\nM=M-1\nA=M\nD=M\n@sqrt$LABEL\nD;JNE\n"},
	}

	for i, test := range tests {
		c := New()
		c.SetFunctionName(test.functionName)
		c.WriteIf(test.label)
		if c.writer.String() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, c.writer.String(), test.out)
		}
	}
}

func TestWriteReturn(t *testing.T) {
	c := New()
	c.WriteReturn()

	actual := c.writer.String()
	expected := strings.Join([]string{
		"@LCL",
		"D=M",
		"@R13",
		"M=D",
		"D=M",
		"@5",
		"A=D-A",
		"D=M",
		"@R14",
		"M=D",
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"@ARG",
		"A=M",
		"M=D",
		"@ARG",
		"D=M+1",
		"@SP",
		"M=D",
		"@R13",
		"A=M-1",
		"D=M",
		"@THAT",
		"M=D",
		"@2",
		"D=A",
		"@R13",
		"A=M-D",
		"D=M",
		"@THIS",
		"M=D",
		"@3",
		"D=A",
		"@R13",
		"A=M-D",
		"D=M",
		"@ARG",
		"M=D",
		"@4",
		"D=A",
		"@R13",
		"A=M-D",
		"D=M",
		"@LCL",
		"M=D",
		"@R14",
		"A=M",
		"0;JMP",
	}, "\n") + "\n"

	if actual != expected {
		t.Errorf("got: %v wanted: %v", actual, expected)
	}
}

type functionCallTest struct {
	functionName string
	numLocals    int
	out          string
}

func TestWriteFunction(t *testing.T) {
	tests := []functionCallTest{
		{"Sqrt", 1, "(Sqrt)\n@SP\nA=M\nM=0\n@SP\nM=M+1\n"},
		{"Power", 2, "(Power)\n@SP\nA=M\nM=0\n@SP\nM=M+1\n@SP\nA=M\nM=0\n@SP\nM=M+1\n"},
	}

	for i, test := range tests {
		c := New()
		c.WriteFunction(test.functionName, test.numLocals)
		if c.writer.String() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, c.writer.String(), test.out)
		}
	}
}
