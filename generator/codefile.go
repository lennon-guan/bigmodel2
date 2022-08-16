package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
)

type CodeFile struct {
	fileset *token.FileSet
	file    *ast.File
	content []byte
}

func NewCodeFile(filename string) (*CodeFile, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	return &CodeFile{
		fileset: fset,
		file:    f,
		content: content,
	}, nil
}

func (cf *CodeFile) GetText(node ast.Node) string {
	p0 := node.Pos()
	begin := cf.fileset.File(p0).Position(p0).Offset
	p1 := node.End()
	end := cf.fileset.File(p1).Position(p1).Offset
	re := regexp.MustCompile(`\s+`)
	return string(re.ReplaceAll(cf.content[begin:end], []byte(" ")))
}
