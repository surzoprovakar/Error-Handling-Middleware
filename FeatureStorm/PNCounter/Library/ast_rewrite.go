// Apply basic AST rewriting to Go code, without using external packages.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	//"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", os.Stdin, 0)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			id, ok := x.Fun.(*ast.Ident)
			if ok {
				if id.Name == "SetVal" {
					id.Name = "setValueWrap"
				}
				if id.Name == "NewCounter" {
					id.Name = " createWrap"
				}
				// if id.Name == "SetRemoteVal" {
				// 	id.Name = "setRemoteWrap"
				// }
			}

			s, ok := x.Fun.(*ast.SelectorExpr)
			if ok {
				if s.Sel.Name == "SetVal" {
					s.Sel.Name = "setValueWrap"
				}
				if s.Sel.Name == "SetRemoteVal" {
					s.Sel.Name = "setRemoteWrap"
				}
				if s.Sel.Name == "Merge" {
					s.Sel.Name = "mergeWrap"
				}
			}
			// case *ast.FuncDecl:
			// 	if x.Name.Name == "SetVal" {
			// 		x.Name.Name = "setValueWrap"
			// 	}
			// 	if x.Name.Name == "SetRemoteVal" {
			// 		x.Name.Name = "setRemoteWrap"
			// 	}

			//name := x.Fun.(*ast.Ident).Name
			//fmt.Println("Alt name is", name)
			/*
				case *ast.FuncDecl:
					if x.Name.Name == "setValue" {
						//keep it the same
						//x.Name.Name = "setValueWrap"
					}
					newCallStmt := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fmt",
								},
								Sel: &ast.Ident{
									Name: "Println",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"instrumentation"`,
								},
							},
						},
					}
					x.Body.List = append([]ast.Stmt{newCallStmt}, x.Body.List...)
			*/
		}

		return true
	})

	//fmt.Println("Modified AST:")
	printer.Fprint(os.Stdout, fset, file)
}
