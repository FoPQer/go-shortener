// Command reset generates Reset methods for marked structs in the module.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

const marker = "generate:reset"

// main resolves the module root and runs the generator.
func main() {
	root, err := findModuleRoot()
	if err != nil {
		fatal(err)
	}

	if err := generateResetMethods(root); err != nil {
		fatal(err)
	}
}

// fatal prints an error to stderr and exits with a non-zero status.
func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// findModuleRoot walks upward from the current directory until it finds go.mod.
func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

// generateResetMethods loads all packages in the module and writes reset.gen.go files.
func generateResetMethods(root string) error {
	packagesConfig := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Dir: root,
		Env: os.Environ(),
	}

	pkgs, err := packages.Load(packagesConfig, "./...")
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("no packages found under %s", root)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return fmt.Errorf("package loading failed")
	}

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].PkgPath < pkgs[j].PkgPath
	})

	for _, pkg := range pkgs {
		if err := generatePackageFile(pkg); err != nil {
			return err
		}
	}

	return nil
}

type structSpec struct {
	name   string
	fields []fieldSpec
}

type fieldSpec struct {
	expr string
	typ  types.Type
}

// generatePackageFile renders reset methods for a single package.
func generatePackageFile(pkg *packages.Package) error {
	structs := collectStructs(pkg)
	if len(structs) == 0 {
		return removeGeneratedFile(pkg)
	}

	sort.Slice(structs, func(i, j int) bool {
		return structs[i].name < structs[j].name
	})

	src, err := renderPackage(pkg, structs)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(packageDir(pkg), "reset.gen.go")
	if existing, err := os.ReadFile(outputPath); err == nil && bytes.Equal(existing, src) {
		return nil
	}

	return os.WriteFile(outputPath, src, 0o644)
}

// removeGeneratedFile deletes the generated file when no marked structs remain.
func removeGeneratedFile(pkg *packages.Package) error {
	outputPath := filepath.Join(packageDir(pkg), "reset.gen.go")
	if _, err := os.Stat(outputPath); err == nil {
		return os.Remove(outputPath)
	}
	return nil
}

// packageDir returns the filesystem directory for the package.
func packageDir(pkg *packages.Package) string {
	if len(pkg.GoFiles) > 0 {
		return filepath.Dir(pkg.GoFiles[0])
	}
	if len(pkg.Syntax) > 0 {
		return filepath.Dir(pkg.Fset.Position(pkg.Syntax[0].Pos()).Filename)
	}
	return ""
}

// collectStructs finds all struct declarations marked with generate:reset.
func collectStructs(pkg *packages.Package) []structSpec {
	structs := make([]structSpec, 0)

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !hasMarker(genDecl.Doc, typeSpec.Doc, typeSpec.Comment) {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				obj, ok := pkg.TypesInfo.Defs[typeSpec.Name]
				if !ok || obj == nil {
					continue
				}
				namedType, ok := obj.Type().(*types.Named)
				if !ok {
					continue
				}

				structs = append(structs, structSpec{
					name:   namedType.Obj().Name(),
					fields: collectFields(pkg, structType),
				})
			}
		}
	}

	return structs
}

// hasMarker reports whether any supplied comment group contains the marker text.
func hasMarker(groups ...*ast.CommentGroup) bool {
	for _, group := range groups {
		if group != nil && strings.Contains(group.Text(), marker) {
			return true
		}
	}
	return false
}

// collectFields converts a struct AST into resettable field specifications.
func collectFields(pkg *packages.Package, structType *ast.StructType) []fieldSpec {
	fields := make([]fieldSpec, 0)
	for _, field := range structType.Fields.List {
		fieldType := pkg.TypesInfo.TypeOf(field.Type)
		if fieldType == nil {
			continue
		}

		if len(field.Names) == 0 {
			fields = append(fields, fieldSpec{
				expr: anonymousFieldName(field.Type),
				typ:  fieldType,
			})
			continue
		}

		for _, name := range field.Names {
			if name.Name == "_" {
				continue
			}
			fields = append(fields, fieldSpec{
				expr: name.Name,
				typ:  fieldType,
			})
		}
	}

	return fields
}

// anonymousFieldName returns a usable selector name for an embedded field.
func anonymousFieldName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	case *ast.StarExpr:
		return anonymousFieldName(t.X)
	default:
		return ""
	}
}

// renderPackage builds the Go source for a package-level reset file.
func renderPackage(pkg *packages.Package, structs []structSpec) ([]byte, error) {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "// Code generated by cmd/reset; DO NOT EDIT.")
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "package %s\n\n", pkg.Name)

	for _, spec := range structs {
		fmt.Fprintf(&buf, "func (v *%s) Reset() {\n", spec.name)
		fmt.Fprintln(&buf, "\tif v == nil {")
		fmt.Fprintln(&buf, "\t\treturn")
		fmt.Fprintln(&buf, "\t}")

		for _, field := range spec.fields {
			writeResetStatement(&buf, "v."+field.expr, field.typ, 1)
		}

		fmt.Fprintln(&buf, "}")
		fmt.Fprintln(&buf)
	}

	processed, err := imports.Process(filepath.Join(packageDir(pkg), "reset.gen.go"), buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return processed, nil
}

// writeResetStatement emits the reset logic for one field expression.
func writeResetStatement(buf *bytes.Buffer, expr string, typ types.Type, indentLevel int) {
	indent := strings.Repeat("\t", indentLevel)
	if typ == nil {
		return
	}

	switch t := typ.Underlying().(type) {
	case *types.Pointer:
		fmt.Fprintf(buf, "%sif %s != nil {\n", indent, expr)
		if hasResetMethod(t.Elem()) {
			fmt.Fprintf(buf, "%s\t%s.Reset()\n", indent, expr)
		} else {
			writeResetStatement(buf, derefExpr(expr), t.Elem(), indentLevel+1)
		}
		fmt.Fprintf(buf, "%s}\n", indent)
	case *types.Slice:
		fmt.Fprintf(buf, "%s%s = %s[:0]\n", indent, expr, expr)
	case *types.Map:
		fmt.Fprintf(buf, "%sclear(%s)\n", indent, expr)
	case *types.Struct:
		if hasResetMethod(typ) {
			fmt.Fprintf(buf, "%s(%s).Reset()\n", indent, expr)
			return
		}
		fmt.Fprintf(buf, "%s%s = %s\n", indent, expr, zeroExpr(typ))
	default:
		fmt.Fprintf(buf, "%s%s = %s\n", indent, expr, zeroExpr(typ))
	}
}

// derefExpr returns a parenthesized pointer dereference expression.
func derefExpr(expr string) string {
	return "(*" + expr + ")"
}

// hasResetMethod reports whether typ or *typ has a Reset method.
func hasResetMethod(typ types.Type) bool {
	if typ == nil {
		return false
	}

	if methodSetHas(typ, "Reset") {
		return true
	}

	return methodSetHas(types.NewPointer(typ), "Reset")
}

// methodSetHas checks whether typ exposes a method with the requested name.
func methodSetHas(typ types.Type, name string) bool {
	methods := types.NewMethodSet(typ)
	for i := 0; i < methods.Len(); i++ {
		if methods.At(i).Obj().Name() == name {
			return true
		}
	}
	return false
}

// typeString formats a type name for generated source.
func typeString(typ types.Type) string {
	return types.TypeString(typ, func(other *types.Package) string {
		if other == nil {
			return ""
		}
		return other.Name()
	})
}

// zeroExpr returns a Go expression representing the zero value for typ.
func zeroExpr(typ types.Type) string {
	if typ == nil {
		return "nil"
	}

	if named, ok := typ.(*types.Named); ok {
		return zeroExpr(named.Underlying())
	}

	switch t := typ.Underlying().(type) {
	case *types.Basic:
		return basicZeroExpr(t)
	case *types.Struct, *types.Array:
		return typeString(typ) + "{}"
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Interface, *types.Signature:
		return "nil"
	default:
		return typeString(typ) + "{}"
	}
}

// basicZeroExpr returns the literal zero value for a basic type.
func basicZeroExpr(basic *types.Basic) string {
	switch {
	case basic.Info()&types.IsBoolean != 0:
		return "false"
	case basic.Info()&types.IsString != 0:
		return `""`
	case basic.Info()&(types.IsInteger|types.IsFloat|types.IsComplex) != 0:
		return "0"
	default:
		return "nil"
	}
}
