package profile

import (
	"fmt"
	"github.com/gocomply/oscalkit/types/oscal/validation_root"
	"strings"
)

func (p *Profile) GetDocumentFragment(uri string) (*validation_root.Resource, error) {
	if strings.HasPrefix(uri, "#") {
		uri = uri[1:]
	}
	if p.BackMatter == nil {
		return nil, fmt.Errorf("cannot resolve %s within Profile %s, back-matter information missing", uri, p.Uuid)
	}
	return p.BackMatter.GetResourceByUuid(uri), nil
}
