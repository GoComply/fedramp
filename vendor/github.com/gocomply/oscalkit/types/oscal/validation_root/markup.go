package validation_root

import (
	"encoding/json"
	"strings"
)

// Markup ...
type Markup struct {
	Raw       string `xml:",innerxml" yaml:"raw,omitempty"`
	PlainText string `xml:"-"`
}

func MarkupFromPlain(plain string) *Markup {
	plain = strings.ReplaceAll(plain, "&", "&amp;")
	plain = strings.ReplaceAll(plain, "<", "&lt;")
	plain = strings.ReplaceAll(plain, "<", "&gt;")
	return &Markup{
		Raw: "<p>" + plain + "</p>",
	}
}

func (m *Markup) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &m.PlainText); err != nil {
		return err
	}
	return nil
}
