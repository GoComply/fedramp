package docx_helper

import (
	"github.com/jbowtie/gokogiri/xml"
)

func ParagraphSetText(pNode xml.Node, text string) error {
	existingR, err := pNode.Search(".//w:r")
	if err != nil {
		return err
	}
	for _, rNode := range existingR {
		rNode.Remove()
	}
	err = pNode.AddChild(`<w:r><w:t></w:t></w:r>`)
	if err != nil {
		return err
	}
	textCell := pNode.LastChild().FirstChild()
	return textCell.SetContent(text)
}
