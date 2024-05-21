package main

import (
	"flag"
	"go/ast"
	"log"
)

var (
	constType = flag.String("const", "", "seek for constants only with this type")
	dictType  = flag.String("dict", "", "desired type for dictionary entries")
	dictName  = flag.String("name", "", "desired name for dictionary")
	output    = flag.String("output", "", "set output file")
)

// enumutil takes files with global constants declarations from Go file
// and create or modify another Go file, where creates data types with prefix
// for constants and create dictionary with all constants.
// If const type is not provided, it uses first found.

func main() {
	log.SetFlags(0)
	log.SetPrefix("enumutil: ")

	flag.Parse()

	if constType == nil {
		log.Fatal("there is should be const type")
	}
	if dictType == nil {
		log.Fatal("there is should be dictionary type")
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("there is should at least one file")
	}

	// get all global constants & type of constants & desired type for map
	var struct_type *ast.StructType
	var consts []string

	for i := range args {
		cs, err := ParseIn(args[i], *constType, *dictType, &struct_type)
		if err != nil {
			log.Fatal(err)
		}

		consts = append(consts, cs...)
	}

	// parse out if exists and create new otherwise

}
