package cmd

import (
	"github.com/GoComply/fedramp/pkg/templater"
	"github.com/urfave/cli"
	"os"
)

// Execute ...
func Execute() error {
	app := cli.NewApp()
	app.Name = "fedramp"
	app.Usage = "OSCAL-FedRAMP Workbench"
	app.Commands = []cli.Command{
		convert,
		openControl,
	}

	return app.Run(os.Args)
}

var convert = cli.Command{
	Name:      "convert",
	Usage:     "Convert OSCAL SSP to FedRAMP Document",
	ArgsUsage: "[ssp.oscal.xml] [output.docx]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "template, t",
			Usage:       "Fedramp docx template",
			Destination: &template,
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() != 2 {
			return cli.NewExitError("Exactly 2 arguments are required", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		sspFile, outputFile := c.Args()[0], c.Args()[1]
		return templater.ConvertFile(sspFile, template, outputFile)
	},
}

var template string
