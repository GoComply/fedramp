package fedramp

import (
	"fmt"
	"github.com/GoComply/fedramp/pkg/fedramp/common"
	"github.com/docker/oscalkit/pkg/oscal/constants"
	"github.com/docker/oscalkit/pkg/oscal_source"
	ssp "github.com/docker/oscalkit/types/oscal/system_security_plan"
)

type SSP struct {
	plan     ssp.SystemSecurityPlan
	baseline *Baseline
}

func NewSSP(sspSource *oscal_source.OSCALSource) (*SSP, error) {
	var result SSP
	var err error
	o := sspSource.OSCAL()
	if o.DocumentType() != constants.SSPDocument {
		return nil, fmt.Errorf("Provided OSCAL file is not system-security-plan")
	}
	result.plan = *o.SystemSecurityPlan

	baseline := result.Level()
	if baseline == common.LevelUnknown {
		return nil, fmt.Errorf("Unrecognized FedRAMP profile URL: %s", result.plan.ImportProfile.Href)
	}

	result.baseline, err = NewBaseline(baseline)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *SSP) Level() common.BaselineLevel {
	if p.plan.ImportProfile == nil {
		return common.LevelUnknown
	}
	for baseline, href := range common.ProfileUrls {
		if href == p.plan.ImportProfile.Href {
			return baseline
		}
	}
	return common.LevelUnknown
}
