package templater

import (
	"fmt"
	"github.com/GoComply/fedramp/pkg/fedramp"
	"github.com/GoComply/fedramp/pkg/templater/template"
	"github.com/docker/oscalkit/pkg/oscal_source"
)

func Convert(sspSource *oscal_source.OSCALSource, outputPath string) error {
	plan, err := fedramp.NewSSP(sspSource)
	if err != nil {
		return err
	}

	doc, err := template.NewTemplate(plan.Level())
	if err != nil {
		return err
	}

	err = fillInSSP(doc, plan)
	if err != nil {
		return err
	}
	return doc.Save(outputPath)
}

func ConvertFile(oscalSSPFilePath, template, outputPath string) error {
	source, err := oscal_source.Open(oscalSSPFilePath)
	if err != nil {
		return nil
	}
	defer source.Close()
	return Convert(source, outputPath)
}

func fillInSSP(doc *template.Template, plan *fedramp.SSP) error {
	tables, err := doc.ControlSummaryInformations()
	if err != nil {
		return err
	}

	for _, table := range tables {
		name, err := table.ControlName()
		if err != nil {
			return err
		}
		fmt.Println(name)
	}
	return nil
}
