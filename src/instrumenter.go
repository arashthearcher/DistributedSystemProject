package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"reflect"
	"strings"
)

var fset *token.FileSet
var astFile *ast.File

func main() {
	src := `
// This is the package comment.
package main

import (
	"fmt"
)


// This comment is associated with the hello constant.
const hello = "Hello, World!" // line comment 1

// This comment is associated with the foo variable.
var foo = hello // line comment 2

//@dump This comment is associated with the main function.
func main() {
	//@dump
	fmt.Println(hello) // line comment 3
	fmt.Println("a")
	hello = b
	c = hello
	fmt.Println(foo) //@dump
	//@dump
}
`
	initializeInstrumenter(src)
	addImports()
	printer.Fprint(os.Stdout, fset, astFile)
}

func initializeInstrumenter(src string) {
	// Create the AST by parsing src.
	fset = token.NewFileSet() // positions are relative to fset
	astFile, _ = parser.ParseFile(fset, "test_program.go", src, parser.ParseComments)

	// Print the AST.
	ast.Print(fset, astFile)

	global_objs := astFile.Scope.Objects

	var global_vars string = ""
	for identifier, _ := range global_objs {
		global_vars += fmt.Sprintf("%v, ", identifier)
	}

	fmt.Println("Global Vars: " + global_vars)

	filePos := fset.File(astFile.Package)
	path, _ := astutil.PathEnclosingInterval(astFile, filePos.Pos(410), filePos.Pos(412))

	var pathStr string = ""
	for _, astnode := range path {
		pathStr += fmt.Sprintf("%v:%v -> ", astnode.Pos(), reflect.TypeOf(astnode))
		//fmt.Println(astutil.NodeDescription(astnode))
	}

	fmt.Println(pathStr)
	//ast.Walk(new(ImportVisitor), astFile)
	//printer.Fprint(os.Stdout, fset, f)

}

type ImportVisitor struct{}

func (v *ImportVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.Ident:
		fmt.Println(t.Pos())
		fmt.Println(t.Name)
	case *ast.FuncDecl:
		t.Name = ast.NewIdent(strings.Title(t.Name.Name))
	case *ast.Comment:
		if strings.Contains(t.Text, "@dump") {
			fmt.Println("dump encountered!")
			fmt.Println(t.Text)

		}
	case *ast.AssignStmt:
		fmt.Println(fmt.Sprintf("%v", t.Rhs[0]) + " -> " + fmt.Sprintf("%v", t.Lhs[0]))

	}

	return v
}

func GetAccessibleVarsInScope(lineNumber int) []string {
	return nil
}

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
			return im
		}
		return nil
	}
	return im
}
