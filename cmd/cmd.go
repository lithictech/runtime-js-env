package cmd

import (
	"github.com/lithictech/runtime-js-env/jsenv"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Execute() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "index",
				Aliases: shortf("i"),
				Value:   "index.html",
				Usage:   "Path to the index.html file. Default to index.html in pwd.",
			},
			&cli.StringFlag{
				Name:    "window-var-name",
				Aliases: shortf("w"),
				Value:   jsenv.DefaultConfig.WindowVarName,
				Usage:   "Attribute name for the config object.",
			},
			&cli.StringSliceFlag{
				Name:    "env-prefixes",
				Aliases: shortf("p"),
				Value:   cli.NewStringSlice(jsenv.DefaultConfig.EnvPrefixes...),
				Usage:   "Environment variable prefixes to copy into the config object.",
			},
			&cli.StringFlag{
				Name:    "indent",
				Aliases: shortf("t"),
				Value:   jsenv.DefaultConfig.Indent,
				Usage:   "Indentation for each line in the config script tag.",
			},
		},
		Action: func(c *cli.Context) error {
			opts := jsenv.Config{
				WindowVarName: c.String("window-var-name"),
				EnvPrefixes:   c.StringSlice("env-prefixes"),
				Indent:        c.String("indent"),
			}
			return jsenv.InstallAt(c.String("index"), opts)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func shortf(s string) []string {
	return []string{s}
}
