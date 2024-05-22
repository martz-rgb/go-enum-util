package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

var nameUtil = "enumutil"

var (
	constType   = flag.String("const", "", "seek for constants only with this type")
	suffix      = flag.String("suffix", "", "suffix which is added to const name to name corresponding struct")
	dictType    = flag.String("dict", "", "desired type for dictionary entries")
	dictName    = flag.String("name", "", "desired name for dictionary")
	packageName = flag.String("package", "", "package if output file is empty")
	output      = flag.String("output", "", "set output file")
)

// enumutil takes files with global constants declarations from Go file
// and create or modify another Go file, where creates data types with prefix
// for constants and create dictionary with all constants.
// If const type is not provided, it uses first found.

func main() {
	log.SetFlags(0)
	log.SetPrefix("enumutil: ")

	flag.Parse()

	if len(*output) == 0 {
		log.Fatal("there is should be output file")
	}
	if len(*suffix) == 0 {
		log.Fatal("there is should be suffix")
	}
	if len(*constType) == 0 {
		log.Fatal("there is should be const type")
	}
	if len(*dictName) == 0 {
		log.Fatal("there is should be dictionary name")
	}
	if len(*dictType) == 0 {
		log.Fatal("there is should be dictionary type")
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("there is should at least one file")
	}

	// get all global constants & type of constants & desired type for map
	var struct_type *ast.TypeSpec
	var consts []string

	for i := range args {
		cs, st, err := ParseIn(args[i], *constType, *dictType)
		if err != nil {
			log.Fatal(err)
		}

		if st != nil {
			struct_type = st
		}

		consts = append(consts, cs...)
	}

	// generate structs
	sc := Analysis(struct_type)
	dict, types, err := Generate(consts, *suffix, *constType, *dictName, *dictType, sc)
	if err != nil {
		log.Fatal(err)
	}

	// parse output if exists and create new otherwise
	var out *ast.File

	if _, err := os.Stat(*output); err == nil {
		fset := token.NewFileSet()
		out, err = parser.ParseFile(fset, *output, nil, 0)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if len(*packageName) == 0 {
			log.Fatal("package should be provided if output file is empty")
		}

		out = &ast.File{
			Name: &ast.Ident{
				Name: *packageName,
			},
		}
	}

	out, err = Rearrange(out, dict, types)
	if err != nil {
		log.Fatal(err)
	}

	// print to file
	err = PrintOut(*output, out)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("succesfully generated file")
}
