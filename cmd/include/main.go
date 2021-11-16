// Package main contains the entry point for the go-include generator.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/halimath/include/internal/include"
)

func main() {
	outFile := flag.String("out", "", "Name of the output file to write to")
	buildTag := flag.String("buildtag", "include", "Build tag to deactivate in generate source")

	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "%s: Missing input file\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	inFile := flag.Arg(0)

	out, err := include.IncludeFile(inFile, include.Options{
		BuildTag: *buildTag,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(2)
	}

	if *outFile == "" {
		fmt.Print(string(out))
		return
	}

	if err := os.WriteFile(*outFile, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Failed to write %s: %s\n", os.Args[0], *outFile, err)
		os.Exit(3)
	}
}
