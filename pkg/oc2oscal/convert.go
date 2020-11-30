package oc2oscal

import (
	"fmt"
	"os"
	"time"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/gocomply/fedramp/pkg/oc2oscal/masonry"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/oscal_source"
	"github.com/gocomply/oscalkit/pkg/uuid"
	"github.com/gocomply/oscalkit/types/oscal"
	ssp "github.com/gocomply/oscalkit/types/oscal/system_security_plan"
	"github.com/gocomply/oscalkit/types/oscal/validation_root"
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
	log "github.com/sirupsen/logrus"
)

func Convert(repoUri, outputDirectory string, format constants.DocumentFormat) error {
	workspace, err := masonry.Open(repoUri)
	if err != nil {
		return err
	}

	_, err = os.Stat(outputDirectory)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputDirectory, 0755)
	}
	if err != nil {
		return err
	}

	fedrampBaselines, err := fedramp.AvailableBaselines()
	if err != nil {
		return err
	}

	for _, component := range workspace.GetAllComponents() {
		controls, err := NewComponent(component)
		if err != nil {
			return err
		}
		for _, baseline := range fedrampBaselines {
			log.Debugf("Converting opencontrols for %s to FedRAMP %s", component.GetKey(), baseline.Level.Name())
			err = convertComponent(baseline, controls, outputDirectory, format)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func convertComponent(baseline fedramp.Baseline, component *Component, outputDirectory string, format constants.DocumentFormat) error {
	plan, err := fedramp.GSATemplate()
	if err != nil {
		return err
	}
	plan.Metadata.Title = &ssp.Title{PlainText: "FedRAMP System Security Plan (SSP)"}
	plan.Metadata.LastModified = validation_root.LastModified(time.Now().Format(constants.FormatDatetimeTz))
	plan.Metadata.Version = validation_root.Version("0.0.1")
	plan.Metadata.OscalVersion = validation_root.OscalVersion(constants.LatestOscalVersion)

	plan.ImportProfile = &ssp.ImportProfile{
		Href: baseline.ProfileURL(),
	}
	plan.SystemCharacteristics = convertSystemCharacteristics(component)
	sspComponent, err := buildSspComponent(component)
	if err != nil {
		return err
	}
	user := ssp.User{
		RoleIds: []ssp.RoleId{
			"generator",
		}}
	err = uuid.Refresh(&user)
	if err != nil {
		return err
	}

	plan.SystemImplementation = &ssp.SystemImplementation{
		Users:      []ssp.User{user},
		Components: []ssp.Component{sspComponent},
	}
	plan.ControlImplementation, err = convertControlImplementation(baseline, component, &sspComponent)
	if err != nil {
		return err
	}
	err = uuid.Refresh(plan)
	if err != nil {
		return err
	}
	filePath := outputDirectory + "/" + component.GetKey() + "-fedramp-" + baseline.Level.Name() + "." + format.String()
	err = writeSSP(plan, filePath, format)
	if err != nil {
		return err
	}
	return validate(filePath)
}

func buildSspComponent(oc *Component) (ssp.Component, error) {
	component := ssp.Component{
		ComponentType: "system",
		Title:         validation_root.ML("This system"),
		Description:   validation_root.MML("The entire system as depicted in the system authorization boundary"),
		Status:        &ssp.Status{State: "under-development"},
	}
	err := uuid.Refresh(&component)
	return component, err
}

func validate(filePath string) error {
	os, err := oscal_source.Open(filePath)
	if err != nil {
		return fmt.Errorf("Cannot read %s for validation: %s", filePath, err)
	}
	defer os.Close()
	err = os.Validate()
	if err != nil {
		return fmt.Errorf("Cannot validate %s: %s", filePath, err)
	}
	log.Debugf("Exported file has validated successfully: %s", filePath)
	return err
}

func convertControlImplementation(baseline fedramp.Baseline, component *Component, sspComponent *ssp.Component) (*ssp.ControlImplementation, error) {
	var ci ssp.ControlImplementation
	ci.Description = validation_root.MML("FedRAMP SSP Template Section 13")
	ci.ImplementedRequirements = make([]ssp.ImplementedRequirement, 0)

	if len(baseline.Controls()) != 0 {
		return nil, fmt.Errorf("Fedramp %s includes direct controls, those are not implemented yet", baseline.Level.Name())
	}

	for _, grp := range baseline.ControlGroups() {
		if len(grp.Groups) != 0 {
			return nil, fmt.Errorf("Fedramp %s includes nested control groups (inside group/@id=%s), those are not implemented yet", baseline.Level.Name(), grp.Id)
		}

		for _, ctrl := range grp.Controls {
			sat := component.GetSatisfy(ctrl.Id)
			if sat == nil {
				if baseline.Level.Name() == "High" {
					log.Warnf("Did not found control response for %s in %s\n", ctrl.Id, component.GetKey())
				}
				continue
			}
			stmts, err := convertStatements(ctrl.Id, sat.GetNarratives(), sspComponent)
			if err != nil {
				return nil, err
			}
			ir := ssp.ImplementedRequirement{
				ControlId: ctrl.Id,
				Annotations: []ssp.Annotation{
					fedrampImplementationStatus(sat.GetImplementationStatus()),
				},
				Statements: stmts,
			}
			err = uuid.Refresh(&ir)
			if err != nil {
				return nil, err
			}
			ci.ImplementedRequirements = append(ci.ImplementedRequirements, ir)

			for _, subctrl := range ctrl.Controls {
				if len(subctrl.Controls) != 0 {
					return nil, fmt.Errorf("3 layers of nested controls detected within %s", subctrl.Id)
				}
				sat = component.GetSatisfy(subctrl.Id)
				if sat == nil {
					if baseline.Level.Name() == "High" {
						log.Warnf("Did not found control response for %s in %s\n", subctrl.Id, component.GetKey())
					}
					continue
				}
				stmts, err := convertStatements(subctrl.Id, sat.GetNarratives(), sspComponent)
				if err != nil {
					return nil, err
				}
				ir := ssp.ImplementedRequirement{
					ControlId: subctrl.Id,
					Annotations: []ssp.Annotation{
						fedrampImplementationStatus(sat.GetImplementationStatus()),
					},
					Statements: stmts,
				}
				err = uuid.Refresh(&ir)
				if err != nil {
					return nil, err
				}
				ci.ImplementedRequirements = append(ci.ImplementedRequirements, ir)
			}
		}
	}
	return &ci, nil
}

func convertStatements(id string, narratives []common.Section, sspComponent *ssp.Component) ([]ssp.Statement, error) {
	var res []ssp.Statement
	if len(narratives) == 1 {
		stmt, err := newStatement(id, "", narratives[0].GetText(), sspComponent)
		return append(res, *stmt), err
	}

	for _, narrative := range narratives {
		stmt, err := newStatement(id, narrative.GetKey(), narrative.GetText(), sspComponent)
		if err != nil {
			return nil, err
		}
		res = append(res, *stmt)
	}
	return res, nil
}

func newStatement(controlId, narrativeId, narrative string, sspComponent *ssp.Component) (*ssp.Statement, error) {
	narrativeSuffix := ""
	if narrativeId != "" {
		narrativeSuffix = "." + narrativeId
	}
	byComponent := ssp.ByComponent{
		Description:   validation_root.MML("Describe how is the software component satisfying the control."),
		Remarks:       validation_root.MML(narrative),
		ComponentUuid: sspComponent.Uuid,
	}
	err := uuid.Refresh(&byComponent)
	if err != nil {
		return nil, err
	}
	statement := ssp.Statement{
		StatementId:  fmt.Sprintf("%s_stmt%s", controlId, narrativeSuffix),
		ByComponents: []ssp.ByComponent{byComponent},
	}
	err = uuid.Refresh(&statement)
	return &statement, err
}

func fedrampImplementationStatus(status string) ssp.Annotation {
	// Based on "Guide to OSCAL-based FedRAMP System Security Plans" (Version 1.0, November 27, 2019)
	// 5.3. Implementation Status (page 53)
	if status == "not applicable" {
		status = "not-applicable"
	}
	return ssp.Annotation{
		Name:  "implementation-status",
		Ns:    "https://fedramp.gov/ns/oscal",
		Value: status,
	}
}

func convertSystemCharacteristics(component *Component) *ssp.SystemCharacteristics {
	var syschar ssp.SystemCharacteristics
	syschar.SystemIds = []ssp.SystemId{
		ssp.SystemId{
			IdentifierType: "https://fedramp.gov",
			Id:             "F00000000",
		},
	}
	syschar.SystemName = ssp.SystemName(component.GetName())
	syschar.SystemNameShort = ssp.SystemNameShort(component.GetKey())
	syschar.Description = validation_root.MML("Automatically generated OSCAL SSP from OpenControl guidance for " + component.GetName())
	syschar.SecuritySensitivityLevel = ssp.SecuritySensitivityLevel("low")
	syschar.SystemInformation = staticSystemInformation()
	syschar.SecurityImpactLevel = &ssp.SecurityImpactLevel{
		SecurityObjectiveConfidentiality: ssp.SecurityObjectiveConfidentiality("fips-199-moderate"),
		SecurityObjectiveIntegrity:       ssp.SecurityObjectiveIntegrity("fips-199-moderate"),
		SecurityObjectiveAvailability:    ssp.SecurityObjectiveAvailability("fips-199-moderate"),
	}
	syschar.Status = &ssp.Status{
		State: "operational",
	}
	syschar.AuthorizationBoundary = &ssp.AuthorizationBoundary{
		Description: validation_root.MML("A holistic, top-level explanation of the FedRAMP authorization boundary."),
	}
	return &syschar
}

func staticSystemInformation() *ssp.SystemInformation {
	var sysinf ssp.SystemInformation
	sysinf.InformationTypes = []ssp.InformationType{
		ssp.InformationType{
			Title:       validation_root.ML("Information Type Name"),
			Description: validation_root.MML("This item is useless nevertheless required."),
			ConfidentialityImpact: &ssp.ConfidentialityImpact{
				Base: "fips-199-moderate",
			},
			IntegrityImpact: &ssp.IntegrityImpact{
				Base: "fips-199-moderate",
			},
			AvailabilityImpact: &ssp.AvailabilityImpact{
				Base: "fips-199-moderate",
			},
		},
	}
	return &sysinf
}

func writeSSP(plan *ssp.SystemSecurityPlan, outputFile string, format constants.DocumentFormat) error {
	destFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Error opening output file %s: %s", outputFile, err)
	}
	defer destFile.Close()

	output := oscal.OSCAL{SystemSecurityPlan: plan}
	err = output.Write(destFile, format, true)
	if err != nil {
		return fmt.Errorf("Cannot write %s: %s", outputFile, err)
	}
	return nil
}
