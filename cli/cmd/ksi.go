package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/urfave/cli"
)

var ksiCommand = cli.Command{
	Name:  "ksi",
	Usage: "FedRAMP 20x Key Security Indicators operations",
	Subcommands: []cli.Command{
		ksiValidateCommand,
		ksiReportCommand,
		ksiProposalCommand,
	},
}

var ksiValidateCommand = cli.Command{
	Name:      "validate",
	Usage:     "Validate Key Security Indicators",
	ArgsUsage: "[service-id]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file for the KSI validation report",
			Value: "ksi-report.json",
		},
		cli.BoolFlag{
			Name:  "automated",
			Usage: "Run automated validations only",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Service ID is required", 1)
		}
		
		serviceID := c.Args()[0]
		report := fedramp.NewKSIReport(serviceID)
		
		// Simulate KSI validations
		fmt.Printf("Validating KSIs for service: %s\n", serviceID)
		
		// Example validations
		validations := []struct {
			ksiID     string
			automated bool
			evidence  []fedramp.KSIEvidence
		}{
			{
				ksiID:     "KSI-CNA",
				automated: true,
				evidence: []fedramp.KSIEvidence{
					{
						Type:        "configuration",
						Description: "DoS protection enabled via CloudFlare",
						Reference:   "cloudflare-config-2025-01",
						Source:      "CloudFlare API",
					},
					{
						Type:        "scan_result",
						Description: "Container scanning shows immutable deployments",
						Reference:   "trivy-scan-2025-01",
						Source:      "CI/CD Pipeline",
					},
				},
			},
			{
				ksiID:     "KSI-SC",
				automated: true,
				evidence: []fedramp.KSIEvidence{
					{
						Type:        "certificate",
						Description: "TLS 1.3 enforced on all endpoints",
						Reference:   "ssl-labs-scan-2025-01",
						Source:      "SSL Labs API",
					},
				},
			},
			{
				ksiID:     "KSI-IAM",
				automated: false,
				evidence: []fedramp.KSIEvidence{
					{
						Type:        "policy",
						Description: "MFA policy enforced for all users",
						Reference:   "iam-policy-doc-v2",
						Source:      "Identity Provider",
					},
				},
			},
		}
		
		for _, v := range validations {
			if c.Bool("automated") && !v.automated {
				continue
			}
			
			validation := fedramp.ValidateKSI(v.ksiID, v.evidence, v.automated)
			if validation != nil {
				report.AddValidation(validation)
				fmt.Printf("✓ Validated %s: %s\n", v.ksiID, validation.Status)
			}
		}
		
		// Export report
		jsonData, err := report.ToJSON()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating report: %v", err), 1)
		}
		
		outputFile := c.String("output")
		if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing report: %v", err), 1)
		}
		
		fmt.Printf("\nKSI Validation Report Summary:\n")
		fmt.Printf("Total KSIs: %d\n", report.Summary.TotalKSIs)
		fmt.Printf("Automated: %d\n", report.Summary.AutomatedCount)
		fmt.Printf("Compliance Score: %.1f%%\n", report.Summary.ComplianceScore)
		fmt.Printf("Report saved to: %s\n", outputFile)
		
		return nil
	},
}

var ksiReportCommand = cli.Command{
	Name:  "report",
	Usage: "Generate continuous reporting data",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service-id",
			Usage: "Service offering ID",
			Value: "CSO-DEMO-001",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file",
			Value: "continuous-report.json",
		},
	},
	Action: func(c *cli.Context) error {
		serviceID := c.String("service-id")
		manager := fedramp.NewContinuousReportingManager(serviceID)
		
		// Generate continuous report
		reportData, err := manager.GenerateContinuousReport()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating report: %v", err), 1)
		}
		
		outputFile := c.String("output")
		if err := os.WriteFile(outputFile, reportData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing report: %v", err), 1)
		}
		
		fmt.Printf("Continuous monitoring report generated: %s\n", outputFile)
		return nil
	},
}

var ksiProposalCommand = cli.Command{
	Name:  "proposal",
	Usage: "Generate continuous reporting proposal",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service-id",
			Usage: "Service offering ID",
			Value: "CSO-DEMO-001",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file",
			Value: "continuous-reporting-proposal.json",
		},
	},
	Action: func(c *cli.Context) error {
		serviceID := c.String("service-id")
		manager := fedramp.NewContinuousReportingManager(serviceID)
		
		// Generate proposal
		proposal := manager.GenerateProposal()
		
		// Convert to JSON
		jsonData, err := json.MarshalIndent(proposal, "", "  ")
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating proposal: %v", err), 1)
		}
		
		outputFile := c.String("output")
		if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing proposal: %v", err), 1)
		}
		
		fmt.Printf("Continuous Reporting Proposal Generated:\n")
		fmt.Printf("Service ID: %s\n", serviceID)
		fmt.Printf("Automated Coverage: %.1f%%\n", proposal.CoveragePercentage)
		fmt.Printf("Automated KSIs: %d\n", len(proposal.AutomatedKSIs))
		fmt.Printf("Proposal saved to: %s\n", outputFile)
		
		return nil
	},
} 