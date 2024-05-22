package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func CompareStructType(name string, spec ast.Spec) (*ast.TypeSpec, bool) {
	t, ok := spec.(*ast.TypeSpec)
	if !ok {
		// error
		return nil, false
	}

	if strings.Compare(name, t.Name.String()) == 0 {
		return t, true
	}

	return nil, false
}

func CompareValueType(name string, spec ast.Spec) ([]string, bool) {
	v, ok := spec.(*ast.ValueSpec)
	if !ok {
		return nil, false
	}

	i, ok := v.Type.(*ast.Ident)
	if !ok {
		return nil, false
	}

	if strings.Compare(name, i.Name) == 0 {
		names := make([]string, len(v.Names))
		for i := range v.Names {
			names[i] = v.Names[i].Name
		}

		return names, true
	}

	return nil, false
}

func ParseIn(filename string, constType string, dictType string) ([]string, *ast.TypeSpec, error) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var consts []string
	var st *ast.TypeSpec

	ast.Inspect(tree, func(x ast.Node) bool {
		decl, ok := x.(*ast.GenDecl)
		if !ok {
			return true
		}

		if decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if t, ok := CompareStructType(dictType, spec); ok {
					st = t
				}
			}
		}

		if decl.Tok == token.CONST {
			for _, spec := range decl.Specs {
				if names, ok := CompareValueType(constType, spec); ok {
					consts = append(consts, names...)
				}
			}
		}
		return false
	})

	return consts, st, nil
}
