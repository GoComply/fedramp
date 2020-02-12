package templater

import (
	"github.com/GoComply/fedramp/bundled"
	"github.com/GoComply/fedramp/pkg/fedramp"
	"github.com/docker/oscalkit/pkg/oscal_source"
)

func Convert(sspSource *oscal_source.OSCALSource, template string) error {
	plan, err := fedramp.NewSSP(sspSource)
	if err != nil {
		return err
	}

	docx, err := bundled.TemplateDOCX(plan.Level())
	if err != nil {
		return err
	}
	defer docx.Close()

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
