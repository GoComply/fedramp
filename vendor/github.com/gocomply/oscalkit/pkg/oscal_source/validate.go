package oscal_source

import (
	"errors"
	"github.com/gocomply/oscalkit/pkg/bundled"
	"github.com/gocomply/oscalkit/pkg/json_validation"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/xml_validation"
)

type validator func(schemaPath, inputFile string) error

func (s *OSCALSource) Validate() error {
	validate := s.relevantValidator()
	if validate == nil {
		return errors.New("No validator available this file type")
	}
	schema, err := s.relevantSchema()
	if err != nil {
		return err
	}
	defer schema.Cleanup()
	return validate(schema.Path, s.UserPath)
}

func (s *OSCALSource) relevantSchema() (*bundled.BundledFile, error) {
	return bundled.Schema(s.DocumentFormat(), s.OSCAL().DocumentType())
}

func (s *OSCALSource) relevantValidator() validator {
	switch s.DocumentFormat() {
	case constants.XmlFormat:
		return xml_validation.Validate
	case constants.JsonFormat:
		return json_validation.Validate
	}
	return nil
}
