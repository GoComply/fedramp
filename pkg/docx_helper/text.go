package docx_helper

import (
	"github.com/jbowtie/gokogiri/xml"
	"strings"
)

// ConcatTextNodesList will concatenate the text from an array of text nodes and trim any whitespace from the final result.
func ConcatTextNodesList(textNodes []xml.Node) string {
	result := ""
	for _, textNode := range textNodes {
		result += textNode.Content()
	}
	return strings.TrimSpace(result)
}

// ConcatTextNodes will find all <w:t/> child nodes and concatenate the text of those and trim any whitespaces from the final result.
func ConcatTextNodes(node xml.Node) (string, error) {
	textNodes, err := node.Search(".//w:t")
	if err != nil {
		return "", err
	}
	return ConcatTextNodesList(textNodes), nil
}
