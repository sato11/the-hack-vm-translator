package parser

import (
	"bytes"
	"testing"
)

type newTest struct {
	reader string
	lines  []string
}

func TestNew(t *testing.T) {
	tests := []newTest{
		{"", []string{}},
		{"\n", []string{}},
		{"\n\n", []string{}},
		{"// comment", []string{}},
		{"push constant 0", []string{"push constant 0"}},
		{"// comment\npush constant 0", []string{"push constant 0"}},
		{"push constant 0\npop local 0", []string{"push constant 0", "pop local 0"}},
	}
	for i, test := range tests {
		b := bytes.NewBufferString(test.reader)
		p := New(b)
		if len(p.lines) != len(test.lines) {
			t.Errorf("#%d: got %v wanted %v", i, len(p.lines), len(test.lines))
		} else {
			for j := range p.lines {
				if p.lines[j] != test.lines[j] {
					t.Errorf("#%d: got: %v wanted: %v", j, p.lines[j], test.lines[j])
				}
			}
		}
	}
}

type hasMoreCommandsTest struct {
	lines []string
	out   bool
}

func TestHasMoreCommands(t *testing.T) {
	tests := []hasMoreCommandsTest{
		{[]string{}, false},
		{[]string{"push constant 0"}, true},
		{[]string{"push constant 0", "pop local 0"}, true},
	}
	for i, test := range tests {
		p := &Parser{"", test.lines}
		if p.HasMoreCommands() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, p.HasMoreCommands(), test.out)
		}
	}
}

type advanceTest struct {
	before  []string
	after   []string
	command string
}

func TestAdvance(t *testing.T) {
	tests := []advanceTest{
		{[]string{"push constant 0"}, []string{}, "push constant 0"},
		{[]string{"push constant 0", "pop local 0"}, []string{"pop local 0"}, "push constant 0"},
	}
	for i, test := range tests {
		p := &Parser{"", test.before}
		p.Advance()

		if p.currentCommand != test.command {
			t.Errorf("#%d: got: %v wanted: %v", i, p.currentCommand, test.command)
		}

		if len(p.lines) != len(test.after) {
			t.Errorf("#%d: got: %v wanted: %v", i, p.lines, test.after)
		} else {
			for j := range p.lines {
				if p.lines[j] != test.after[j] {
					t.Errorf("#%d: got: %v wanted: %v", j, p.lines[j], test.after[j])
				}
			}
		}
	}
}

type commandTypeTest struct {
	command string
	out     CommandTypes
}

func TestCommandType(t *testing.T) {
	tests := []commandTypeTest{
		{"push constant 0", PushCommand},
		{"pop location 1", PopCommand},
		{"label loop", LabelCommand},
		{"goto loop", GotoCommand},
		{"if-goto end", IfCommand},
		{"function mult 2", FunctionCommand},
		{"call mult 2 5", CallCommand},
		{"return", ReturnCommand},
		{"add", ArithmeticCommand},
		{"sub", ArithmeticCommand},
		{"lt", ArithmeticCommand},
	}
	for i, test := range tests {
		p := &Parser{test.command, []string{}}
		if p.CommandType() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, p.CommandType(), test.out)
		}
	}
}

type commandTest struct {
	command string
	out     string
}

func TestCommand(t *testing.T) {
	tests := []commandTest{
		{"push constant 0", "push"},
		{"pop location 1", "pop"},
		{"label loop", "label"},
		{"goto loop", "goto"},
		{"if-goto end", "if-goto"},
		{"function mult 2", "function"},
		{"call mult 2 5", "call"},
		{"return", "return"},
		{"add", "add"},
		{"sub", "sub"},
		{"lt", "lt"},
	}
	for i, test := range tests {
		p := &Parser{test.command, []string{}}
		if p.Command() != test.out {
			t.Errorf("#%d: got: %v wanted: %v", i, p.Command(), test.out)
		}
	}
}

func TestArg1(t *testing.T) {
	tests := []commandTest{
		{"push constant 0", "constant"},
		{"pop location 1", "location"},
		{"label loop", "loop"},
		{"goto loop", "loop"},
		{"if-goto end", "end"},
		{"function mult 2", "mult"},
		{"call mult 2 5", "mult"},
		{"add", "add"},
		{"sub", "sub"},
		{"lt", "lt"},
	}
	for i, test := range tests {
		p := &Parser{test.command, []string{}}
		if p.Arg1() != test.out {
			t.Errorf("#%d: got: %v wanted %v", i, p.Arg1(), test.out)
		}
	}
}

func TestArg2(t *testing.T) {
	tests := []commandTest{
		{"push constant 0", "0"},
		{"pop location 1", "1"},
		{"function mult 2", "2"},
		{"call mult 2 5", "2"},
	}
	for i, test := range tests {
		p := &Parser{test.command, []string{}}
		if p.Arg2() != test.out {
			t.Errorf("#%d: got: %v wanted %v", i, p.Arg2(), test.out)
		}
	}
}
