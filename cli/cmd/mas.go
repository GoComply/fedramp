package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/urfave/cli"
)

var masCommand = cli.Command{
	Name:  "mas",
	Usage: "Minimum Assessment Standard operations for R5.MAS",
	Subcommands: []cli.Command{
		masCreateCommand,
		masAddMethodCommand,
		masCompleteCommand,
		masSummaryCommand,
	},
}

var masCreateCommand = cli.Command{
	Name:      "create",
	Usage:     "Create a new MAS assessment",
	ArgsUsage: "[service-id] [assessment-type]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "3pao",
			Usage: "3PAO organization name",
			Value: "SecureAssess Inc",
		},
		cli.StringFlag{
			Name:  "lead",
			Usage: "Lead assessor name",
			Value: "Jane Doe",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file",
			Value: "mas-assessment.json",
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() != 2 {
			return cli.NewExitError("Service ID and assessment type are required", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		serviceID := c.Args()[0]
		assessmentType := c.Args()[1]
		
		// Validate assessment type
		var masType fedramp.MASAssessmentType
		switch assessmentType {
		case "initial":
			masType = fedramp.MASInitial
		case "annual":
			masType = fedramp.MASAnnual
		case "significant_change":
			masType = fedramp.MASSignificant
		case "incident_response":
			masType = fedramp.MASIncident
		default:
			return cli.NewExitError(fmt.Sprintf("Invalid assessment type: %s", assessmentType), 1)
		}
		
		// Create assessment
		assessment := fedramp.NewMASAssessment(serviceID, masType)
		
		// Set 3PAO information
		assessment.ThreePAO = fedramp.AssessmentOrganization{
			Name:              c.String("3pao"),
			A2LAAccreditation: "R311-2025",
			LeadAssessor:      c.String("lead"),
			TeamMembers:       []string{"John Smith", "Alice Johnson"},
			ContactEmail:      "assessments@secureassess.example",
		}
		
		// Set scope
		assessment.Scope = fedramp.AssessmentScope{
			FullAssessment:     masType == fedramp.MASInitial,
			Locations:          []string{"AWS us-east-1", "AWS us-west-2"},
			DataTypes:          []string{"CUI", "PII"},
			UserPopulation:     5000,
		}
		
		// Export
		data, err := assessment.ToJSON()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating assessment: %v", err), 1)
		}
		
		outputFile := c.String("output")
		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("MAS Assessment created:\n")
		fmt.Printf("  Assessment ID: %s\n", assessment.AssessmentID)
		fmt.Printf("  Type: %s\n", assessment.AssessmentType)
		fmt.Printf("  3PAO: %s\n", assessment.ThreePAO.Name)
		fmt.Printf("  Saved to: %s\n", outputFile)
		
		// Show requirements
		fmt.Printf("\nRequired assessment activities:\n")
		for _, req := range assessment.GetRequirements() {
			fmt.Printf("  - %s: %s\n", req.Name, req.Description)
			fmt.Printf("    Methods: %v\n", req.RequiredMethods)
			fmt.Printf("    Evidence: %v\n", req.MinimumEvidence)
		}
		
		return nil
	},
}

var masAddMethodCommand = cli.Command{
	Name:      "add-method",
	Usage:     "Add an assessment method to a MAS assessment",
	ArgsUsage: "[assessment-file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "Method type (interview, examine, test)",
			Value: "test",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "Method description",
		},
		cli.BoolFlag{
			Name:  "automated",
			Usage: "Whether the method is automated",
		},
		cli.StringSliceFlag{
			Name:  "controls",
			Usage: "Controls covered by this method",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Assessment file path is required", 1)
		}
		
		// Read existing assessment
		data, err := os.ReadFile(c.Args()[0])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading assessment: %v", err), 1)
		}
		
		var assessment fedramp.MASAssessment
		if err := json.Unmarshal(data, &assessment); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing assessment: %v", err), 1)
		}
		
		// Add method
		method := fedramp.AssessmentMethod{
			MethodType:      c.String("type"),
			Description:     c.String("description"),
			Automated:       c.Bool("automated"),
			StartTime:       time.Now(),
			EndTime:         time.Now().Add(2 * time.Hour),
			ControlsCovered: c.StringSlice("controls"),
		}
		
		if method.Automated {
			method.ToolsUsed = []string{"Automated Scanner v2.0"}
		}
		
		assessment.AddMethod(method)
		
		// Save updated assessment
		updatedData, err := assessment.ToJSON()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating JSON: %v", err), 1)
		}
		
		if err := os.WriteFile(c.Args()[0], updatedData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("Added assessment method:\n")
		fmt.Printf("  Type: %s\n", method.MethodType)
		fmt.Printf("  Automated: %v\n", method.Automated)
		fmt.Printf("  Controls: %v\n", method.ControlsCovered)
		
		return nil
	},
}

var masCompleteCommand = cli.Command{
	Name:      "complete",
	Usage:     "Mark a MAS assessment as complete",
	ArgsUsage: "[assessment-file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "attestor",
			Usage: "Attestor name",
			Value: "Jane Doe",
		},
		cli.StringFlag{
			Name:  "title",
			Usage: "Attestor title",
			Value: "Principal Assessor",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Assessment file path is required", 1)
		}
		
		// Read assessment
		data, err := os.ReadFile(c.Args()[0])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading assessment: %v", err), 1)
		}
		
		var assessment fedramp.MASAssessment
		if err := json.Unmarshal(data, &assessment); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing assessment: %v", err), 1)
		}
		
		// Add sample findings
		assessment.AddFinding(fedramp.AssessmentFinding{
			FindingID:      "FIND-001",
			ControlID:      "AC-2",
			Severity:       "moderate",
			Status:         "open",
			Description:    "Service accounts lack regular review",
			Recommendation: "Implement quarterly service account reviews",
			DateIdentified: time.Now(),
		})
		
		// Add attestation
		assessment.Attestation = &fedramp.AssessmentAttestation{
			AttestorName:  c.String("attestor"),
			AttestorTitle: c.String("title"),
			Organization:  assessment.ThreePAO.Name,
			Date:          time.Now(),
			Statement:     "I attest that this assessment was conducted in accordance with FedRAMP MAS requirements",
			Signature:     fmt.Sprintf("Signed-%s-%s", c.String("attestor"), time.Now().Format("20060102")),
		}
		
		// Complete assessment
		assessment.Complete()
		
		// Validate completeness
		if err := assessment.ValidateCompleteness(); err != nil {
			fmt.Printf("Warning: Assessment may be incomplete: %v\n", err)
		}
		
		// Save
		updatedData, err := assessment.ToJSON()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating JSON: %v", err), 1)
		}
		
		if err := os.WriteFile(c.Args()[0], updatedData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("Assessment completed:\n")
		fmt.Printf("  Status: %s\n", assessment.Status)
		fmt.Printf("  Findings: %d\n", len(assessment.Findings))
		fmt.Printf("  Attestor: %s\n", assessment.Attestation.AttestorName)
		
		return nil
	},
}

var masSummaryCommand = cli.Command{
	Name:      "summary",
	Usage:     "Generate summary of a MAS assessment",
	ArgsUsage: "[assessment-file]",
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Assessment file path is required", 1)
		}
		
		// Read assessment
		data, err := os.ReadFile(c.Args()[0])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading assessment: %v", err), 1)
		}
		
		var assessment fedramp.MASAssessment
		if err := json.Unmarshal(data, &assessment); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing assessment: %v", err), 1)
		}
		
		// Generate summary
		summary := assessment.GenerateSummary()
		
		fmt.Printf("\nMAS Assessment Summary\n")
		fmt.Printf("======================\n")
		fmt.Printf("Assessment ID: %s\n", summary["assessment_id"])
		fmt.Printf("Type: %s\n", summary["type"])
		fmt.Printf("Status: %s\n", summary["status"])
		fmt.Printf("Duration: %s\n", summary["duration"])
		fmt.Printf("3PAO: %s\n", summary["3pao"])
		fmt.Printf("\nFindings:\n")
		fmt.Printf("  Total: %d\n", summary["total_findings"])
		if findings, ok := summary["findings_by_severity"].(map[string]int); ok {
			for severity, count := range findings {
				fmt.Printf("  %s: %d\n", severity, count)
			}
		}
		fmt.Printf("\nEvidence:\n")
		fmt.Printf("  Total collected: %d\n", summary["evidence_collected"])
		fmt.Printf("  Automated: %d\n", summary["automated_evidence"])
		fmt.Printf("  Manual: %d\n", summary["manual_evidence"])
		
		return nil
	},
} 