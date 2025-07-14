package frmr

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// CombineFRMRDocuments combines multiple FRMR documents into a single document
// This is useful for creating combined requirement sets (e.g., all 20x Low requirements)
func CombineFRMRDocuments(docs ...*FRMRDocument) (*FRMRDocument, error) {
	if len(docs) == 0 {
		return nil, fmt.Errorf("no documents provided")
	}

	// Start with the first document as base
	combined := &FRMRDocument{
		Schema: docs[0].Schema,
		ID:     "combined-frmr",
		Info: Info{
			Name:           "Combined FRMR Requirements",
			ShortName:      "COMBINED",
			CurrentRelease: docs[0].Info.CurrentRelease,
			Types:          []string{},
		},
		FRD: make(map[string][]FRDItem),
		FRR: make(map[string]FRRBase),
		FRA: make(map[string][]FRAItem),
		KSI: make(map[string]KSIItem),
	}

	// Track unique types
	typeSet := make(map[string]bool)

	// Combine all documents
	for _, doc := range docs {
		// Combine types
		for _, t := range doc.Info.Types {
			typeSet[t] = true
		}

		// Combine FRD
		for key, items := range doc.FRD {
			combined.FRD[key] = append(combined.FRD[key], items...)
		}

		// Combine FRR
		for key, base := range doc.FRR {
			if existing, ok := combined.FRR[key]; ok {
				// Merge requirements
				existing.Base.Requirements = append(existing.Base.Requirements, base.Base.Requirements...)
				combined.FRR[key] = existing
			} else {
				combined.FRR[key] = base
			}
		}

		// Combine FRA
		for key, items := range doc.FRA {
			combined.FRA[key] = append(combined.FRA[key], items...)
		}

		// Combine KSI
		for key, ksi := range doc.KSI {
			combined.KSI[key] = ksi
		}
	}

	// Convert type set to slice
	for t := range typeSet {
		combined.Info.Types = append(combined.Info.Types, t)
	}

	return combined, nil
}

// MarkdownExporter handles conversion of FRMR documents to Markdown
type MarkdownExporter struct {
	Template *template.Template
}

// NewMarkdownExporter creates a new markdown exporter with the default template
func NewMarkdownExporter() *MarkdownExporter {
	tmpl := template.New("markdown").Funcs(template.FuncMap{
		"join": strings.Join,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"replace": strings.ReplaceAll,
	})

	// Default template (simplified version of FedRAMP's Handlebars template)
	tmplText := `# {{.Info.Name}}

**Current Release:** {{.Info.CurrentRelease}}

{{if .KSI}}
## Key Security Indicators

{{range $id, $ksi := .KSI}}
### {{$ksi.ID}}: {{$ksi.Name}}

{{$ksi.Indicator}}

**Requirements:**
{{range $ksi.Requirements}}
- **{{.ID}}**: {{.Statement}}{{if .Controls}}
  - Controls: {{range $i, $c := .Controls}}{{if $i}}, {{end}}{{$c.ControlID}}{{end}}{{end}}
{{end}}
{{end}}
{{end}}

{{if .FRD}}
## Definitions

{{range $category, $items := .FRD}}
### {{$category}}
{{range $items}}
- **{{.Term}}** ({{.ID}}): {{.Definition}}{{if .Note}}
  - Note: {{.Note}}{{end}}
{{end}}
{{end}}
{{end}}

{{if .FRR}}
## Requirements

{{range $category, $base := .FRR}}
### {{$category}}

{{$base.Base.Application}}

{{range $base.Base.Requirements}}
- **{{.ID}}**: {{.Statement}}
  - Affects: {{join .Affects ", "}}
  - Key Word: {{.PrimaryKeyWord}}
{{end}}
{{end}}
{{end}}

{{if .FRA}}
## Assistance

{{range $category, $items := .FRA}}
### {{$category}}
{{range $items}}
#### {{.ID}}: {{.Title}}

{{.Description}}

{{if .Examples}}
**Examples:**
{{range .Examples}}
- {{.}}
{{end}}
{{end}}
{{end}}
{{end}}
{{end}}
`

	tmpl.Parse(tmplText)
	return &MarkdownExporter{Template: tmpl}
}

// Export converts an FRMR document to Markdown
func (m *MarkdownExporter) Export(doc *FRMRDocument, w io.Writer) error {
	return m.Template.Execute(w, doc)
}

// FilterFRMR filters an FRMR document based on criteria
type FilterOptions struct {
	// Filter by impact levels
	ImpactLevels []string
	// Filter by specific KSI IDs
	KSIIDs []string
	// Filter by requirement IDs
	RequirementIDs []string
	// Include only specific types (FRD, FRR, FRA, KSI)
	Types []string
}

// FilterDocument creates a filtered copy of an FRMR document
func FilterDocument(doc *FRMRDocument, opts FilterOptions) *FRMRDocument {
	filtered := &FRMRDocument{
		Schema: doc.Schema,
		ID:     doc.ID + "-filtered",
		Info:   doc.Info,
		FRD:    make(map[string][]FRDItem),
		FRR:    make(map[string]FRRBase),
		FRA:    make(map[string][]FRAItem),
		KSI:    make(map[string]KSIItem),
	}

	// Check if type should be included
	includeType := func(t string) bool {
		if len(opts.Types) == 0 {
			return true
		}
		for _, typ := range opts.Types {
			if strings.EqualFold(typ, t) {
				return true
			}
		}
		return false
	}

	// Filter FRD
	if includeType("FRD") {
		filtered.FRD = doc.FRD // FRD doesn't have impact levels, so copy all
	}

	// Filter FRR
	if includeType("FRR") {
		for key, base := range doc.FRR {
			filteredBase := FRRBase{
				Base: struct {
					Application  string    `json:"application"`
					ID           string    `json:"id"`
					Requirements []FRRItem `json:"requirements"`
				}{
					Application: base.Base.Application,
					ID:          base.Base.ID,
					Requirements: []FRRItem{},
				},
			}

			for _, req := range base.Base.Requirements {
				include := true
				
				// Filter by impact levels
				if len(opts.ImpactLevels) > 0 && len(req.AppliedImpactLevels) > 0 {
					include = false
					for _, reqLevel := range req.AppliedImpactLevels {
						for _, filterLevel := range opts.ImpactLevels {
							if strings.EqualFold(reqLevel, filterLevel) {
								include = true
								break
							}
						}
						if include {
							break
						}
					}
				}

				// Filter by requirement IDs
				if include && len(opts.RequirementIDs) > 0 {
					include = false
					for _, id := range opts.RequirementIDs {
						if req.ID == id {
							include = true
							break
						}
					}
				}

				if include {
					filteredBase.Base.Requirements = append(filteredBase.Base.Requirements, req)
				}
			}

			if len(filteredBase.Base.Requirements) > 0 {
				filtered.FRR[key] = filteredBase
			}
		}
	}

	// Filter FRA
	if includeType("FRA") {
		filtered.FRA = doc.FRA // FRA doesn't have impact levels, so copy all
	}

	// Filter KSI
	if includeType("KSI") {
		for key, ksi := range doc.KSI {
			// Filter by KSI IDs
			if len(opts.KSIIDs) > 0 {
				include := false
				for _, id := range opts.KSIIDs {
					if ksi.ID == id {
						include = true
						break
					}
				}
				if !include {
					continue
				}
			}

			// Filter requirements by impact level
			filteredKSI := ksi
			if len(opts.ImpactLevels) > 0 {
				filteredReqs := []KSIRequirement{}
				for _, req := range ksi.Requirements {
					include := false
					for _, reqLevel := range req.AppliedImpactLevels {
						for _, filterLevel := range opts.ImpactLevels {
							if strings.EqualFold(reqLevel, filterLevel) {
								include = true
								break
							}
						}
						if include {
							break
						}
					}
					if include {
						filteredReqs = append(filteredReqs, req)
					}
				}
				filteredKSI.Requirements = filteredReqs
			}

			if len(filteredKSI.Requirements) > 0 {
				filtered.KSI[key] = filteredKSI
			}
		}
	}

	return filtered
}

// ValidateSchema validates an FRMR document against the FedRAMP schema
func ValidateSchema(doc *FRMRDocument) []string {
	errors := []string{}

	// Check required fields
	if doc.Info.Name == "" {
		errors = append(errors, "Info.Name is required")
	}
	if doc.Info.ShortName == "" {
		errors = append(errors, "Info.ShortName is required")
	}
	if doc.Info.CurrentRelease == "" {
		errors = append(errors, "Info.CurrentRelease is required")
	}
	if len(doc.Info.Types) == 0 {
		errors = append(errors, "Info.Types must contain at least one type")
	}

	// Validate KSI structure
	for id, ksi := range doc.KSI {
		if ksi.ID == "" {
			errors = append(errors, fmt.Sprintf("KSI %s: ID is required", id))
		}
		if ksi.Name == "" {
			errors = append(errors, fmt.Sprintf("KSI %s: Name is required", id))
		}
		if ksi.Indicator == "" {
			errors = append(errors, fmt.Sprintf("KSI %s: Indicator is required", id))
		}
		if len(ksi.Requirements) == 0 {
			errors = append(errors, fmt.Sprintf("KSI %s: Must have at least one requirement", id))
		}

		for i, req := range ksi.Requirements {
			if req.ID == "" {
				errors = append(errors, fmt.Sprintf("KSI %s Requirement %d: ID is required", id, i))
			}
			if req.Statement == "" {
				errors = append(errors, fmt.Sprintf("KSI %s Requirement %d: Statement is required", id, i))
			}
			if len(req.AppliedImpactLevels) == 0 {
				errors = append(errors, fmt.Sprintf("KSI %s Requirement %s: Must have at least one impact level", id, req.ID))
			}
		}
	}

	// Validate FRR structure
	for category, base := range doc.FRR {
		if base.Base.ID == "" {
			errors = append(errors, fmt.Sprintf("FRR %s: Base.ID is required", category))
		}
		if base.Base.Application == "" {
			errors = append(errors, fmt.Sprintf("FRR %s: Base.Application is required", category))
		}

		for i, req := range base.Base.Requirements {
			if req.ID == "" {
				errors = append(errors, fmt.Sprintf("FRR %s Requirement %d: ID is required", category, i))
			}
			if req.Statement == "" {
				errors = append(errors, fmt.Sprintf("FRR %s Requirement %d: Statement is required", category, i))
			}
		}
	}

	return errors
}

// GenerateEvidenceTemplate creates a template evidence file for KSI validation
func GenerateEvidenceTemplate(doc *FRMRDocument) map[string]interface{} {
	evidence := make(map[string]interface{})
	
	// Add metadata
	evidence["_metadata"] = map[string]interface{}{
		"generated_for": doc.Info.Name,
		"release":       doc.Info.CurrentRelease,
		"description":   "Template evidence file for KSI validation. Set each requirement ID to true/false based on implementation status.",
	}

	// Add all KSI requirements
	for _, ksi := range doc.KSI {
		for _, req := range ksi.Requirements {
			evidence[req.ID] = false // Default to false
		}
	}

	return evidence
}

// ExportEvidenceTemplate writes an evidence template to a writer
func ExportEvidenceTemplate(doc *FRMRDocument, w io.Writer) error {
	template := GenerateEvidenceTemplate(doc)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(template)
} 