package templater

import (
	"github.com/GoComply/fedramp/bundled"
	"github.com/GoComply/fedramp/pkg/fedramp"
	"github.com/GoComply/fedramp/pkg/fedramp/common"
	"github.com/docker/oscalkit/pkg/oscal_source"
)

func Convert(sspSource *oscal_source.OSCALSource, template string) error {
	_, err := fedramp.NewSSP(sspSource)
	if err != nil {
		return err
	}

	docx, err := bundled.TemplateDOCX(common.LevelModerate)
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
