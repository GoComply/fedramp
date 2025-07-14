package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/urfave/cli"
)

var ssadCommand = cli.Command{
	Name:  "ssad",
	Usage: "Storing and Sharing Authorization Data operations for R5.SSAD",
	Subcommands: []cli.Command{
		ssadCreateCommand,
		ssadAddDocCommand,
		ssadFinalizeCommand,
		ssadShareCommand,
	},
}

var ssadCreateCommand = cli.Command{
	Name:      "create",
	Usage:     "Create a new SSAD package",
	ArgsUsage: "[service-id]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "csp",
			Usage: "CSP name",
			Value: "CloudNative Inc",
		},
		cli.StringFlag{
			Name:  "impact",
			Usage: "Impact level (Low, Moderate, High)",
			Value: "Moderate",
		},
		cli.StringFlag{
			Name:  "auth-type",
			Usage: "Authorization type (JAB, Agency, FedRAMP Tailored)",
			Value: "Agency",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file",
			Value: "ssad-package.json",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Service ID is required", 1)
		}
		
		serviceID := c.Args()[0]
		
		// Create metadata
		metadata := fedramp.SSADMetadata{
			Title:             fmt.Sprintf("FedRAMP Authorization Package - %s", serviceID),
			Description:       "Complete FedRAMP authorization documentation package",
			AuthorizationType: c.String("auth-type"),
			ImpactLevel:       c.String("impact"),
			AuthorizationDate: time.Now(),
			ExpirationDate:    time.Now().AddDate(3, 0, 0), // 3 years
			CSPName:           c.String("csp"),
			PackageFormat:     "JSON",
			Tags:              []string{"fedramp", "authorization", c.String("impact")},
			Keywords:          []string{serviceID, c.String("csp"), "cloud"},
		}
		
		// Create package
		pkg := fedramp.NewSSADPackage(serviceID, metadata)
		pkg.AccessControl.Owner = c.String("csp")
		pkg.AccessControl.DataClassification = "Controlled Unclassified Information (CUI)"
		
		// Export
		data, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating package: %v", err), 1)
		}
		
		outputFile := c.String("output")
		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("SSAD Package created:\n")
		fmt.Printf("  Package ID: %s\n", pkg.PackageID)
		fmt.Printf("  Service: %s\n", pkg.ServiceOfferingID)
		fmt.Printf("  CSP: %s\n", pkg.Metadata.CSPName)
		fmt.Printf("  Impact Level: %s\n", pkg.Metadata.ImpactLevel)
		fmt.Printf("  Saved to: %s\n", outputFile)
		
		return nil
	},
}

var ssadAddDocCommand = cli.Command{
	Name:      "add-doc",
	Usage:     "Add a document to SSAD package",
	ArgsUsage: "[package-file] [doc-type] [doc-path]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "title",
			Usage: "Document title",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "Document version",
			Value: "1.0",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 3 {
			return cli.NewExitError("Package file, document type, and document path are required", 1)
		}
		
		packageFile := c.Args()[0]
		docType := c.Args()[1]
		docPath := c.Args()[2]
		
		// Read package
		data, err := os.ReadFile(packageFile)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading package: %v", err), 1)
		}
		
		var pkg fedramp.SSADPackage
		if err := json.Unmarshal(data, &pkg); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing package: %v", err), 1)
		}
		
		// Get document info
		docInfo, err := os.Stat(docPath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading document: %v", err), 1)
		}
		
		// Create document entry
		doc := fedramp.SSADDocument{
			DocumentID:   fmt.Sprintf("DOC-%s-%s", docType, time.Now().Format("20060102")),
			Title:        c.String("title"),
			Type:         docType,
			Format:       "JSON", // In practice, detect from file
			Version:      c.String("version"),
			CreatedDate:  docInfo.ModTime(),
			LastModified: docInfo.ModTime(),
			Author:       pkg.Metadata.CSPName,
			Size:         docInfo.Size(),
			Location:     docPath,
			AccessLevel:  "restricted",
		}
		
		if doc.Title == "" {
			doc.Title = fmt.Sprintf("%s Document", docType)
		}
		
		// Add to package
		if err := pkg.AddDocument(docType, doc); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error adding document: %v", err), 1)
		}
		
		// Save updated package
		updatedData, err := json.MarshalIndent(&pkg, "", "  ")
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating JSON: %v", err), 1)
		}
		
		if err := os.WriteFile(packageFile, updatedData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("Added document to package:\n")
		fmt.Printf("  Type: %s\n", docType)
		fmt.Printf("  Title: %s\n", doc.Title)
		fmt.Printf("  Size: %d bytes\n", doc.Size)
		fmt.Printf("  Hash: %s\n", doc.Hash)
		
		return nil
	},
}

var ssadFinalizeCommand = cli.Command{
	Name:      "finalize",
	Usage:     "Finalize SSAD package for distribution",
	ArgsUsage: "[package-file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "signed-by",
			Usage: "Person signing the package",
			Value: "Authorizing Official",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Package file is required", 1)
		}
		
		// Read package
		data, err := os.ReadFile(c.Args()[0])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading package: %v", err), 1)
		}
		
		var pkg fedramp.SSADPackage
		if err := json.Unmarshal(data, &pkg); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing package: %v", err), 1)
		}
		
		// Finalize
		if err := pkg.Finalize(c.String("signed-by")); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error finalizing package: %v", err), 1)
		}
		
		// Save
		finalData, err := json.MarshalIndent(&pkg, "", "  ")
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating JSON: %v", err), 1)
		}
		
		if err := os.WriteFile(c.Args()[0], finalData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("Package finalized:\n")
		fmt.Printf("  Status: %s\n", pkg.Status)
		fmt.Printf("  Package Hash: %s\n", pkg.IntegrityCheck.PackageHash)
		fmt.Printf("  Signed By: %s\n", pkg.IntegrityCheck.SignedBy)
		fmt.Printf("  Signature Date: %s\n", pkg.IntegrityCheck.SignatureDate.Format("2006-01-02"))
		
		return nil
	},
}

var ssadShareCommand = cli.Command{
	Name:      "share",
	Usage:     "Share SSAD package with an entity",
	ArgsUsage: "[package-file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "with",
			Usage: "Entity to share with",
			Value: "FedRAMP PMO",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Entity type (agency, 3pao, public)",
			Value: "agency",
		},
		cli.StringFlag{
			Name:  "access",
			Usage: "Access level (read, write, admin)",
			Value: "read",
		},
		cli.IntFlag{
			Name:  "days",
			Usage: "Number of days access is valid",
			Value: 90,
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Package file is required", 1)
		}
		
		// Read package
		data, err := os.ReadFile(c.Args()[0])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error reading package: %v", err), 1)
		}
		
		var pkg fedramp.SSADPackage
		if err := json.Unmarshal(data, &pkg); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error parsing package: %v", err), 1)
		}
		
		// Create sharee
		expiration := time.Now().AddDate(0, 0, c.Int("days"))
		sharee := fedramp.SSADSharee{
			EntityID:       fmt.Sprintf("ENT-%s-%s", c.String("type"), time.Now().Format("20060102")),
			EntityType:     c.String("type"),
			Name:           c.String("with"),
			Email:          fmt.Sprintf("contact@%s.gov", c.String("with")),
			AccessLevel:    c.String("access"),
			ExpirationDate: &expiration,
		}
		
		// Share
		pkg.ShareWith(sharee)
		
		// Log distribution
		pkg.LogDistribution(sharee.Name, "Authorization review", "portal")
		
		// Save
		sharedData, err := json.MarshalIndent(&pkg, "", "  ")
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error generating JSON: %v", err), 1)
		}
		
		if err := os.WriteFile(c.Args()[0], sharedData, 0644); err != nil {
			return cli.NewExitError(fmt.Sprintf("Error writing file: %v", err), 1)
		}
		
		fmt.Printf("Package shared:\n")
		fmt.Printf("  Shared with: %s\n", sharee.Name)
		fmt.Printf("  Entity type: %s\n", sharee.EntityType)
		fmt.Printf("  Access level: %s\n", sharee.AccessLevel)
		fmt.Printf("  Expires: %s\n", sharee.ExpirationDate.Format("2006-01-02"))
		
		return nil
	},
} 