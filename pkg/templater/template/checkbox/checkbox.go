package checkbox

import (
	"fmt"
	"github.com/GoComply/fedramp/pkg/docx_helper"
	"github.com/jbowtie/gokogiri/xml"
)

type NotFound struct {
	err string
}

func (e *NotFound) Error() string {
	return e.err
}

type CheckBox struct {
	node      xml.Node
	checkMark xml.Node
	textNodes []xml.Node
}

func Parse(paragraph xml.Node) (*CheckBox, error) {
	checkBoxTag, err := findCheckBoxTag(paragraph)
	if err != nil {
		return nil, err
	}

	textNodes, err := paragraph.Search(".//w:t")
	if len(textNodes) < 1 || err != nil {
		return nil, fmt.Errorf("Could not find any <w:t/> elements after checkbox: %v", err)
	}

	return &CheckBox{node: paragraph, checkMark: checkBoxTag, textNodes: textNodes}, nil
}

func findCheckBoxTag(paragraph xml.Node) (xml.Node, error) {
	checkBoxes, err := paragraph.Search("(.//w:checkBox//w:default)|(.//w14:checkbox//w14:checked)")
	if err != nil {
		return nil, err
	} else if len(checkBoxes) != 1 {
		return nil, &NotFound{err: "No check box found"}
	}
	return checkBoxes[0], nil
}

func (cb *CheckBox) Text() string {
	return docx_helper.ConcatTextNodes(cb.textNodes)
}

const (
	boxChecked    = "☒"
	boxNotChecked = "☐"
)

func (cb *CheckBox) SetChecked() {
	cb.checkMark.AttributeList()[0].SetContent("1")
	if len(cb.textNodes) > 0 && cb.textNodes[0].Content() == boxNotChecked {
		cb.textNodes[0].SetContent(boxChecked)
	}
}
