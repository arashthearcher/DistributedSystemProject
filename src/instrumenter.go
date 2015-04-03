package main

import (
	//"bytes"
	"fmt"
	"go/ast"
	//"go/format"
	"go/parser"
	//"go/printer"
	"go/token"
	//"os"
	"strings"
)

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
	instrument(src)

}

func instrument(src string) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "test_program.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)

	ast.Walk(new(ImportVisitor), f)
	//printer.Fprint(os.Stdout, fset, f)

}

type ImportVisitor struct{}

func (v *ImportVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
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
