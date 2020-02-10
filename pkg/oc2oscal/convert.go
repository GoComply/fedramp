package oc2oscal

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GoComply/fedramp/pkg/fedramp"
	"github.com/GoComply/fedramp/pkg/oc2oscal/masonry"
	"github.com/docker/oscalkit/pkg/oscal/constants"
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
		for _, baseline := range fedrampBaselines {
			err = convertComponent(baseline, component, metadata, outputDirectory)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func convertComponent(baseline fedramp.Baseline, component common.Component, metadata ssp.Metadata, outputDirectory string) error {
	var plan ssp.SystemSecurityPlan
	plan.Id = "TODO"
	plan.Metadata = &metadata
	plan.ImportProfile = &ssp.ImportProfile{
		Href: baseline.ProfileURL(),
	}
	plan.SystemCharacteristics = convertSystemCharacteristics(component)
	plan.ControlImplementation = convertControlImplementation(component)
	return writeSSP(plan, outputDirectory+"/"+component.GetKey()+"-fedramp-"+baseline.Level.Name()+".xml")
}

func convertControlImplementation(component common.Component) *ssp.ControlImplementation {
	var ci ssp.ControlImplementation
	ci.Description = validation_root.MarkupFromPlain("FedRAMP SSP Template Section 13")
	ci.ImplementedRequirements = make([]ssp.ImplementedRequirement, 0)
	for _, sat := range component.GetAllSatisfies() {
		id := convertControlId(sat.GetControlKey())

		ci.ImplementedRequirements = append(ci.ImplementedRequirements, ssp.ImplementedRequirement{
			ControlId: id,
			Annotations: []ssp.Annotation{
				fedrampImplementationStatus(sat.GetImplementationStatus()),
			},
			Statements: convertStatements(id, sat.GetNarratives()),
		})
	}
	return &ci
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

func convertControlId(controlKey string) string {
	lower := strings.ToLower(controlKey)
	re := regexp.MustCompile("^([a-z][a-z])-([0-9]+)(\\s+\\(([0-9]+)\\))?$")
	match := re.FindStringSubmatch(lower)
	result := fmt.Sprintf("%s-%s", match[1], match[2])
	if match[4] != "" {
		result = fmt.Sprintf("%s.%s", result, match[4])
	}
	return result

}

func convertSystemCharacteristics(component common.Component) *ssp.SystemCharacteristics {
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
