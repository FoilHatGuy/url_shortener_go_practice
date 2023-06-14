package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func choice(node ast.Node) string {
	switch x := node.(type) {
	// любой идентификатор, например имя переменной
	case *ast.Ident:
		return "Ident(" + x.Name + ")"
	// литерал, например 100500 или "Hello"
	case *ast.BasicLit:
		return "BasicLit(" + x.Value + ")"
	// выражение с левой и правой частью, разделёнными токеном, например err != nil
	case *ast.BinaryExpr:
		return "BinaryExpr(" + x.Op.String() + ")"
	// унарное выражение, например -foo
	case *ast.UnaryExpr:
		return "UnaryExpr(" + x.Op.String() + ")"
	// foo++
	case *ast.IncDecStmt:
		return "IncDecStmt(" + x.Tok.String() + ")"
	// присваивание или короткая форма декларации foo := 5
	case *ast.AssignStmt:
		return "AssignStmt(" + x.Tok.String() + ")"
	default:
		return strings.ReplaceAll(fmt.Sprintf("%T", x), "*ast.", "")
	}
}

func main() {
	fmt.Println(`digraph{graph [dpi=128];`)
	fset := token.NewFileSet()
	flag.Parse()
	f, err := parser.ParseFile(fset, flag.Arg(0), nil, 0)
	if err != nil {
		fmt.Printf("Failed to parse file\n")
		return
	}
	parents := []*ast.Node{}
	fmt.Println(`n0x0[label="ROOT"]`)
	ast.Inspect(f, func(node ast.Node) bool {
		if node != nil {
			var p *ast.Node
			if len(parents) != 0 {
				p = parents[len(parents)-1]
			}
			fmt.Printf("n%p->n%p\n", p, &node)
			fmt.Printf(`n%p[ label=%q ]`+"\n", &node, choice(node))
		}
		if node == nil {
			// возвращаемся к родительскому элементу
			// уменьшаем стек
			parents = parents[:len(parents)-1]
		} else {
			// добавляем в стек узел node
			parents = append(parents, &node)
		}
		return true
	})
	fmt.Println(`}`)
}
