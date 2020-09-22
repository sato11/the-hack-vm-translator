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
