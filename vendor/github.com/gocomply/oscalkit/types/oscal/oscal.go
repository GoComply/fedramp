package oscal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"

	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/types/oscal/catalog"
	"github.com/gocomply/oscalkit/types/oscal/component_definition"
	"github.com/gocomply/oscalkit/types/oscal/profile"
	ssp "github.com/gocomply/oscalkit/types/oscal/system_security_plan"
	yaml "gopkg.in/yaml.v2"
)

const (
	catalogRootElement = "catalog"
	profileRootElement = "profile"
	sspRootElement     = "system-security-plan"
	componentElement   = "component-definition"
)

// OSCAL contains specific OSCAL components
type OSCAL struct {
	XMLName xml.Name         `json:"-" yaml:"-"`
	Catalog *catalog.Catalog `json:"catalog,omitempty" yaml:"catalog,omitempty"`
	// Declarations *Declarations `json:"declarations,omitempty" yaml:"declarations,omitempty"`
	Profile                 *profile.Profile `json:"profile,omitempty" yaml:"profile,omitempty"`
	*ssp.SystemSecurityPlan `xml:"system-security-plan"`
	Component               *component_definition.ComponentDefinition
	documentType            constants.DocumentType
}

func (o *OSCAL) DocumentType() constants.DocumentType {
	if o.Catalog != nil {
		return constants.CatalogDocument
	} else if o.Profile != nil {
		return constants.ProfileDocument
	} else if o.SystemSecurityPlan != nil {
		return constants.SSPDocument
	} else if o.Component != nil {
		return constants.ComponentDocument
	} else {
		return constants.UnknownDocument
	}
}

// MarshalXML marshals either a catalog or a profile
func (o *OSCAL) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if o.Catalog != nil {
		o.XMLName = o.Catalog.XMLName
		if err := e.Encode(o.Catalog); err != nil {
			return err
		}
	} else if o.Profile != nil {
		o.XMLName = o.Profile.XMLName
		if err := e.Encode(o.Profile); err != nil {
			return err
		}
	} else if o.SystemSecurityPlan != nil {
		o.XMLName = o.SystemSecurityPlan.XMLName
		if err := e.Encode(o.SystemSecurityPlan); err != nil {
			return err
		}
	}

	return nil
}

// New returns a concrete OSCAL type from a reader
func New(r io.Reader) (*OSCAL, error) {
	oscalBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	d := xml.NewDecoder(bytes.NewReader(oscalBytes))
	for {
		token, err := d.Token()
		if err != nil || token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			switch startElement.Name.Local {
			case catalogRootElement:
				var catalog catalog.Catalog
				if err := d.DecodeElement(&catalog, &startElement); err != nil {
					return nil, err
				}
				return &OSCAL{Catalog: &catalog}, nil

			case profileRootElement:
				var profile profile.Profile
				if err := d.DecodeElement(&profile, &startElement); err != nil {
					return nil, err
				}
				return &OSCAL{Profile: &profile}, nil
			case sspRootElement:
				var ssp ssp.SystemSecurityPlan
				if err := d.DecodeElement(&ssp, &startElement); err != nil {
					return nil, err
				}
				return &OSCAL{SystemSecurityPlan: &ssp}, nil
			case componentElement:
				var component component_definition.ComponentDefinition
				if err := d.DecodeElement(&component, &startElement); err != nil {
					return nil, err
				}
				return &OSCAL{Component: &component}, nil
			}
		}
	}

	var oscalT map[string]json.RawMessage
	if err := json.Unmarshal(oscalBytes, &oscalT); err == nil {
		for k, v := range oscalT {
			switch k {
			case catalogRootElement:
				var catalog catalog.Catalog
				if err := json.Unmarshal(v, &catalog); err != nil {
					return nil, err
				}
				return &OSCAL{Catalog: &catalog}, nil

			case profileRootElement:
				var profile profile.Profile
				if err := json.Unmarshal(v, &profile); err != nil {
					return nil, err
				}
				return &OSCAL{Profile: &profile}, nil
			}
		}
	}

	return nil, errors.New("Malformed OSCAL. Must be XML or JSON")
}

// XML writes the OSCAL object as XML to the given writer
func (o *OSCAL) XML(w io.Writer, prettify bool) error {
	return o.encode(encodeOptions{"xml", prettify, w})
}

// JSON writes the OSCAL object as JSON to the given writer
func (o *OSCAL) JSON(w io.Writer, prettify bool) error {
	return o.encode(encodeOptions{"json", prettify, w})
}

// YAML writes the OSCAL object as YAML to the given writer
func (o *OSCAL) YAML(w io.Writer) error {
	return o.encode(encodeOptions{format: "yaml", writer: w})
}

type encodeOptions struct {
	format   string
	prettify bool
	writer   io.Writer
}

func (o *OSCAL) encode(options encodeOptions) error {
	switch options.format {
	case "xml":
		e := xml.NewEncoder(options.writer)
		if options.prettify {
			e.Indent("", "  ")
		}

		return e.Encode(o)

	case "json":
		e := json.NewEncoder(options.writer)
		if options.prettify {
			e.SetIndent("", "  ")
		}

		return e.Encode(o)

	case "yaml":
		return yaml.NewEncoder(options.writer).Encode(o)
	}

	return errors.New("Incorrect format specified")
}
