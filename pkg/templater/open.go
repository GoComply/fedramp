package templater

import (
	"fmt"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/gocomply/fedramp/pkg/templater/template"
	"github.com/gocomply/oscalkit/pkg/oscal_source"
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
		return err
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
		controlId, err := table.ControlName()
		if err != nil {
			return err
		}
		responsibleRole, err := table.ResponsibleRole()
		if err != nil {
			return err
		}
		// Implements: 5.2. Responsible Roles and Parameter Assignments
		if err = responsibleRole.SetValue(plan.ResponsibleRoleForControl(controlId)); err != nil {
			return err
		}

		paramRows, err := table.ParameterRows()
		if err != nil {
			return err
		}

		for idx, paramRow := range paramRows {
			paramId, err := paramRow.ControlId()
			if err != nil {
				return fmt.Errorf("%v while trying to parse parameter rows in '%s Control Summary Information' table", err, controlId)
			}
			newValue, err := plan.ParamValue(paramId, idx+1)
			if err != nil {
				return err
			}
			err = paramRow.SetValue(newValue)
			if err != nil {
				return err
			}
		}

		// Implements: 5.3 Implementation Status
		implStatus, err := table.ImplementationStatus()
		if err != nil {
			return err
		}
		if err = implStatus.SetValue(plan.ImplementationStatusForControl(controlId)); err != nil {
			return err
		}

		// TODO: 5.4 Control Origination
	}

	// Implements: 5.5. Control Implementation Descriptions
	cidTables, err := doc.ControlImplementationDescriptions()
	if err != nil {
		return err
	}

	for _, table := range cidTables {
		controlId, err := table.ControlName()
		if err != nil {
			return err
		}
		partRows, err := table.PartRows()
		if err != nil {
			return err
		}
		for _, partRow := range partRows {
			partName, err := partRow.PartName()
			if err != nil {
				return err
			}
			statementText, err := plan.StatementTextForPart(controlId, partName)
			if err != nil {
				return err
			}
			if statementText != "" {
				if err = partRow.SetValue(statementText); err != nil {
					return err
				}
			}
		}
		plain, err := table.Plain()
		if err != nil {
			return err
		}
		if plain {
			txt, err := plan.StatementTextFor(controlId)
			if err != nil {
				return err
			}
			if txt != "" {
				if err = table.SetValue(txt); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
