package parser

import (
	"bufio"
	"io"
	"strings"
)

// Parser handles the parsing of a single .vm file and encapsulates access to the input code.
type Parser struct {
	currentCommand string
	lines          []string
}

// New initializes the parser and gets ready to parse the input stream.
func New(r io.Reader) *Parser {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, "//")[0] // remove comments
		line = strings.Trim(line, " ")      // remove whitespaces
		if line != "" {
			lines = append(lines, line)
		}
	}

	return &Parser{
		"",
		lines,
	}
}

// HasMoreCommands returns true if there are more commands to parse.
func (p *Parser) HasMoreCommands() bool {
	return len(p.lines) != 0
}

// Advance reads the next command from the input and makes it the current command.
// Should be called only if HasMoreCommands() is true.
func (p *Parser) Advance() {
	p.currentCommand = p.lines[0]
	p.lines = p.lines[1:]
}
