package validation_root

import (
	"strings"
)

// Markup ...
type Markup struct {
	Raw string `xml:",innerxml" json:"raw,omitempty" yaml:"raw,omitempty"`
}

func MarkupFromPlain(plain string) *Markup {
	plain = strings.ReplaceAll(plain, "&", "&amp;")
	plain = strings.ReplaceAll(plain, "<", "&lt;")
	plain = strings.ReplaceAll(plain, "<", "&gt;")
	return &Markup{
		Raw: "<p>" + plain + "</p>",
	}
}
