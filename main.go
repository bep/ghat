package main

import (
	"fmt"
	"ghat/src/core"
	"ghat/src/version"
	"os"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var file string

	var directory string

	app := &cli.App{
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:      "version",
				Aliases:   []string{"v"},
				Usage:     "Outputs the application version",
				UsageText: "ghat version",
				Action: func(*cli.Context) error {
					fmt.Println(version.Version)

					return nil
				},
			},
			{
				Name:      "swot",
				Aliases:   []string{"a"},
				Usage:     "updates GHA in a directory",
				UsageText: "ghat swot",
				Action: func(*cli.Context) error {

					if &file == nil {
						err := core.UpdateFile(&file)
						if err != nil {
							return err
						}
					} else {
						_, err := core.Files(&directory)
						if err != nil {
							return err
						}
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Usage:       "GHA file to parse",
						Destination: &file,
						Category:    "files",
					},
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Destination to update GHAs",
						Value:       ".",
						Destination: &directory,
						Category:    "files",
					},
				},
			},
		},
		Name:     "ghat",
		Usage:    "Update GHA dependencies",
		Compiled: time.Time{},
		Authors:  []*cli.Author{{Name: "James Woolfenden", Email: "jim.wolf@duck.com"}},
		Version:  version.Version,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("ghat failure")
	}
}

//for each
//	find all instances in file that match regex
//	uses:(.*)
//	replace with newest version
//save updated
