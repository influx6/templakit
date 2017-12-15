package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gokit/templakit/templa"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/flags"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	"github.com/influx6/moz/ast"
)

func main() {
	flags.Run("templakit", flags.Command{
		Name:      "generate",
		ShortDesc: "Generates code from template source",
		Desc:      "Generates code from template source",
		Action: func(ctx context.Context) error {
			force, _ := ctx.Bag().GetBool("force")
			dest, _ := ctx.Bag().GetString("dest")
			verbose, _ := ctx.Bag().GetBool("verbose")

			logs := metrics.New()
			if verbose {
				logs = metrics.New(custom.StackDisplay(os.Stderr))
			}

			currentdir, err := os.Getwd()
			if err != nil {
				return err
			}

			generators := ast.NewAnnotationRegistryWith(logs)
			generators.Register("@source", templa.StructGenerator)
			generators.Register("@source", templa.InterfaceGenerator)
			generators.Register("@source", templa.PackageGenerator)
			generators.Register("@source", templa.AnyTypeGenerator)

			res, err := ast.ParseAnnotations(logs, currentdir)
			if err != nil {
				return err
			}

			return ast.SimplyParse(dest, logs, generators, force, res...)
		},
		Flags: []flags.Flag{
			&flags.BoolFlag{
				Name: "verbose",
				Desc: "verbose logs all operations out to console.",
			},
			&flags.BoolFlag{
				Name: "force",
				Desc: "force regeneration of packages annotation directives.",
			},
			&flags.StringFlag{
				Name:    "dest",
				Default: "./",
				Desc:    "relative destination for package",
				Validation: func(received string) error {
					if filepath.IsAbs(received) {
						return errors.New("only relative paths not absolute allowed")
					}
					return nil
				},
			},
		},
	})

}
