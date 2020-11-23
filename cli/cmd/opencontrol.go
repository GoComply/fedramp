package cmd

import (
	"github.com/gocomply/fedramp/pkg/oc2oscal"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/urfave/cli"
)

var format string

// ConvertOpenControl ...
var openControl = cli.Command{
	Name:        "opencontrol",
	Usage:       `Convert OpenControl masonry repo into FedRAMP formatted OSCAL`,
	Description: `Convert OpenControl masonry repository into FedRAMP formatted OSCAL SSP Documents`,
	ArgsUsage:   "[masonry-repository] [output-directory]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "format, f",
			Usage:       "Format of the output: xml, json, or yaml",
			Value:       "xml",
			Destination: &format,
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() != 2 {
			return cli.NewExitError("Missing masonry repository or output directory", 1)
		}
		if constants.NewDocumentFormat(format) == constants.UnknownFormat {
			return cli.NewExitError("Unrecognized file format: "+format, 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		err := oc2oscal.Convert(c.Args()[0], c.Args()[1], constants.NewDocumentFormat(format))
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	},
}
