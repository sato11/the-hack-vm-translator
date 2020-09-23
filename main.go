package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sato11/the-hack-vm-translator/codewriter"
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

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	fmt.Printf(string(b))
	return nil
}

// main reads single file when argument is vm file.
// otherwise recursively searches for vm files under the given path.
func main() {
	path := os.Args[1]
	codewriter := codewriter.New()
	codewriter.SetFileName(fmt.Sprintf("%s.asm", filepath.Base(path)))

	if filepath.Ext(path) == ".vm" {
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

	os.Exit(ExitCodeOK)
}
