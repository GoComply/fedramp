package template

import (
	"fmt"
	"github.com/gocomply/fedramp/pkg/docx_helper"
	"github.com/jbowtie/gokogiri/xml"
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
	rows, err := cid.node.Search(".//w:tc[starts-with(normalize-space(.), 'Part')]")
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
