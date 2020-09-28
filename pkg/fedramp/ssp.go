package fedramp

import (
	"fmt"
	"github.com/gocomply/fedramp/pkg/fedramp/common"
	"github.com/gocomply/fedramp/pkg/utils"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/oscal_source"
	ssp "github.com/gocomply/oscalkit/types/oscal/system_security_plan"
)

type SSP struct {
	plan                         ssp.SystemSecurityPlan
	baseline                     *Baseline
	implementedRequirementsCache map[string]ssp.ImplementedRequirement
}

func NewSSP(sspSource *oscal_source.OSCALSource) (*SSP, error) {
	var result SSP
	var err error
	o := sspSource.OSCAL()
	if o.DocumentType() != constants.SSPDocument {
		return nil, fmt.Errorf("Provided OSCAL file is not system-security-plan")
	}
	result.plan = *o.SystemSecurityPlan

	if result.plan.ControlImplementation == nil {
		return nil, fmt.Errorf("SSP is missing control implementation section")

	}
	result.implementedRequirementsCache = make(map[string]ssp.ImplementedRequirement)
	for _, ir := range result.plan.ControlImplementation.ImplementedRequirements {
		result.implementedRequirementsCache[ir.ControlId] = ir
	}

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

func (p *SSP) ResponsibleRoleForControl(controlId string) string {
	ir, found := p.implementedRequirementsCache[utils.ControlKeyToOSCAL(controlId)]
	if !found {
		return "No information available"
	}
	if len(ir.ResponsibleRoles) == 0 {
		return "No information available"
	}

	return ir.ResponsibleRoles[0].RoleId
}

func (p *SSP) ParamValue(controlId string, index int) (string, error) {
	paramId := fmt.Sprintf("%s_prm_%d", controlId, index)
	param, err := p.baseline.FindParam(controlId, paramId)
	if err != nil {
		return "", err
	}
	if param != nil {
		for _, constraint := range param.Constraints {
			return string(constraint.Detail), nil
		}
		if param.Label != "" {
			return "[Assignments: " + string(param.Label) + "]", nil
		}
	}
	return "", nil
}

func (p *SSP) ImplementationStatusForControl(controlId string) ImplementationStatus {
	ir, found := p.implementedRequirementsCache[utils.ControlKeyToOSCAL(controlId)]
	if !found {
		return StatusNoStatus
	}

	for _, annotation := range ir.Annotations {
		if annotation.Name == "implementation-status" && annotation.Ns == "https://fedramp.gov/ns/oscal" {
			return StatusFromOSCAL(annotation.Value)
		}
	}
	return StatusNoStatus
}

func (p *SSP) StatementTextForPart(controlId, partName string) (string, error) {
	oscalControlId := utils.ControlKeyToOSCAL(controlId)
	ir, found := p.implementedRequirementsCache[oscalControlId]
	if !found {
		return "No information available", nil
	}

	stmt, err := findStatement(&ir, oscalControlId, partName)
	if err != nil {
		return "", err
	}

	if stmt == nil && len(partName) > 1 {
		partName = partName[:1]
		stmt, err = findStatement(&ir, oscalControlId, partName)
		if err != nil {
			return "", err
		}
	}
	if stmt != nil {
		return stmt.Description.PlainString(), nil
	}

	return "", nil
}

func findStatement(implementedRequirement *ssp.ImplementedRequirement, oscalControlId, partName string) (*ssp.Statement, error) {
	id := fmt.Sprintf("%s_stmt.%s", oscalControlId, partName)
	for _, stmt := range implementedRequirement.Statements {
		if stmt.StatementId == id {
			return &stmt, nil
		}
	}
	return nil, nil
}
