// ProgramSlicer
package main

import (
	"./cfg"
	"fmt"
	//"github.com/godoctor/godoctor/analysis/dataflow"
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
)

// Create CFG
// For given graph compute dominator and post-dominators
// Create Control Dependence Graph
// Create Data Dependence Graph
// Create Program Dependence Graph

func main() {
	src := `
  package main

  func foo(c int, nums []int) int {
    //START
    a := 1
	b := 2
	if a == 1 {
		b = 1
	} else {
		b = 10
	}
	c = b + a
    //END
  }`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	funcOne := f.Decls[0].(*ast.FuncDecl)
	c := cfg.FromFunc(funcOne)
	_ = c.GetBlocks() // for 100% coverage ;)
	c.InitializeBlocks()
	invC := c.InvertCFG()
	var buf bytes.Buffer
	invC.PrintDot(&buf, fset, func(s ast.Stmt) string {
		if _, ok := s.(*ast.AssignStmt); ok {
			return "!"
		} else {
			return ""
		}
	})
	dot := buf.String()
	fmt.Println(invC.BlockSlice)
	cfg.BuildDomTree(invC)
	cfg.PrintDomTreeDot(&buf, invC, fset)
	dot = buf.String()

	fmt.Println(dot)
	//createMyCFG(c)

}
