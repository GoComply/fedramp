package bundled

import (
	"fmt"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"io"
	"io/ioutil"
	"os"
)

var schemaPaths = map[constants.DocumentFormat]map[constants.DocumentType]string{
	constants.XmlFormat: {
		constants.CatalogDocument:           "/OSCAL/xml/schema/oscal_catalog_schema.xsd",
		constants.ProfileDocument:           "/OSCAL/xml/schema/oscal_profile_schema.xsd",
		constants.SSPDocument:               "/OSCAL/xml/schema/oscal_ssp_schema.xsd",
		constants.ComponentDocument:         "/OSCAL/xml/schema/oscal_component_schema.xsd",
		constants.POAMDocument:              "/OSCAL/xml/schema/oscal_poam_schema.xsd",
		constants.AssessmentPlanDocument:    "/OSCAL/xml/schema/oscal_assessment-plan_schema.xsd",
		constants.AssessmentResultsDocument: "/OSCAL/xml/schema/oscal_assessment-results_schema.xsd",
	},
	constants.JsonFormat: {
		constants.CatalogDocument:           "/OSCAL/json/schema/oscal_catalog_schema.json",
		constants.ProfileDocument:           "/OSCAL/json/schema/oscal_profile_schema.json",
		constants.SSPDocument:               "/OSCAL/json/schema/oscal_ssp_schema.json",
		constants.ComponentDocument:         "/OSCAL/json/schema/oscal_component_schema.json",
		constants.POAMDocument:              "/OSCAL/json/schema/oscal_poam_schema.json",
		constants.AssessmentPlanDocument:    "/OSCAL/json/schema/oscal_assessment-plan_schema.json",
		constants.AssessmentResultsDocument: "/OSCAL/json/schema/oscal_assessment-results_schema.json",
	},
}

func noop() {
	// Hint pkger tool to bundle these files
	pkger.Include("/OSCAL/xml/schema/")
	pkger.Include("/OSCAL/json/schema/")
}

type BundledFile struct {
	Path string
}

func Schema(fileFormat constants.DocumentFormat, oscalComponent constants.DocumentType) (*BundledFile, error) {
	if oscalComponent == constants.UnknownDocument {
		return nil, fmt.Errorf("Unknown document format: %s", oscalComponent.String())
	}

	schemas, ok := schemaPaths[fileFormat]
	if !ok {
		return nil, fmt.Errorf("Cannot find schema for FileFormat %d", fileFormat)
	}
	schemaPath, ok := schemas[oscalComponent]
	if !ok {
		return nil, fmt.Errorf("Cannot find schema for document type '%s'", oscalComponent.String())
	}

	return localBundledFile(pkger.Open(schemaPath))
}

func HtmlXslt() (*BundledFile, error) {
	return localBundledFile(pkger.Open("/OSCAL/src/utils/util/publish/XSLT/oscal-browser-display.xsl"))
}

func localBundledFile(in pkging.File, err error) (*BundledFile, error) {
	if err != nil {
		return nil, err
	}
	defer in.Close()

	out, err := ioutil.TempFile("/tmp", "oscal")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return &BundledFile{Path: out.Name()}, nil
}

func (f *BundledFile) Cleanup() {
	os.Remove(f.Path)
}
