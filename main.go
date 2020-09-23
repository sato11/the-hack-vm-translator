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
			w.WritePushPop(parser.PushCommand, p.Arg1(), index)
		}
	}

	return nil
}

// main reads single file when argument is vm file.
// otherwise recursively searches for vm files under the given path.
func main() {
	path := os.Args[1]
	codewriter := codewriter.New()

	extension := filepath.Ext(path)
	codewriter.SetFileName(fmt.Sprintf("%s.asm", strings.TrimSuffix(path, extension)))

	if extension == ".vm" {
		err := translateFile(path, codewriter)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(ExitCodeError)
		}
	} else {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".vm" {
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
