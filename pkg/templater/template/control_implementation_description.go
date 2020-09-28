package template

import (
	"errors"
	"fmt"
	"github.com/gocomply/fedramp/pkg/docx_helper"
	"github.com/jbowtie/gokogiri/xml"
	"regexp"
)

// ControlImplementationDescription represents single table labeled "What is the solution and how is it implemented?"
type ControlImplementationDescription struct {
	node xml.Node
}

func (t *Template) ControlImplementationDescriptions() ([]ControlImplementationDescription, error) {
	tables, err := t.queryImplementationTables()
	if err != nil {
		return nil, err
	}
	var result []ControlImplementationDescription
	for _, table := range tables {
		result = append(result, ControlImplementationDescription{
			node: table,
		})
	}
	return result, nil
}

func (t *Template) queryImplementationTables() ([]xml.Node, error) {
	return t.xmlDoc.Search(
		"//w:tbl[contains(normalize-space(.), ' What is the solution and how is it implemented?')]",
	)
}

func (cid *ControlImplementationDescription) ControlName() (name string, err error) {
	content, err := cid.queryTableHeader()
	if err != nil {
		return
	}
	name, err = parseControlId(content)
	return
}

func (cid *ControlImplementationDescription) queryTableHeader() (string, error) {
	content, err := docx_helper.ParseTableHeaderContent(cid.node)
	if err != nil {
		return "", fmt.Errorf("Control Name in 'What is the solution and how is it implemented' table not found: %s", err)

	}
	return content, err
}

type PartRow struct {
	node xml.Node
}

func (cid *ControlImplementationDescription) PartRows() ([]PartRow, error) {
	rows, err := cid.node.Search(".//w:tr[starts-with(normalize-space(.), 'Part')]")
	if err != nil {
		return nil, err
	}
	var result []PartRow
	for _, row := range rows {
		result = append(result, PartRow{
			node: row,
		})
	}
	return result, nil
}

func (pr *PartRow) PartName() (string, error) {
	nodes, err := pr.node.Search(".//w:t[starts-with(normalize-space(.), 'Part')]")
	if err != nil {
		return "", err
	}
	if len(nodes) != 1 {
		return "", fmt.Errorf("Could not find Part text field in table named: 'What is the solution and how is it implemented?'")
	}
	txt, err := docx_helper.ConcatTextNodes(pr.node)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`^Part\s+([^:]*):*\s*$`)
	match := re.FindStringSubmatch(txt)
	if len(match) == 0 {
		return "", fmt.Errorf("Could not locate Part ID in text: '%s'", txt)
	}
	return match[1], nil
}

func (pr *PartRow) SetValue(partResponse string) error {
	paragraphNodes, err := pr.node.Search(".//w:p")
	if err != nil {
		return fmt.Errorf("Cannot search for paragraphs node within Part row: %s", err)
	}
	if len(paragraphNodes) != 2 {
		return errors.New("Cannot edit Part row, expected 2 paragraphs node")
	}
	return docx_helper.ParagraphSetText(paragraphNodes[1], partResponse)
}
