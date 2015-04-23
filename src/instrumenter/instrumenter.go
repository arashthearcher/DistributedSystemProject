package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	//"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	//"os"
	//"reflect"
	"strings"
)

var fset *token.FileSet
var astFile *ast.File

func main() {
	initializeInstrumenter()
	dumpNodes := GetDumpNodes()
	addImports()

	for _, dumps := range dumpNodes {
		fmt.Println(GetAccessibleVarsInScope(int(dumps.Slash), astFile))
	}

	//printer.Fprint(os.Stdout, fset, astFile)
}

func initializeInstrumenter() {
	// Create the AST by parsing src.
	fset = token.NewFileSet() // positions are relative to fset
	astFile, _ = parser.ParseFile(fset, "../TestPrograms/serverUDP.go", nil, parser.ParseComments)

	// Print the AST.
	ast.Print(fset, astFile)

	//collectVars(149, astFile)

	//fmt.Println(pathStr)

	//ast.Walk(new(ImportVisitor), astFile)
	//printer.Fprint(os.Stdout, fset, f)

}

func GetAccessibleVarsInScope(start int, file *ast.File) []string {
	var results []string

	global_objs := file.Scope.Objects

	for identifier, _ := range global_objs {
		results = append(results, fmt.Sprintf("%v, ", identifier))
	}

	filePos := fset.File(file.Package)
	path, _ := astutil.PathEnclosingInterval(file, filePos.Pos(start), filePos.Pos(start+2))

	for _, astnode := range path {

		switch t := astnode.(type) {
		case *ast.BlockStmt:

			stmts := t.List
			for _, stmtnode := range stmts {
				switch t := stmtnode.(type) {
				case *ast.DeclStmt:
					idents := t.Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Names
					for _, identifier := range idents {
						results = append(results, fmt.Sprintf("%v, ", identifier.Name))

					}

				}
			}

		}
	}

	return results
}

func GetDumpNodes() []*ast.Comment {
	var dumpNodes []*ast.Comment
	for _, commentGroup := range astFile.Comments {
		for _, comment := range commentGroup.List {
			if strings.Contains(comment.Text, "@dump") {
				dumpNodes = append(dumpNodes, comment)
			}
		}
	}
	return dumpNodes
}

//type DumpVisitor struct{}

//func (v DumpVisitor) Visit(node ast.Node) (w ast.Visitor) {
//	//fmt.Println(node)
//	switch t := node.(type) {
//	case *ast.Comment:

//		if strings.Contains(t.Text, "@dump") {
//			fmt.Println("dump encountered !!")
//			dumpNodes = append(dumpNodes, t)
//		}
//	}

//	return v
//}

// returns dump code that should replace that specific line number
func GenerateDumpCode(vars []string, lineNumber int) string {
	var buffer bytes.Buffer

	// write vars' values
	buffer.WriteString(fmt.Sprintf("vars%d := []interface{}{", lineNumber))

	for i := 0; i < len(vars)-1; i++ {
		buffer.WriteString(fmt.Sprintf("%s,", vars[i]))
	}
	buffer.WriteString(fmt.Sprintf("%s}\n", vars[len(vars)-1]))

	// write vars' names
	buffer.WriteString(fmt.Sprintf("varsName%d := []string{", lineNumber))

	for i := 0; i < len(vars)-1; i++ {
		buffer.WriteString(fmt.Sprintf("\"%s\",", vars[i]))
	}
	buffer.WriteString(fmt.Sprintf("\"%s\"}\n", vars[len(vars)-1]))

	buffer.WriteString(fmt.Sprintf("point%d := createPoint(vars%d, varNames%d, %d)", lineNumber, lineNumber, lineNumber, lineNumber))
	buffer.WriteString(fmt.Sprintf("encoder.Encode(point%d)", lineNumber))

	return buffer.String()
}

//func AddRequiredStructsToProgram() {
//	code := `

//var encoder gob.Encoder
//func InstrumenterInit() {
//	fileW, _ := os.Create("log.txt")
//	encoder = gob.NewEncoder(fileW)
//}
//var logger *govec.GoLog
//func createPoint(vars []interface{}, varNames []string, lineNumber int) Point {

//	length := len(varNames)
//	dumps := make([]NameValuePair, length)
//	for i := 0; i < length; i++ {
//		dumps[i].VarName = varNames[i]
//		dumps[i].Value = vars[i]
//		dumps[i].Type = reflect.TypeOf(vars[i]).String()
//	}
//	point := Point{dumps, lineNumber, logger.currentVC}

//	return point
//}

//type Point struct {
//	Dump        []NameValuePair
//	LineNumber  string
//	vectorClock []byte
//}

//type NameValuePair struct {
//	VarName string
//	Value   interface{}
//	Type    string
//}

//func (nvp NameValuePair) String() string {
//	return fmt.Sprintf("(%s,%s,%s)", nvp.VarName, nvp.Value, nvp.Type)
//}

//func (p Point) String() string {
//	return fmt.Sprintf("%d : %s", p.LineNumber, p.Dump)
//}`
//}

func addImports() {
	packagesToImport := []string{"\"bytes\"", "\"encoding/gob\"", "\"reflect\"", "\"./govec\""}
	im := ImportAdder{packagesToImport}
	ast.Walk(im, astFile)
	ast.SortImports(fset, astFile)

}

type ImportAdder struct {
	PackagesToImport []string
}

func (im ImportAdder) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.GenDecl:
		if t.Tok == token.IMPORT {
			newSpecs := make([]ast.Spec, len(t.Specs)+len(im.PackagesToImport))
			for i, spec := range t.Specs {
				newSpecs[i] = spec
			}
			for i, spec := range im.PackagesToImport {
				newPackage := &ast.BasicLit{token.NoPos, token.STRING, spec}
				newSpecs[len(t.Specs)+i] = &ast.ImportSpec{nil, nil, newPackage, nil, token.NoPos}
			}

			t.Specs = newSpecs
			return nil
		}
	}
	return im
}
