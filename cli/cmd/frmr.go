package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gocomply/fedramp/pkg/fedramp/frmr"
	"github.com/urfave/cli"
)

// FRMR returns the CLI command for FRMR operations
func FRMR() cli.Command {
	return cli.Command{
		Name:  "frmr",
		Usage: "Work with FedRAMP Machine Readable (FRMR) documents",
		Subcommands: []cli.Command{
			{
				Name:      "fetch",
				Usage:     "Fetch FRMR documents from the official FedRAMP repository",
				ArgsUsage: "[ksi|mas|scn]",
				Action:    fetchFRMR,
			},
			{
				Name:      "validate",
				Usage:     "Validate KSI requirements against evidence",
				ArgsUsage: "<frmr-file> <evidence-file>",
				Action:    validateFRMR,
			},
			{
				Name:      "info",
				Usage:     "Display information about a FRMR document",
				ArgsUsage: "<frmr-file>",
				Action:    infoFRMR,
			},
			{
				Name:      "export",
				Usage:     "Export FRMR data in various formats",
				ArgsUsage: "<frmr-file>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "format",
						Value: "json",
						Usage: "Output format (json, yaml, markdown)",
					},
					cli.StringFlag{
						Name:  "output",
						Usage: "Output file (default: stdout)",
					},
				},
				Action: exportFRMR,
			},
			{
				Name:      "combine",
				Usage:     "Combine multiple FRMR documents into one",
				ArgsUsage: "<frmr-file1> <frmr-file2> [...]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "output",
						Usage: "Output file (default: combined-frmr.json)",
						Value: "combined-frmr.json",
					},
				},
				Action: combineFRMR,
			},
			{
				Name:      "filter",
				Usage:     "Filter FRMR document by criteria",
				ArgsUsage: "<frmr-file>",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "impact",
						Usage: "Filter by impact levels (Low, Moderate, High)",
					},
					cli.StringSliceFlag{
						Name:  "ksi",
						Usage: "Filter by KSI IDs",
					},
					cli.StringSliceFlag{
						Name:  "type",
						Usage: "Include only specific types (FRD, FRR, FRA, KSI)",
					},
					cli.StringFlag{
						Name:  "output",
						Usage: "Output file (default: stdout)",
					},
				},
				Action: filterFRMR,
			},
			{
				Name:      "evidence-template",
				Usage:     "Generate an evidence template for KSI validation",
				ArgsUsage: "<frmr-file>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "output",
						Usage: "Output file (default: evidence-template.json)",
						Value: "evidence-template.json",
					},
				},
				Action: evidenceTemplateFRMR,
			},
			{
				Name:      "schema-validate",
				Usage:     "Validate FRMR document against schema",
				ArgsUsage: "<frmr-file>",
				Action:    schemaValidateFRMR,
			},
		},
	}
}

const frmrBaseURL = "https://raw.githubusercontent.com/FedRAMP/docs/main/"

var frmrFiles = map[string]string{
	"ksi": "FRMR.KSI.key-security-indicators.json",
	"mas": "FRMR.MAS.minimum-assessment-standard.json",
	"scn": "FRMR.SCN.significant-change-notifications.json",
}

func fetchFRMR(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("Please specify document type: ksi, mas, or scn", 1)
	}

	docType := strings.ToLower(c.Args()[0])
	filename, ok := frmrFiles[docType]
	if !ok {
		return cli.NewExitError(fmt.Sprintf("Unknown document type: %s", docType), 1)
	}

	url := frmrBaseURL + filename
	fmt.Printf("Fetching %s from %s...\n", docType, url)

	resp, err := http.Get(url)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to fetch document: %v", err), 1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cli.NewExitError(fmt.Sprintf("Failed to fetch document: HTTP %d", resp.StatusCode), 1)
	}

	// Parse to validate
	doc, err := frmr.ParseFRMR(resp.Body)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse document: %v", err), 1)
	}

	// Save to file
	outputFile := filename
	file, err := os.Create(outputFile)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create file: %v", err), 1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(doc); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to write file: %v", err), 1)
	}

	fmt.Printf("Successfully saved %s to %s\n", doc.Info.Name, outputFile)
	fmt.Printf("Current Release: %s\n", doc.Info.CurrentRelease)

	return nil
}

func validateFRMR(c *cli.Context) error {
	if c.NArg() < 2 {
		return cli.NewExitError("Please specify FRMR file and evidence file", 1)
	}

	// Load FRMR document
	frmrFile, err := os.Open(c.Args()[0])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open FRMR file: %v", err), 1)
	}
	defer frmrFile.Close()

	doc, err := frmr.ParseFRMR(frmrFile)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse FRMR document: %v", err), 1)
	}

	// Load evidence
	evidenceFile, err := os.Open(c.Args()[1])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open evidence file: %v", err), 1)
	}
	defer evidenceFile.Close()

	var evidence map[string]bool
	decoder := json.NewDecoder(evidenceFile)
	if err := decoder.Decode(&evidence); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse evidence file: %v", err), 1)
	}

	// Validate KSIs
	if len(doc.KSI) == 0 {
		return cli.NewExitError("No KSIs found in document", 1)
	}

	fmt.Printf("Validating %s\n", doc.Info.Name)
	fmt.Printf("Release: %s\n\n", doc.Info.CurrentRelease)

	totalKSIs := 0
	metKSIs := 0

	for _, ksi := range doc.KSI {
		result := ksi.ValidateKSI(evidence)
		totalKSIs++
		if result.FullyMet {
			metKSIs++
		}

		status := "❌"
		if result.FullyMet {
			status = "✅"
		}

		fmt.Printf("%s %s: %s\n", status, result.KSIID, result.KSIName)
		fmt.Printf("   Requirements: %d/%d met\n", result.MetCount, result.TotalCount)

		if !result.FullyMet {
			fmt.Println("   Unmet requirements:")
			for _, req := range result.Requirements {
				if !req.Met {
					fmt.Printf("   - %s: %s\n", req.ID, req.Statement)
				}
			}
		}
		fmt.Println()
	}

	fmt.Printf("Overall: %d/%d KSIs fully met (%.1f%%)\n", 
		metKSIs, totalKSIs, float64(metKSIs)/float64(totalKSIs)*100)

	return nil
}

func infoFRMR(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("Please specify FRMR file", 1)
	}

	file, err := os.Open(c.Args()[0])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %v", err), 1)
	}
	defer file.Close()

	doc, err := frmr.ParseFRMR(file)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse document: %v", err), 1)
	}

	fmt.Printf("Document: %s (%s)\n", doc.Info.Name, doc.Info.ShortName)
	fmt.Printf("Current Release: %s\n", doc.Info.CurrentRelease)
	fmt.Printf("Types: %s\n\n", strings.Join(doc.Info.Types, ", "))

	// Show current release details
	release, err := doc.GetCurrentRelease()
	if err == nil {
		fmt.Printf("Release %s:\n", release.ID)
		fmt.Printf("  Published: %s\n", release.PublishedDate)
		fmt.Printf("  Description: %s\n", release.Description)
		
		if release.Effective.TwentyX != nil {
			fmt.Printf("  20x Effective: %s\n", release.Effective.TwentyX.Timeline.Pilot.StartDate)
		}
		if release.Effective.R5 != nil {
			fmt.Printf("  R5 Effective: %s\n", release.Effective.R5.Timeline.Pilot.StartDate)
		}
	}

	// Show statistics
	fmt.Printf("\nContent:\n")
	if len(doc.KSI) > 0 {
		fmt.Printf("  KSIs: %d\n", len(doc.KSI))
		totalReqs := 0
		for _, ksi := range doc.KSI {
			totalReqs += len(ksi.Requirements)
		}
		fmt.Printf("  Total KSI Requirements: %d\n", totalReqs)
	}
	if len(doc.FRD) > 0 {
		totalDefs := 0
		for _, defs := range doc.FRD {
			totalDefs += len(defs)
		}
		fmt.Printf("  Definitions (FRD): %d\n", totalDefs)
	}
	if len(doc.FRR) > 0 {
		totalRules := 0
		for _, rrrBase := range doc.FRR {
			totalRules += len(rrrBase.Base.Requirements)
		}
		fmt.Printf("  Rules (FRR): %d\n", totalRules)
	}
	if len(doc.FRA) > 0 {
		totalAssist := 0
		for _, assists := range doc.FRA {
			totalAssist += len(assists)
		}
		fmt.Printf("  Assistance (FRA): %d\n", totalAssist)
	}

	return nil
}

func exportFRMR(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("Please specify FRMR file", 1)
	}

	file, err := os.Open(c.Args()[0])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %v", err), 1)
	}
	defer file.Close()

	doc, err := frmr.ParseFRMR(file)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse document: %v", err), 1)
	}

	format := c.String("format")
	output := c.String("output")

	var writer *os.File
	if output == "" {
		writer = os.Stdout
	} else {
		writer, err = os.Create(output)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Failed to create output file: %v", err), 1)
		}
		defer writer.Close()
	}

	switch format {
	case "json":
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(doc)
	
	case "markdown":
		return exportMarkdown(writer, doc)
	
	default:
		return cli.NewExitError(fmt.Sprintf("Unsupported format: %s", format), 1)
	}
}

func exportMarkdown(w *os.File, doc *frmr.FRMRDocument) error {
	exporter := frmr.NewMarkdownExporter()
	return exporter.Export(doc, w)
}

func combineFRMR(c *cli.Context) error {
	if c.NArg() < 2 {
		return cli.NewExitError("Please specify at least two FRMR files to combine", 1)
	}

	// Load all documents
	docs := make([]*frmr.FRMRDocument, 0, c.NArg())
	for i := 0; i < c.NArg(); i++ {
		file, err := os.Open(c.Args()[i])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Failed to open file %s: %v", c.Args()[i], err), 1)
		}
		defer file.Close()

		doc, err := frmr.ParseFRMR(file)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Failed to parse %s: %v", c.Args()[i], err), 1)
		}
		docs = append(docs, doc)
	}

	// Combine documents
	combined, err := frmr.CombineFRMRDocuments(docs...)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to combine documents: %v", err), 1)
	}

	// Write output
	output := c.String("output")
	file, err := os.Create(output)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create output file: %v", err), 1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(combined); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to write combined document: %v", err), 1)
	}

	fmt.Printf("Combined %d documents into %s\n", len(docs), output)
	return nil
}

func filterFRMR(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("Please specify FRMR file", 1)
	}

	file, err := os.Open(c.Args()[0])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %v", err), 1)
	}
	defer file.Close()

	doc, err := frmr.ParseFRMR(file)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse document: %v", err), 1)
	}

	// Build filter options
	opts := frmr.FilterOptions{
		ImpactLevels: c.StringSlice("impact"),
		KSIIDs:       c.StringSlice("ksi"),
		Types:        c.StringSlice("type"),
	}

	// Apply filter
	filtered := frmr.FilterDocument(doc, opts)

	// Write output
	var writer *os.File
	output := c.String("output")
	if output == "" {
		writer = os.Stdout
	} else {
		writer, err = os.Create(output)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Failed to create output file: %v", err), 1)
		}
		defer writer.Close()
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(filtered)
}

func evidenceTemplateFRMR(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("Please specify FRMR file", 1)
	}

	file, err := os.Open(c.Args()[0])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %v", err), 1)
	}
	defer file.Close()

	doc, err := frmr.ParseFRMR(file)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse document: %v", err), 1)
	}

	// Create output file
	output := c.String("output")
	outFile, err := os.Create(output)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create output file: %v", err), 1)
	}
	defer outFile.Close()

	// Generate and write template
	if err := frmr.ExportEvidenceTemplate(doc, outFile); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to generate evidence template: %v", err), 1)
	}

	fmt.Printf("Evidence template saved to %s\n", output)
	fmt.Printf("Edit the file and set each requirement ID to true/false based on your implementation.\n")
	return nil
}

func schemaValidateFRMR(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("Please specify FRMR file", 1)
	}

	file, err := os.Open(c.Args()[0])
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %v", err), 1)
	}
	defer file.Close()

	doc, err := frmr.ParseFRMR(file)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to parse document: %v", err), 1)
	}

	// Validate schema
	errors := frmr.ValidateSchema(doc)
	
	if len(errors) == 0 {
		fmt.Println("✅ Document is valid!")
		return nil
	}

	fmt.Printf("❌ Document has %d validation errors:\n\n", len(errors))
	for i, err := range errors {
		fmt.Printf("%d. %s\n", i+1, err)
	}
	
	return cli.NewExitError("Document validation failed", 1)
} 