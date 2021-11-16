// Package include provides functions and types that execute code inclusion in go files.
package include

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/imports"
)

const DEFAULT_BUILD_TAG = "include"

// Options defines the options to be passed when including files.
type Options struct {
	// BuildTag defines the build tag to toggle when generating code.
	BuildTag string

	// WorkingDir defines the current working directory used to resolve relative file name. If not set
	// defaults to the process' current working directory.
	WorkingDir string
}

// IncludeFile performs the inclusion reading source code from src. It returns the generated source as well
// as any error occured during processing.
func IncludeFile(src string, options Options) ([]byte, error) {
	content, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %s: %s", src, err)
	}

	return Include(src, content, options)
}

// Include processes the given source and returns the processed source as well as any error. Filename is only
// used in error messages.
func Include(filename string, source []byte, options Options) ([]byte, error) {
	if options.BuildTag == "" {
		options.BuildTag = DEFAULT_BUILD_TAG
	}

	var out bytes.Buffer

	out.WriteString("//Code generated by include. DO NOT EDIT.\n\n")

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filename, source, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source file %s: %s", filename, err)
	}

	var lastWriteOffset int

	var inspectErr error

	ast.Inspect(fileAst, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		f, ok := n.(*ast.File)
		if !ok {
			return false
		}

		for _, cg := range f.Comments {
			for _, c := range cg.List {
				if strings.Contains(c.Text, "go:build "+options.BuildTag) ||
					strings.Contains(c.Text, "+build "+options.BuildTag) ||
					strings.Contains(c.Text, "go:generate include") {
					out.Write(source[lastWriteOffset : c.Pos()-1])
					lastWriteOffset = int(c.End())
				}
			}
		}

		out.WriteString(fmt.Sprintf("//go:build !%s\n", options.BuildTag))
		out.WriteString(fmt.Sprintf("// +build !%s\n", options.BuildTag))

		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, s := range g.Specs {
				v, ok := s.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, n := range v.Values {
					c, ok := n.(*ast.CallExpr)
					if !ok {
						continue
					}

					f, ok := c.Fun.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					ident, ok := f.X.(*ast.Ident)
					if !ok {
						continue
					}
					if ident.Name != "include" {
						continue
					}

					if len(c.Args) != 1 {
						continue
					}

					var filename string

					if fn, ok := c.Args[0].(*ast.BasicLit); !ok {
						inspectErr = fmt.Errorf("unsupported argument when calling include.%s in %s: only strings are supported", f.Sel.Name, filename)
						return false
					} else {
						if fn.Value[0] != '"' || fn.Value[len(fn.Value)-1] != '"' {
							inspectErr = fmt.Errorf("unsupported argument when calling include.%s in %s: only strings are supported", f.Sel.Name, filename)
							return false
						}
						filename = fn.Value[1 : len(fn.Value)-1]
					}

					var replacement string

					if f.Sel.Name == "String" {
						replacement, err = fileContentAsString(filename, options)
						if err != nil {
							inspectErr = err
							return false
						}
					} else if f.Sel.Name == "Bytes" {
						replacement, err = fileContentAsBytes(filename, options)
						if err != nil {
							inspectErr = err
							return false
						}
					} else {
						inspectErr = fmt.Errorf("unsupported function call include.%s in %s", f.Sel.Name, filename)
						return false
					}

					out.Write(source[lastWriteOffset : c.Pos()-1])
					out.WriteString(replacement)
					lastWriteOffset = int(c.End())
				}
			}
		}

		return false
	})

	if inspectErr != nil {
		return nil, inspectErr
	}

	out.Write(source[lastWriteOffset:])

	c, err := imports.Process(filename, out.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to process imports: %s", err)
	}

	return c, nil
}

func fileContentAsString(filename string, options Options) (string, error) {
	content, err := readFile(filename, options)
	if err != nil {
		return "", err
	}

	return "`" + strings.ReplaceAll(string(content), "`", "\\`") + "`", nil
}

func fileContentAsBytes(filename string, options Options) (string, error) {
	content, err := readFile(filename, options)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("[]byte{\n")

	for i, b := range content {
		builder.WriteString(strconv.Itoa(int(b)))
		builder.WriteRune(',')
		if i > 0 && i%20 == 0 {
			builder.WriteRune('\n')
		}
	}

	builder.WriteString("\n}\n")

	return builder.String(), nil
}

func readFile(filename string, options Options) ([]byte, error) {
	if options.WorkingDir != "" {
		filename = filepath.Join(options.WorkingDir, filename)
	}
	fqfn, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", filename, err)
	}

	content, err := os.ReadFile(fqfn)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", filename, err)
	}

	return content, nil
}
