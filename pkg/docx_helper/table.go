package docx_helper

import (
	"errors"
	"fmt"
	"github.com/jbowtie/gokogiri/xml"
	"strings"
)

const libxml2_copy_constant = 2

func ParseTableHeaderContent(table xml.Node) (content string, err error) {
	nodes, err := table.Search(".//w:tr")
	if err != nil {
		return
	}
	if len(nodes) == 0 {
		err = errors.New("Could not locate table header: no w:tr elements found")
		return
	}
	// we only care about the first match
	content = nodes[0].Content()

	return
}

func RowReplaceText(row xml.Node, newText string) error {
	paragraphNodes, err := row.Search(".//w:p")
	if err != nil {
		return err
	}
	if len(paragraphNodes) != 1 {
		return fmt.Errorf("Could not update a table row: found '%d' paragraph(s) in while expecting only 1.", len(paragraphNodes))
	}
	templateParagraph := paragraphNodes[0]

	for _, text := range strings.Split(newText, "\n\n") {
		clone := templateParagraph.Duplicate(libxml2_copy_constant)
		err = ParagraphSetText(clone, text)
		if err != nil {
			return err
		}
		err = templateParagraph.Parent().AddChild(clone)
		if err != nil {
			return err
		}
	}

	templateParagraph.Remove()
	return nil
}
