package cmd

import (
	"fmt"
	"os"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/urfave/cli"
)

var scnCommand = cli.Command{
	Name:  "scn",
	Usage: "Significant Change Notification operations for R5.SCN BIR",
	Subcommands: []cli.Command{
		scnCreateCommand,
		scnValidateCommand,
		scnExportCommand,
		scnListCommand,
	},
}

var scnCreateCommand = cli.Command{
	Name:      "create",
	Usage:     "Create a new Significant Change Notification",
	ArgsUsage: "[service-id] [change-type] [description] [reason]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file for the SCN JSON",
			Value: "scn.json",
		},
		cli.StringFlag{
			Name:  "3pao",
			Usage: "Third Party Assessment Organization name",
		},
		cli.StringFlag{
			Name:  "poam",
			Usage: "Related POA&M identifier",
		},
		cli.StringFlag{
			Name:  "approver-name",
			Usage: "Name of the approver",
		},
		cli.StringFlag{
			Name:  "approver-title",
			Usage: "Title of the approver",
		},
		cli.StringSliceFlag{
			Name:  "affected-controls",
			Usage: "List of affected control IDs (can be specified multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "affected-components",
			Usage: "List of affected components (can be specified multiple times)",
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() != 4 {
			return cli.NewExitError("Exactly 4 arguments are required: service-id, change-type, description, reason", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		serviceID := c.Args()[0]
		changeType := c.Args()[1]
		description := c.Args()[2]
		reason := c.Args()[3]

		// Create new SCN
		scn := fedramp.NewSCN(serviceID, changeType, description, reason)

		// Set optional fields
		if threePAO := c.String("3pao"); threePAO != "" {
			scn.ThreePAOName = threePAO
		}
		if poam := c.String("poam"); poam != "" {
			scn.RelatedPOAM = poam
		}
		if approverName := c.String("approver-name"); approverName != "" {
			scn.ApproverName = approverName
		}
		if approverTitle := c.String("approver-title"); approverTitle != "" {
			scn.ApproverTitle = approverTitle
		}

		// Add affected controls
		for _, control := range c.StringSlice("affected-controls") {
			scn.AddAffectedControl(control)
		}

		// Add affected components
		for _, component := range c.StringSlice("affected-components") {
			scn.AddAffectedComponent(component)
		}

		// Classify SCN type
		if err := scn.ClassifySCNType(); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error classifying SCN type: %v", err), 1)
		}

		// Export to JSON
		jsonData, err := scn.ToJSON()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error converting SCN to JSON: %v", err), 1)
		}

		// Write to file
		outputFile := c.String("output")
		if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing SCN to file: %v", err), 1)
		}

		fmt.Printf("SCN created successfully and saved to %s\n", outputFile)
		fmt.Printf("SCN Type: %s\n", scn.SCNType)
		fmt.Printf("Service ID: %s\n", scn.ServiceOfferingID)
		fmt.Printf("Change Type: %s\n", scn.ChangeType)

		return nil
	},
}

var scnValidateCommand = cli.Command{
	Name:      "validate",
	Usage:     "Validate a Significant Change Notification JSON file",
	ArgsUsage: "[scn-file.json]",
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Exactly 1 argument is required: path to SCN JSON file", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		scnFile := c.Args()[0]

		// Read SCN file
		data, err := os.ReadFile(scnFile)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading SCN file: %v", err), 1)
		}

		// Parse SCN
		var scn fedramp.SignificantChangeNotification
		if err := scn.FromJSON(data); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing SCN JSON: %v", err), 1)
		}

		// Validate SCN
		if err := scn.ValidateForSubmission(); err != nil {
			fmt.Printf("❌ SCN validation failed: %v\n", err)
			return cli.NewExitError("SCN validation failed", 1)
		}

		fmt.Printf("✅ SCN validation successful\n")
		fmt.Printf("Service ID: %s\n", scn.ServiceOfferingID)
		fmt.Printf("SCN Type: %s\n", scn.SCNType)
		fmt.Printf("Status: %s\n", scn.Status)
		fmt.Printf("Affected Controls: %v\n", scn.ControlsAffected)
		fmt.Printf("Affected Components: %v\n", scn.ComponentsAffected)

		return nil
	},
}

var scnExportCommand = cli.Command{
	Name:  "export",
	Usage: "Export SCN in various formats for agency notification",
	Subcommands: []cli.Command{
		{
			Name:      "summary",
			Usage:     "Export SCN as human-readable summary",
			ArgsUsage: "[scn-file.json] [output.txt]",
			Before: func(c *cli.Context) error {
				if c.NArg() != 2 {
					return cli.NewExitError("Exactly 2 arguments are required: SCN file and output file", 1)
				}
				return nil
			},
			Action: func(c *cli.Context) error {
				scnFile := c.Args()[0]
				outputFile := c.Args()[1]

				// Read and parse SCN
				data, err := os.ReadFile(scnFile)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Error reading SCN file: %v", err), 1)
				}

				var scn fedramp.SignificantChangeNotification
				if err := scn.FromJSON(data); err != nil {
					return cli.NewExitError(fmt.Sprintf("Error parsing SCN JSON: %v", err), 1)
				}

				// Generate summary
				summary := generateSCNSummary(&scn)

				// Write summary
				if err := os.WriteFile(outputFile, []byte(summary), 0644); err != nil {
					return cli.NewExitError(fmt.Sprintf("Error writing summary: %v", err), 1)
				}

				fmt.Printf("SCN summary exported to %s\n", outputFile)
				return nil
			},
		},
	},
}

var scnListCommand = cli.Command{
	Name:  "list",
	Usage: "List and manage multiple SCNs",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service-id",
			Usage: "Filter by service ID",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Filter by SCN type (adaptive, transformative, impact-change)",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Filter by status",
		},
	},
	Action: func(c *cli.Context) error {
		// This is a placeholder for SCN management functionality
		// In a real implementation, this would read from a database or file system
		fmt.Println("SCN List functionality - placeholder for R5.SCN BIR implementation")
		fmt.Println("This would integrate with FedRAMP's authorization data storage")
		return nil
	},
}

func generateSCNSummary(scn *fedramp.SignificantChangeNotification) string {
	summary := fmt.Sprintf(`SIGNIFICANT CHANGE NOTIFICATION SUMMARY
===========================================

Service Offering ID: %s
Change Type: %s
SCN Classification: %s
Status: %s

Description:
%s

Reason for Change:
%s

Affected Security Controls: %v
Affected Components: %v

`, scn.ServiceOfferingID, scn.ChangeType, scn.SCNType, scn.Status,
		scn.ShortDescription, scn.ReasonForChange,
		scn.ControlsAffected, scn.ComponentsAffected)

	// Add type-specific information
	switch scn.SCNType {
	case fedramp.SCNAdaptive:
		if scn.DateOfChange != nil {
			summary += fmt.Sprintf("Date of Change: %s\n", scn.DateOfChange.Format("2006-01-02"))
		}
		if scn.VerificationSteps != "" {
			summary += fmt.Sprintf("Verification Steps: %s\n", scn.VerificationSteps)
		}
		if scn.NewRisks != "" {
			summary += fmt.Sprintf("New Risks Identified: %s\n", scn.NewRisks)
		}

	case fedramp.SCNTransformative:
		if scn.PlannedChangeDate != nil {
			summary += fmt.Sprintf("Planned Change Date: %s\n", scn.PlannedChangeDate.Format("2006-01-02"))
		}
		if scn.ThreePAOName != "" {
			summary += fmt.Sprintf("3PAO: %s\n", scn.ThreePAOName)
		}
		if scn.RollbackPlan != "" {
			summary += fmt.Sprintf("Rollback Plan: %s\n", scn.RollbackPlan)
		}
	}

	// Add approval information
	if scn.ApproverName != "" {
		summary += fmt.Sprintf("\nApproved by: %s (%s)\n", scn.ApproverName, scn.ApproverTitle)
	}

	summary += fmt.Sprintf("\nCreated: %s\nLast Updated: %s\n",
		scn.CreatedAt.Format("2006-01-02 15:04:05"),
		scn.UpdatedAt.Format("2006-01-02 15:04:05"))

	return summary
} 