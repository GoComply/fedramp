package docx_helper

import (
	"github.com/jbowtie/gokogiri/xml"
	"strings"
)

const libxml2_copy_constant = 2

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

	text = strings.TrimSpace(text)
	return textCell.SetContent(text)
}

func ParagraphReplaceWithText(originalParagraph xml.Node, newText string) error {
	for _, text := range strings.Split(newText, "\n\n") {
		clone := originalParagraph.Duplicate(libxml2_copy_constant)
		if err := ParagraphSetText(clone, text); err != nil {
			return err
		}
		if err := originalParagraph.Parent().AddChild(clone); err != nil {
			return err
		}
	}

	originalParagraph.Remove()
	return nil

}
