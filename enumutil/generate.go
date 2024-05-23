package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
)

func Generate(consts []string, suffix string, constType string, dictName string, dictType string, sc *StructConfig) (*ast.ValueSpec, []ast.TypeSpec, error) {
	types := make([]ast.TypeSpec, len(consts))
	elements := make([]ast.Expr, len(consts))

	for i, c := range consts {
		types[i] = *generateType(c, suffix)
		elements[i] = sc.generateElement(c, suffix)
	}

	dict := &ast.ValueSpec{
		Names: []*ast.Ident{
			{Name: dictName},
		},

		Values: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.MapType{
					Key:   &ast.Ident{Name: constType},
					Value: &ast.Ident{Name: dictType},
				},
				Elts: elements,
			},
		},
	}

	return dict, types, nil
}

func generateType(c string, suffix string) *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: fmt.Sprintf("%s%s", c, suffix),
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{},
		},
	}
}

type StructConfig struct {
	AddType []string
}

func Analysis(t *ast.TypeSpec) *StructConfig {
	if t == nil {
		return &StructConfig{}
	}
	st, ok := t.Type.(*ast.StructType)
	if !ok {
		return &StructConfig{}
	}

	config := &StructConfig{}

	for _, field := range st.Fields.List {
		if field.Tag != nil {
			tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
			if err != nil {
				continue
			}

			tag, err := tags.Get(nameUtil)
			if err != nil {
				continue
			}

			switch tag.Name {
			case "type":
				names := make([]string, len(field.Names))
				for i, n := range field.Names {
					names[i] = n.String()
				}
				config.AddType = append(config.AddType, names...)
			}
		}
	}

	return config
}

func (sc *StructConfig) generateElement(c string, suffix string) ast.Expr {
	elem := &ast.KeyValueExpr{
		Key:   &ast.Ident{Name: c},
		Value: &ast.CompositeLit{},
	}

	if len(sc.AddType) == 0 {
		return elem
	}

	fields := make([]ast.Expr, len(sc.AddType))

	for i, t := range sc.AddType {
		fields[i] = &ast.KeyValueExpr{
			Key: &ast.Ident{Name: t},
			Value: &ast.CallExpr{
				Fun: &ast.ParenExpr{X: &ast.StarExpr{
					X: &ast.Ident{Name: fmt.Sprintf("%s%s", c, suffix)},
				}},
				Args: []ast.Expr{
					&ast.Ident{Name: "nil"},
				},
				Ellipsis: token.NoPos,
			},
		}
	}

	elem.Value = &ast.CompositeLit{
		Elts: fields,
	}

	return elem
}
