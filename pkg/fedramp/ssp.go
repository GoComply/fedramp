package fedramp

import (
	"fmt"
	"github.com/docker/oscalkit/pkg/oscal/constants"
	"github.com/docker/oscalkit/pkg/oscal_source"
	ssp "github.com/docker/oscalkit/types/oscal/system_security_plan"
)

type SSP struct {
	responses ssp.SystemSecurityPlan
	baseline  Baseline
}

func NewSSP(sspSource *oscal_source.OSCALSource) (*SSP, error) {
	var result SSP
	o := sspSource.OSCAL()
	if o.DocumentType() != constants.SSPDocument {
		return nil, fmt.Errorf("Provided OSCAL file is not system-security-plan")
	}
	result.responses = *o.SystemSecurityPlan
	return &result, nil
}
