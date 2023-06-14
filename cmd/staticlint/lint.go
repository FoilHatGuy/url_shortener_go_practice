package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitFromMainAnalyzer = &analysis.Analyzer{
	Name: "exitFromMain",
	Doc:  "check for the cas of os.Exit in main()of package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			fmt.Printf("%+v\n", node)

			//nolint:gocritic
			switch x := node.(type) {
			case *ast.ExprStmt: // выражение
				call, ok := x.X.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				pkg, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if pkg.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(pkg.Pos(), "calling os.Exit() from main() of package main is forbidden")
				}
			}
			return true
		})
	}
	return nil, nil
}
