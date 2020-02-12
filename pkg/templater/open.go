package templater

import (
	"github.com/GoComply/fedramp/pkg/fedramp"
	"github.com/GoComply/fedramp/pkg/templater/template"
	"github.com/docker/oscalkit/pkg/oscal_source"
)

func Convert(sspSource *oscal_source.OSCALSource, tmplt string) error {
	plan, err := fedramp.NewSSP(sspSource)
	if err != nil {
		return err
	}

	_, err = template.NewTemplate(plan.Level())
	if err != nil {
		return err
	}

	return nil
}

func ConvertFile(oscalSSPFilePath, template, outputPath string) error {
	source, err := oscal_source.Open(oscalSSPFilePath)
	if err != nil {
		return nil
	}
	defer source.Close()
	return Convert(source, template)
}
