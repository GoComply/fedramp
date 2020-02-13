package template

import (
	"errors"
	"github.com/jbowtie/gokogiri/xml"
	"regexp"
)

// ControlSummaryInformation represents single table labeled "Control Summary Information"
type ControlSummaryInformation struct {
	table xml.Node
}

func (t *Template) ControlSummaryInformations() ([]ControlSummaryInformation, error) {
	tables, err := t.querySummaryTables()
	if err != nil {
		return nil, err
	}

	var result []ControlSummaryInformation
	for _, table := range tables {
		result = append(result, ControlSummaryInformation{
			table: table,
		})
	}
	return result, nil
}

func (t *Template) querySummaryTables() ([]xml.Node, error) {
	return t.xmlDoc.Search(
		"//w:tbl[contains(normalize-space(.), 'Control Summary') or contains(normalize-space(.), 'Control Enhancement Summary')]",
	)
}

func (csi *ControlSummaryInformation) ControlName() (name string, err error) {
	content, err := csi.queryTableHeader()
	if err != nil {
		return
	}

	if content == "CM2 (7)Control Summary Information" {
		// Workaround typo in the 8/28/2018 version of FedRAMP-SSP-High-Baseline-Template.docx
		content = "CM-2 (7)Control Summary Information"
	}

	// matches controls and control enhancements, e.g. `AC-2`, `AC-2 (1)`, etc.
	regex := regexp.MustCompile(`[A-Z]{2}-\d+( +\(\d+\))?`)
	name = regex.FindString(content)
	if name == "" {
		err = errors.New("control name not found for " + content)
	}
	return
}

func (csi *ControlSummaryInformation) queryTableHeader() (content string, err error) {
	nodes, err := csi.table.Search(".//w:tr")
	if err != nil {
		return
	}
	if len(nodes) == 0 {
		err = errors.New("Could not locate control name in the heading of 'Control Summary Information' table")
		return
	}
	// we only care about the first match
	content = nodes[0].Content()

	return
}
