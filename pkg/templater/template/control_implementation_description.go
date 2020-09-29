package template

import (
	"errors"
	"fmt"
	"github.com/gocomply/fedramp/pkg/docx_helper"
	"github.com/jbowtie/gokogiri/xml"
	"regexp"
	"strings"
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

func (cid *ControlImplementationDescription) ReqRows() ([]PartRow, error) {
	rows, err := cid.node.Search(".//w:tr[starts-with(normalize-space(.), 'Req')]")
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

func (cid *ControlImplementationDescription) Plain() (bool, error) {
	p, err := cid.PartRows()
	if err != nil {
		return false, err
	}
	if len(p) > 0 {
		return false, nil
	}
	r, err := cid.ReqRows()
	if err != nil {
		return false, err
	}
	return len(r) == 0, nil
}

func (cid *ControlImplementationDescription) SetValue(response string) error {
	rows, err := cid.node.Search(".//w:tr")
	if err != nil {
		return err
	}
	if len(rows) != 2 {
		return fmt.Errorf("Could not update 'What is the solution and how is it implemented' table: found '%d' rows while expecting 2.", len(rows))
	}
	paragraphNodes, err := rows[1].Search(".//w:p")
	if err != nil {
		return err
	}
	if len(paragraphNodes) != 1 {
		return fmt.Errorf("Could not update 'What is the solution and how is it implemented' table: found '%d' paragraph(s) in the last row while expecting only 1.", len(paragraphNodes))
	}
	return docx_helper.ParagraphSetText(paragraphNodes[0], response)
}

func (pr *PartRow) PartName() (string, error) {
	tcNodes, err := pr.node.Search(".//w:tc")
	if err != nil {
		return "", err
	}
	if len(tcNodes) != 2 {
		return "", fmt.Errorf("Could not parse 'Part' row, expected 2 <w:tc/> elements but got %d; %s", len(tcNodes), pr.node)
	}

	nodes, err := tcNodes[0].Search(".//w:t[starts-with(normalize-space(.), 'Part')]")
	if err != nil {
		return "", err
	}
	if len(nodes) != 1 {
		return "", fmt.Errorf("Could not find Part text field in table named: 'What is the solution and how is it implemented?'")
	}
	txt, err := docx_helper.ConcatTextNodes(tcNodes[0])
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`^Part\s+([^:]*):*\s*$`)
	match := re.FindStringSubmatch(txt)
	if len(match) == 0 {
		return "", fmt.Errorf("Could not locate Part ID in text: '%s'", txt)
	}
	if strings.Contains(match[1], " ") || len(match[1]) > 2 {
		return "", fmt.Errorf("Suspicious '%s' Part ID found in the text '%s'", match[1], txt)
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
