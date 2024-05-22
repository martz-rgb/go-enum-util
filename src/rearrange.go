package main

import (
	"go/ast"
	"go/token"
	"strings"
)

func Rearrange(tree *ast.File, dict *ast.ValueSpec, types []ast.TypeSpec) (*ast.File, error) {
	index := len(tree.Decls) - 1

	for index > -1 {
		switch v := tree.Decls[index].(type) {
		case *ast.GenDecl:
			if v.Tok == token.TYPE {
				// i feel like len always would be 1, but okay
				j := len(v.Specs) - 1

				for j > -1 {
					var ok bool
					types, ok = matchType(types, v.Specs[j])

					if ok {
						v.Specs = append(v.Specs[:j], v.Specs[j+1:]...)
					}
					j--
				}

				if len(v.Specs) == 0 {
					tree.Decls = append(tree.Decls[:index], tree.Decls[index+1:]...)
				}
			}

			if v.Tok == token.VAR {
				j := len(v.Specs) - 1

				for j > -1 {
					pos := matchDict(dict, v.Specs[j])

					if pos >= 0 {
						s := v.Specs[j].(*ast.ValueSpec)

						s.Names = append(s.Names[:pos], s.Names[pos+1:]...)
						s.Values = append(s.Values[:pos], s.Values[pos+1:]...)

						if len(s.Names) == 0 {
							v.Specs = append(v.Specs[:j], v.Specs[j+1:]...)
						}
					}

					j--
				}

				if len(v.Specs) == 0 {
					tree.Decls = append(tree.Decls[:index], tree.Decls[index+1:]...)
				}
			}

			index--
		}
	}

	types_decls := toDecls(types)

	tree.Decls = append(tree.Decls, types_decls...)
	tree.Decls = append(tree.Decls, &ast.GenDecl{
		Tok:   token.VAR,
		Specs: []ast.Spec{dict},
	})

	return tree, nil
}

func matchType(types []ast.TypeSpec, spec ast.Spec) ([]ast.TypeSpec, bool) {
	t, ok := spec.(*ast.TypeSpec)
	if !ok {
		return types, false
	}

	name := t.Name.String()

	flag := false
	for i := range types {
		if strings.Compare(types[i].Name.String(), name) == 0 {
			types[i] = *t
			flag = true
		}
	}

	return types, flag
}

func matchDict(dict *ast.ValueSpec, spec ast.Spec) int {
	v, ok := spec.(*ast.ValueSpec)
	if !ok {
		return -1
	}

	// dangerous, but it should be so
	name := dict.Names[0].String()

	for i := range v.Names {
		if strings.Compare(name, v.Names[i].String()) == 0 {
			return i
		}
	}

	return -1
}

func toDecls(types []ast.TypeSpec) []ast.Decl {
	decls := make([]ast.Decl, len(types))

	for i := range types {
		decls[i] = &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{&types[i]},
		}
	}

	return decls
}
