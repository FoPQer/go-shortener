package main

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports use of built-in panic and calls to os.Exit / log.Fatal(f|ln)
// outside the main function of the main package.
var Analyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc: "reports use of built-in panic and os.Exit/log.Fatal(f|ln) " +
		"called outside the main function of the main package",
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	isMainPkg := pass.Pkg.Name() == "main"

	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Body == nil {
				continue
			}

			inMainFunc := isMainPkg &&
				funcDecl.Recv == nil &&
				funcDecl.Name.Name == "main"

			inspectBody(pass, funcDecl.Body, inMainFunc)
		}
	}

	return nil, nil
}

// inspectBody walks a function body and reports forbidden patterns.
// inMain is true only when the body belongs to the top-level main() of package main.
func inspectBody(pass *analysis.Pass, body ast.Node, inMain bool) {
	ast.Inspect(body, func(n ast.Node) bool {
		// Handle closures: they are never "in main".
		if lit, ok := n.(*ast.FuncLit); ok {
			inspectBody(pass, lit.Body, false)
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		reportPanic(pass, call)
		if !inMain {
			reportForbiddenExit(pass, call)
		}

		return true
	})
}

// reportPanic reports use of the built-in panic function.
func reportPanic(pass *analysis.Pass, call *ast.CallExpr) {
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "panic" {
		return
	}
	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return
	}
	if _, ok := obj.(*types.Builtin); ok {
		pass.Reportf(call.Pos(), "use of built-in panic is forbidden")
	}
}

// reportForbiddenExit reports os.Exit and log.Fatal/Fatalf/Fatalln calls.
func reportForbiddenExit(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	fn := pass.TypesInfo.Uses[sel.Sel]
	if fn == nil || fn.Pkg() == nil {
		return
	}

	pkg := fn.Pkg().Path()
	name := fn.Name()

	switch {
	case pkg == "os" && name == "Exit":
		pass.Reportf(call.Pos(), "os.Exit called outside main function of main package")
	case pkg == "log" && (name == "Fatal" || name == "Fatalf" || name == "Fatalln"):
		pass.Reportf(call.Pos(), "log.%s called outside main function of main package", name)
	}
}
