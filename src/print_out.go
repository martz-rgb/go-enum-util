package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

func PrintOut(filename string, tree *ast.File) error {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err := format.Node(&buf, fset, tree)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
