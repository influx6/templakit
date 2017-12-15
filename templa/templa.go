package templa

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

// StructGenerator defines a struct level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// It provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @makeFor
//
func StructGenerator(toDir string, an ast.AnnotationDeclaration, ty ast.StructDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation  ast.AnnotationDeclaration
		PkgDeclr    ast.PackageDeclaration
		Package     ast.Package
		StructDeclr ast.StructDeclaration
	}{
		PkgDeclr:    pkgDeclr,
		Annotation:  an,
		Package:     pkg,
		StructDeclr: ty,
	})
}

// InterfaceGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// It provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @makeFor
//
func InterfaceGenerator(toDir string, an ast.AnnotationDeclaration, ty ast.InterfaceDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation     ast.AnnotationDeclaration
		PkgDeclr       ast.PackageDeclaration
		Package        ast.Package
		InterfaceDeclr ast.InterfaceDeclaration
	}{
		PkgDeclr:       pkgDeclr,
		Annotation:     an,
		Package:        pkg,
		InterfaceDeclr: ty,
	})
}

// PackageGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// It provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @makeFor
//
func PackageGenerator(toDir string, an ast.AnnotationDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation ast.AnnotationDeclaration
		PkgDeclr   ast.PackageDeclaration
		Package    ast.Package
	}{
		PkgDeclr:   pkgDeclr,
		Annotation: an,
		Package:    pkg,
	})
}

// AnyTypeGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// It provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @makeFor
//
func AnyTypeGenerator(toDir string, an ast.AnnotationDeclaration, ty ast.TypeDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation ast.AnnotationDeclaration
		TypeDeclr  ast.TypeDeclaration
		Package    ast.Package
		PkgDeclr   ast.PackageDeclaration
	}{
		Annotation: an,
		TypeDeclr:  ty,
		PkgDeclr:   pkgDeclr,
		Package:    pkg,
	})
}

func handleGeneration(toDir string, an ast.AnnotationDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package, binding interface{}) ([]gen.WriteDirective, error) {
	templaterID, ok := an.Params["id"]
	if !ok {
		return nil, errors.New("No source id provided")
	}

	// Get all templaters AnnotationDeclaration.
	templaters := pkg.AnnotationsFor("source")
	if len(templaters) == 0 {
		return nil, errors.New("No  found in package")
	}

	var target ast.AnnotationDeclaration

	// Search for source with associated ID, if not found, return error, if multiple found, use the first.
	for _, target = range templaters {
		if target.Params["id"] != templaterID {
			continue
		}

		break
	}

	var templateData string

	switch len(target.Template) == 0 {
	case true:
		templateFilePath, dok := target.Params["file"]
		if !dok && target.Template == "" {
			return nil, errors.New("Expected Template from annotation or provide `file => 'path_to_template`")
		}

		baseDir := filepath.Dir(pkgDeclr.FilePath)
		templateFile := filepath.Join(baseDir, templateFilePath)

		data, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to find template file: %+q", err)
		}

		templateData = string(data)
	case false:
		templateData = target.Template
	}

	var directives []gen.WriteDirective

	genName := strings.ToLower(target.Params["gen"])

	fileName, ok := an.Params["filename"]
	if !ok {
		fileName = fmt.Sprintf("%s_impl_gen.go", strings.ToLower(an.Name))
	}

	typeGen := gen.Block(gen.SourceTextWith(templateData, template.FuncMap{
		"sel":                    an.Param,
		"params":                 an.Param,
		"attrs":                  an.Attr,
		"hasArg":                 an.HasArg,
		"annotationDefer":        func() bool { return an.Defer },
		"annotationTemplate":     func() string { return an.Template },
		"annotationParams":       func() map[string]string { return an.Params },
		"annotationAttrs":        func() map[string]interface{} { return an.Attrs },
		"annotationArguments":    func() []string { return an.Arguments },
		"targetSel":              target.Param,
		"targetParams":           target.Param,
		"targetAttrs":            target.Attr,
		"targetHasArg":           target.HasArg,
		"targetDefer":            func() bool { return target.Defer },
		"targetTemplate":         func() string { return target.Template },
		"targetArguments":        func() []string { return target.Arguments },
		"targetAnnotationParams": func() map[string]string { return target.Params },
		"targetAnnotationAttrs":  func() map[string]interface{} { return target.Attrs },
	}, binding))

	switch genName {
	case "partial_test.go":

		var packageName string

		switch len(an.Params["packageName"]) == 0 {
		case true:
			packageName = ast.WhichPackage(toDir, pkg)
		case false:
			packageName = target.Params["packageName"]
		}

		packageName = fmt.Sprintf("%s_test", packageName)

		pkgGen := gen.Block(

			gen.Package(
				gen.Name(packageName),
				typeGen,
			),
		)

		directives = append(directives, gen.WriteDirective{
			FileName:     fileName,
			DontOverride: true,
			Writer:       fmtwriter.New(pkgGen, true, true),
		})

	case "partial.go":

		var packageName string

		switch len(an.Params["packageName"]) == 0 {
		case true:
			packageName = ast.WhichPackage(toDir, pkg)
		case false:
			packageName = target.Params["packageName"]
		}

		pkgGen := gen.Block(

			gen.Package(
				gen.Name(packageName),
				typeGen,
			),
		)

		directives = append(directives, gen.WriteDirective{
			FileName:     fileName,
			DontOverride: true,
			Writer:       fmtwriter.New(pkgGen, true, true),
		})

	case "go":
		directives = append(directives, gen.WriteDirective{
			FileName:     fileName,
			DontOverride: true,

			Writer: fmtwriter.New(typeGen, true, true),
		})

	default:
		directives = append(directives, gen.WriteDirective{
			Writer:       typeGen,
			DontOverride: true,
			FileName:     fileName,
		})
	}

	return directives, nil
}
