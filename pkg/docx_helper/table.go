package docx_helper

import (
	"errors"
	"github.com/jbowtie/gokogiri/xml"
)

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
