package profile

import (
	"errors"
	"net/url"
	"strings"
)

// IsHttpResource returns true if import refers to http resource
func (i *Import) IsHttpResource() bool {
	url, err := url.Parse(i.Href)
	if err != nil {
		return false
	}
	return strings.HasPrefix(url.Scheme, "http")
}

// Validates that profile import contains valid href
func (i *Import) ValidateHref() error {
	if i.Href == "" {
		return errors.New("href cannot be empty")
	}
	_, err := url.Parse(i.Href)
	if err != nil {
		return err
	}
	return nil
}
