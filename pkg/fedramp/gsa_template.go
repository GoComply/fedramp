package fedramp

import (
	"fmt"
	"github.com/gocomply/fedramp/bundled"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/oscal_source"
	ssp "github.com/gocomply/oscalkit/types/oscal/system_security_plan"
)

func GSATemplate() (*ssp.SystemSecurityPlan, error) {
	file, err := bundled.TemplateOSCAL()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	source, err := oscal_source.OpenFromReader(file.Name(), file)
	if err != nil {
		return nil, err
	}
	defer source.Close()
	oscal := source.OSCAL()
	if oscal.DocumentType() != constants.SSPDocument {
		return nil, fmt.Errorf("Could not initiate FedRAMP. Expected catalog element in %s", file.Name())
	}
	return oscal.SystemSecurityPlan, nil

}
