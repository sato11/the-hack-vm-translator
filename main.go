package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sato11/the-hack-vm-translator/codewriter"
	"github.com/sato11/the-hack-vm-translator/parser"
)

// ExitCodeOK and ExitCodeError represent respectively a status code.
const (
	ExitCodeOK int = iota
	ExitCodeError
)

func translateFile(path string, w *codewriter.CodeWriter) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	p := parser.New(f)
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.ArithmeticCommand:
			w.WriteArithmetic(p.Command())
		case parser.PushCommand, parser.PopCommand:
			index, err := strconv.Atoi(p.Arg2())
			if err != nil {
				return err
			}
			w.WritePushPop(p.CommandType(), p.Arg1(), index)
		case parser.LabelCommand:
			w.WriteLabel(p.Arg1())
		case parser.GotoCommand:
			w.WriteGoto(p.Arg1())
		case parser.IfCommand:
			w.WriteIf(p.Arg1())
		case parser.CallCommand:
			numArgs, err := strconv.Atoi(p.Arg2())
			if err != nil {
				return err
			}
			w.WriteCall(p.Arg1(), numArgs)
		case parser.ReturnCommand:
			w.WriteReturn()
		case parser.FunctionCommand:
			numLocals, err := strconv.Atoi(p.Arg2())
			if err != nil {
				return err
			}
			w.WriteFunction(p.Arg1(), numLocals)
		}
	}

	return nil
}

// main reads single file when argument is vm file.
// otherwise recursively searches for vm files under the given path.
func main() {
	path := os.Args[1]
	codewriter := codewriter.New()
	codewriter.Setup()

	extension := filepath.Ext(path)

	if extension == ".vm" {
		codewriter.SetFileName(fmt.Sprintf("%s.asm", strings.TrimSuffix(path, extension)))
		err := translateFile(path, codewriter)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(ExitCodeError)
		}
	} else {
		filename := filepath.Join(fmt.Sprintf("%s", path), fmt.Sprintf("%s.asm", path))
		codewriter.SetFileName(filename)
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".vm" {
				codewriter.SetNamespace(strings.TrimSuffix(filepath.Base(path), ".vm"))
				err := translateFile(path, codewriter)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(ExitCodeError)
				}
			}
			return nil
		})
		if err != nil {
			os.Exit(ExitCodeError)
		}
	}

	codewriter.Save()
	os.Exit(ExitCodeOK)
}
