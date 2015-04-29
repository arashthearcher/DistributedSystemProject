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
	"testing"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"

	"github.com/godoctor/godoctor/analysis/cfg"
	"github.com/godoctor/godoctor/analysis/dataflow"
)

const (
	START = 0
	END   = 100000000
)

var fset *token.FileSet
var astFile *ast.File
var c *CFGWrapper

func main() {
	initializeInstrumenter()
	dumpNodes := GetDumpNodes()
	addImports()

	for _, dumps := range dumpNodes {
		//fmt.Println(GetAccessibleVarsInScope(int(dumps.Slash), astFile))
		fmt.Println(GetAccessedVarsInScope(dumps, astFile, c.f))

	}
	fmt.Println(detectSendReceive(astFile))
}

func initializeInstrumenter() {
	src_location := "../TestPrograms/serverUDP.go"
	// Create the AST by parsing src.
	fset = token.NewFileSet() // positions are relative to fset
	astFile, _ = parser.ParseFile(fset, src_location, nil, parser.ParseComments)

	// Print the AST.

	c = getWrapper(nil, src_location)

	ast.Print(fset, astFile)

	//fmt.Println(pathStr)

	//ast.Walk(new(ImportVisitor), astFile)
	//printer.Fprint(os.Stdout, fset, f)

}

func detectSendReceive(f *ast.File) []*ast.Node {
	var results []*ast.Node
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			switch y := x.Fun.(type) {
			case *ast.SelectorExpr:
				left, ok := y.X.(*ast.Ident)
				if ok && left.Name == "conn" && y.Sel.Name == "ReadFrom" || y.Sel.Name == "WriteTo" {
					fmt.Println(left.Name, y.Sel.Name)
					results = append(results, &n)
				}
			}
			return true
		}

		return true
	})
	return results
}

func GetAccessedVarsInScope(dumpNode *ast.Comment, f *ast.File, g *ast.File) []string {
	var results []string
	//filePos := fset.File(file.Package)
	//path, _ := astutil.PathEnclosingInterval(f, dumpNode.Pos(), dumpNode.End())
	path2, _ := astutil.PathEnclosingInterval(g, dumpNode.Pos(), dumpNode.End())

	var stmts []ast.Stmt

	for _, astnode := range path2 {

		funcDecl, ok := astnode.(*ast.FuncDecl)
		if ok { // skip import decl if exists

			ast.Inspect(funcDecl, func(n ast.Node) bool {
				switch x := n.(type) {
				case ast.Stmt:
					switch x.(type) {
					case *ast.BlockStmt:
						return true
					}
					if x.Pos() < dumpNode.Pos() {
						stmts = append(stmts, x)
					}
					//v[i] = x
					//stmts[x] = i
					//i++
				case *ast.FuncLit:
					// skip statements in anonymous functions
					return false
				}
				return true
			})
		}

		//fmt.Println("Decl:::%v", astutil.NodeDescription(astnode))
		//switch t := astnode.(type) {
		//case *ast.BlockStmt:

		//	stmts := t.List
		//	for _, stmtnode := range stmts {
		//		switch t := stmtnode.(type) {
		//		case *ast.DeclStmt:
		//			idents := t.Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Names
		//			for _, identifier := range idents {
		//				fmt.Println("Ident::%v, ", identifier.Name)

		//			}

		//		}
		//	}

		//}
	}
	//fmt.Println(stmts)
	_, uses := dataflow.ReferencedVars(stmts, c.prog.Created[0])

	//actualUse := make(map[*types.Var]struct{})
	for u, _ := range uses {
		results = append(results, u.Name())
	}

	return results

	//c := getWrapper(nil, "../TestPrograms/serverUDP.go")

	//for _, decl := range f.Decls {
	//	//fmt.Println(GetAccessibleVarsInScope(int(dumps.Slash), astFile))
	//	funcDecl, ok := decl.(*ast.FuncDecl)
	//	if ok { // skip import decl if exists
	//		fmt.Println("FnDecl:::%v", astutil.NodeDescription(funcDecl))
	//	}
	//}

	//firstFunc, ok := f.Decls[0].(*ast.FuncDecl)
	//if !ok { // skip import decl if exists
	//	firstFunc = f.Decls[1].(*ast.FuncDecl) // panic here if no first func
	//}
	//cfg := cfg.FromFunc(firstFunc)
	//v := make(map[int]ast.Stmt)
	//stmts := make(map[ast.Stmt]int)
	////objs := make(map[string]*types.Var)
	////objNames := make(map[*types.Var]string)
	//i := 1
	//ast.Inspect(firstFunc, func(n ast.Node) bool {
	//	switch x := n.(type) {
	//	case ast.Stmt:
	//		switch x.(type) {
	//		case *ast.BlockStmt:
	//			return true
	//		}
	//		v[i] = x
	//		stmts[x] = i
	//		i++
	//	case *ast.FuncLit:
	//		// skip statements in anonymous functions
	//		return false
	//	}
	//	return true
	//})
	//v[END] = cfg.Exit
	//v[START] = cfg.Entry
	//stmts[cfg.Entry] = START
	//stmts[cfg.Exit] = END
	//if len(v) != len(cfg.Blocks()) {
	//	fmt.Errorf("expected %d vertices, got %d --construction error", len(v), len(cfg.Blocks()))
	//}

	////c.expectUses(t, START, 2, "c")
	////end := 11
	////start := START

	////c.printAST()
	////blocks := c.cfg.Blocks()
	////info := c.prog.Created[0]
	////in, _ := ReachingDefs(c.cfg, info)
	////ins := in[c.exp[s]]

	//if _, ok := c.stmts[c.exp[0]]; !ok {
	//	fmt.Println("did not find start", 0)
	//	return
	//}
	//if _, ok := c.stmts[dumpNode]; !ok {
	//	fmt.Println("did not find end", dumpNode)
	//	return
	//}

	//var stmts []ast.Stmt
	//for i := start; i <= end; i++ { // include end
	//	stmts = append(stmts, c.exp[i])
	//}

}

func GetAccessibleVarsInScope(start int, file *ast.File) []string {

	var results []string

	global_objs := astFile.Scope.Objects

	for identifier, _ := range global_objs {
		results = append(results, fmt.Sprintf("%v, ", identifier))
	}

	filePos := fset.File(astFile.Package)
	path, _ := astutil.PathEnclosingInterval(astFile, filePos.Pos(start), filePos.Pos(start+2))

	for _, astnode := range path {
		//fmt.Println("%v", astutil.NodeDescription(astnode))
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

type CFGWrapper struct {
	cfg      *cfg.CFG
	prog     *loader.Program
	exp      map[int]ast.Stmt
	stmts    map[ast.Stmt]int
	objs     map[string]*types.Var
	objNames map[*types.Var]string
	fset     *token.FileSet
	f        *ast.File
}

// uses first function in given string to produce CFG
// w/ some other convenient fields for printing in test
// cases when need be...
func getWrapper(t *testing.T, filename string) *CFGWrapper {
	var config loader.Config
	f, err := config.ParseFile(filename, nil)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
		return nil
	}

	config.CreateFromFiles("testing", f)

	prog, err := config.Load()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
		return nil
	}

	firstFunc, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok { // skip import decl if exists
		firstFunc = f.Decls[1].(*ast.FuncDecl) // panic here if no first func
	}
	cfg := cfg.FromFunc(firstFunc)
	v := make(map[int]ast.Stmt)
	stmts := make(map[ast.Stmt]int)
	objs := make(map[string]*types.Var)
	objNames := make(map[*types.Var]string)
	i := 1
	ast.Inspect(firstFunc, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if obj, ok := prog.Created[0].ObjectOf(x).(*types.Var); ok {
				objs[obj.Name()] = obj
				objNames[obj] = obj.Name()
			}
		case ast.Stmt:
			switch x.(type) {
			case *ast.BlockStmt:
				return true
			}
			v[i] = x
			stmts[x] = i
			i++
		case *ast.FuncLit:
			// skip statements in anonymous functions
			return false
		}
		return true
	})
	v[END] = cfg.Exit
	v[START] = cfg.Entry
	stmts[cfg.Entry] = START
	stmts[cfg.Exit] = END
	if len(v) != len(cfg.Blocks()) {
		t.Logf("expected %d vertices, got %d --construction error", len(v), len(cfg.Blocks()))
	}

	return &CFGWrapper{
		cfg:      cfg,
		prog:     prog,
		exp:      v,
		stmts:    stmts,
		objs:     objs,
		objNames: objNames,
		fset:     prog.Fset,
		f:        f,
	}
}

//prints given AST
func (c *CFGWrapper) printAST() {
	ast.Print(c.fset, c.f)
}
