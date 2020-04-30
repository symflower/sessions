package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/zimmski/osutil"
)

// TODO story about program.

// TODO Concurrency is not parallelism https://blog.golang.org/concurrency-is-not-parallelism

func main() {
	// NOTE We need to define the file/directory argument with a flag because the go tools will not allow us to define a "go file" as an argument. Go tools want to compile it.
	var path string
	flag.StringVar(&path, "path", "", "the file or path to analyse")
	flag.Parse()

	// NOTE We currently go through directories and files sequentially.
	// TODO parallel traversing in one process
	// TODO parallel traversing with multiple processes
	log.Print("Traverse path")
	files, err := osutil.FilesRecursive(path)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet() // TODO Should we share a fileSet? How can know if we don't?

	// TODO parallel iterating in one process
	// TODO parallel iterating with multiple processes
	log.Print("Iterate found files")
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") {
			continue
		}

		log.Printf("Handle file %s", file)

		// TODO We read the file content with the parser functionality. Should we do this concurrently?
		// TODO Do we really need to parse the wholte file for the comments? Go comments are really easy by parsing them with a regex.
		ast, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			// NOTE Let's not care about syntax errors.
			log.Print(err)
		}

		for _, commentGroup := range ast.Comments {
			for _, comment := range commentGroup.List {
				fmt.Println(comment.Text)
			}
		}
	}
}
