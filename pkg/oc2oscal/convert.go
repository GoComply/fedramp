package oc2oscal

import (
	"fmt"
	"os"
	"time"

	"github.com/GoComply/fedramp/pkg/fedramp"
	"github.com/GoComply/fedramp/pkg/oc2oscal/masonry"
	"github.com/docker/oscalkit/pkg/oscal/constants"
	"github.com/docker/oscalkit/pkg/oscal_source"
	"github.com/docker/oscalkit/types/oscal"
	ssp "github.com/docker/oscalkit/types/oscal/system_security_plan"
	"github.com/docker/oscalkit/types/oscal/validation_root"
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
)

func Convert(repoUri, outputDirectory string) error {
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

	var metadata ssp.Metadata
	metadata.Title = ssp.Title("FedRAMP System Security Plan (SSP)")
	metadata.LastModified = validation_root.LastModified(time.Now().Format(constants.FormatDatetimeTz))
	metadata.Version = validation_root.Version("0.0.1")
	metadata.OscalVersion = validation_root.OscalVersion(constants.LatestOscalVersion)

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
			err = convertComponent(baseline, controls, metadata, outputDirectory)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func convertComponent(baseline fedramp.Baseline, component *Component, metadata ssp.Metadata, outputDirectory string) error {
	var plan ssp.SystemSecurityPlan
	var err error
	plan.Id = "TODO"
	plan.Metadata = &metadata
	plan.ImportProfile = &ssp.ImportProfile{
		Href: baseline.ProfileURL(),
	}
	plan.SystemCharacteristics = convertSystemCharacteristics(component)
	plan.ControlImplementation, err = convertControlImplementation(baseline, component)
	if err != nil {
		return err
	}
	filePath := outputDirectory + "/" + component.GetKey() + "-fedramp-" + baseline.Level.Name() + ".xml"
	err = writeSSP(plan, filePath)
	if err != nil {
		return err
	}
	return validate(filePath)
}

func validate(filePath string) error {
	os, err := oscal_source.Open(filePath)
	if err != nil {
		return err
	}
	defer os.Close()
	return os.Validate()
}

func convertControlImplementation(baseline fedramp.Baseline, component *Component) (*ssp.ControlImplementation, error) {
	var ci ssp.ControlImplementation
	ci.Description = validation_root.MarkupFromPlain("FedRAMP SSP Template Section 13")
	ci.ImplementedRequirements = make([]ssp.ImplementedRequirement, 0)

	if len(baseline.Controls()) != 0 {
		return nil, fmt.Errorf("Fedramp %s includes direct controls, those are not implemented yet", baseline.Level.Name())
	}

	for _, grp := range baseline.ControlGroups() {
		if len(grp.Groups) != 0 {
			return nil, fmt.Errorf("Fedramp %s includes nested control groups (inside group/@id=), those are not implemented yet", baseline.Level.Name(), grp.Id)
		}

		for _, ctrl := range grp.Controls {
			sat := component.GetSatisfy(ctrl.Id)
			if sat == nil {
				if baseline.Level.Name() == "High" {
					fmt.Printf("Did not found control response for %s in %s\n", ctrl.Id, component.GetKey())
				}
				continue
			}
			ci.ImplementedRequirements = append(ci.ImplementedRequirements, ssp.ImplementedRequirement{
				ControlId: ctrl.Id,
				Annotations: []ssp.Annotation{
					fedrampImplementationStatus(sat.GetImplementationStatus()),
				},
				Statements: convertStatements(ctrl.Id, sat.GetNarratives()),
			})

			for _, subctrl := range ctrl.Controls {
				if len(subctrl.Controls) != 0 {
					return nil, fmt.Errorf("3 layers of nested controls detected within %s", subctrl.Id)
				}
				sat = component.GetSatisfy(subctrl.Id)
				if sat == nil {
					if baseline.Level.Name() == "High" {
						fmt.Printf("Did not found control response for %s in %s\n", subctrl.Id, component.GetKey())
					}
					continue
				}
				ci.ImplementedRequirements = append(ci.ImplementedRequirements, ssp.ImplementedRequirement{
					ControlId: subctrl.Id,
					Annotations: []ssp.Annotation{
						fedrampImplementationStatus(sat.GetImplementationStatus()),
					},
					Statements: convertStatements(subctrl.Id, sat.GetNarratives()),
				})
			}
		}
	}
	return &ci, nil
}

func convertStatements(id string, narratives []common.Section) []ssp.Statement {
	var res []ssp.Statement
	if len(narratives) == 1 {
		return append(res, ssp.Statement{
			StatementId: fmt.Sprintf("%s_stmt", id),
			Description: validation_root.MarkupFromPlain(narratives[0].GetText()),
		})

	}

	for _, narrative := range narratives {
		res = append(res, ssp.Statement{
			StatementId: fmt.Sprintf("%s_stmt.%s", id, narrative.GetKey()),
			Description: validation_root.MarkupFromPlain(narrative.GetText()),
		})

	}
	return res
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
			Value:          "F00000000",
		},
	}
	syschar.SystemName = ssp.SystemName(component.GetName())
	syschar.SystemNameShort = ssp.SystemNameShort(component.GetKey())
	syschar.Description = validation_root.MarkupFromPlain("Automatically generated OSCAL SSP from OpenControl guidance for " + component.GetName())
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
		Description: validation_root.MarkupFromPlain("A holistic, top-level explanation of the FedRAMP authorization boundary."),
	}
	return &syschar
}

func staticSystemInformation() *ssp.SystemInformation {
	var sysinf ssp.SystemInformation
	sysinf.InformationTypes = []ssp.InformationType{
		ssp.InformationType{
			Title:       "Information Type Name",
			Description: validation_root.MarkupFromPlain("This item is useless nevertheless required."),
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

func writeSSP(plan ssp.SystemSecurityPlan, outputFile string) error {
	destFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Error opening output file %s: %s", outputFile, err)
	}
	defer destFile.Close()

	output := oscal.OSCAL{SystemSecurityPlan: &plan}
	return output.XML(destFile, true)
}
