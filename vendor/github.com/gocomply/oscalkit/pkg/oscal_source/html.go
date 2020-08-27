package oscal_source

import (
	"bytes"
	"errors"
	"github.com/gocomply/oscalkit/pkg/bundled"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/xslt"
)

func (s *OSCALSource) HTML() (*bytes.Buffer, error) {
	if s.OSCAL().DocumentType() != constants.CatalogDocument {
		return nil, errors.New("HTML is supported only for OSCAL Catalog")
	}
	transformation, err := bundled.HtmlXslt()
	if err != nil {
		return nil, err
	}
	defer transformation.Cleanup()

	return xslt.Transform(transformation.Path, s.UserPath)
}
