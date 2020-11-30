package cmd

import (
	"github.com/gocomply/fedramp/pkg/templater"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

// Execute ...
func Execute() error {
	app := cli.NewApp()
	app.Name = "gocomply_fedramp"
	app.Usage = "OSCAL-FedRAMP Workbench"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug command output",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}
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
	Before: func(c *cli.Context) error {
		if c.NArg() != 2 {
			return cli.NewExitError("Exactly 2 arguments are required", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		sspFile, outputFile := c.Args()[0], c.Args()[1]
		err := templater.ConvertFile(sspFile, outputFile)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}
