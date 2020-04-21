package docx_helper

import (
	"github.com/jbowtie/gokogiri/xml"
	"strings"
)

// ConcatTextNodes will concatenate the text from an array of text nodes and trim any whitespace from the final result.
func ConcatTextNodes(textNodes []xml.Node) string {
	result := ""
	for _, textNode := range textNodes {
		result += textNode.Content()
	}
	return strings.TrimSpace(result)
}
