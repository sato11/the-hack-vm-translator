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
		line = strings.Trim(line, " \t")    // remove whitespaces
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

// CommandTypes represent the return value for func CommandType()
type CommandTypes int

const (
	// ArithmeticCommand represents arithmetic or logical operation.
	ArithmeticCommand CommandTypes = iota
	// PushCommand represents push command.
	PushCommand
	// PopCommand represents pop command.
	PopCommand
	// LabelCommand represents label command.
	LabelCommand
	// GotoCommand represents goto command.
	GotoCommand
	// IfCommand represents if-goto command.
	IfCommand
	// FunctionCommand represents function command.
	FunctionCommand
	// ReturnCommand represents return command.
	ReturnCommand
	// CallCommand represents call command.
	CallCommand
)

// CommandType returns the type of the current VM command.
// Arithmetic is returned for all the arithmetic commands.
func (p *Parser) CommandType() CommandTypes {
	switch p.Command() {
	case "push":
		return PushCommand
	case "pop":
		return PopCommand
	case "label":
		return LabelCommand
	case "goto":
		return GotoCommand
	case "if-goto":
		return IfCommand
	case "function":
		return FunctionCommand
	case "call":
		return CallCommand
	case "return":
		return ReturnCommand
	default:
		return ArithmeticCommand
	}
}

// Command returns
func (p *Parser) Command() string {
	return strings.Split(p.currentCommand, " ")[0]
}

// Arg1 returns the first argument of the current command.
// In the case of ArithmeticCommand, the command itself (add, sub, etc.) is returned.
// Should not be called if the current command is ReturnCommand.
func (p *Parser) Arg1() string {
	if p.CommandType() == ArithmeticCommand {
		return p.currentCommand
	}
	return strings.Split(p.currentCommand, " ")[1]
}

// Arg2 returns the second argument of the current command.
// Should be called only if the current command is PushCommand, PopCommand, FunctionCommand or CallCommand.
func (p *Parser) Arg2() string {
	return strings.Split(p.currentCommand, " ")[2]
}
