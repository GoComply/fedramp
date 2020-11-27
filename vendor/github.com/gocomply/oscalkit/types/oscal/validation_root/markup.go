package validation_root

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Markup ...
type Markup struct {
	Raw       string `xml:",innerxml" yaml:"raw,omitempty"`
	PlainText string `xml:"-"`
}

// Short-hand method to build "markup-line" text
func ML(plain string) *Markup {
	plain = strings.ReplaceAll(plain, "&", "&amp;")
	plain = strings.ReplaceAll(plain, "<", "&lt;")
	plain = strings.ReplaceAll(plain, "<", "&gt;")
	return &Markup{
		Raw: plain,
	}
}

// Short-hand method to build "markup-multiline" text
func MML(plain string) *Markup {
	plain = strings.ReplaceAll(plain, "&", "&amp;")
	plain = strings.ReplaceAll(plain, "<", "&lt;")
	plain = strings.ReplaceAll(plain, "<", "&gt;")
	return &Markup{
		Raw: "<p>" + plain + "</p>",
	}
}

// PlainText representation
func (m *Markup) PlainString() string {
	if m.PlainText != "" {
		return m.PlainText
	}
	return m.xmlToPlain()
}

func (m *Markup) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &m.PlainText); err != nil {
		return err
	}
	return nil
}

func (m *Markup) MarshalJSON() ([]byte, error) {
	plain := m.PlainString()
	if strings.Contains(plain, "\n") {
		plain = strings.ReplaceAll(plain, "\n", "\\n")
	}
	if strings.Contains(plain, "\t") {
		plain = strings.ReplaceAll(plain, "\t", "\\t")
	}
	if strings.Contains(plain, "\"") {
		plain = strings.ReplaceAll(plain, "\"", "\\\"")
	}
	return []byte("\"" + plain + "\""), nil
}

func (m *Markup) xmlToPlain() string {
	re := regexp.MustCompile(`</p>`)
	s := re.ReplaceAllString(m.Raw, "")

	re = regexp.MustCompile(`<p>`)
	s = re.ReplaceAllString(s, "\t")
	return s
}
