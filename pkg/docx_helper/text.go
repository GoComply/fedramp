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
