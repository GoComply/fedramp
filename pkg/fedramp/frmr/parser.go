package frmr

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// FRMRDocument represents a FedRAMP Machine Readable document
type FRMRDocument struct {
	Schema string                      `json:"$schema"`
	ID     string                      `json:"$id"`
	Info   Info                        `json:"info"`
	FRD    map[string][]FRDItem        `json:"FRD,omitempty"`
	FRR    map[string]FRRBase          `json:"FRR,omitempty"`
	FRA    map[string][]FRAItem        `json:"FRA,omitempty"`
	KSI    map[string]KSIItem          `json:"KSI,omitempty"`
}

// Info contains metadata about the FRMR document
type Info struct {
	Name           string      `json:"name"`
	ShortName      string      `json:"short_name"`
	CurrentRelease string      `json:"current_release"`
	Types          []string    `json:"types"`
	Releases       []Release   `json:"releases"`
	FrontMatter    FrontMatter `json:"front_matter,omitempty"`
}

// Release represents a version release
type Release struct {
	ID              string           `json:"id"`
	PublishedDate   string           `json:"published_date"`
	Description     string           `json:"description"`
	PublicComment   bool             `json:"public_comment"`
	Effective       EffectiveInfo    `json:"effective"`
	RelatedRFCs     []RFC            `json:"related_rfcs,omitempty"`
}

// EffectiveInfo contains effectiveness information
type EffectiveInfo struct {
	TwentyX *TwentyXInfo `json:"20x,omitempty"`
	R5      *R5Info      `json:"r5,omitempty"`
}

// TwentyXInfo contains 20x specific information
type TwentyXInfo struct {
	Timeline        Timeline `json:"timeline"`
	SpecificRelease string   `json:"specific_release"`
	IsOptional      bool     `json:"is_optional"`
	Comment         string   `json:"comment"`
}

// Timeline represents implementation timeline
type Timeline struct {
	Pilot *PhaseInfo `json:"pilot,omitempty"`
}

// PhaseInfo contains phase-specific information
type PhaseInfo struct {
	StartDate  string `json:"start_date"`
	Designator string `json:"designator"`
	Comment    string `json:"comment"`
}

// R5Info contains R5 Balance specific information
type R5Info struct {
	Timeline        Timeline `json:"timeline"`
	SpecificRelease string   `json:"specific_release"`
	IsOptional      bool     `json:"is_optional"`
	Comment         string   `json:"comment"`
}

// RFC represents a Request for Comment
type RFC struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	ID            string `json:"id"`
	URL           string `json:"url"`
	DiscussionURL string `json:"discussion_url"`
	ShortName     string `json:"short_name"`
	FullName      string `json:"full_name"`
}

// FrontMatter contains authority and reference information
type FrontMatter struct {
	Authority []Authority `json:"authority,omitempty"`
}

// Authority represents regulatory authority
type Authority struct {
	Reference    string `json:"reference"`
	ReferenceURL string `json:"reference_url"`
	Description  string `json:"description"`
}

// FRDItem represents a FedRAMP Requirements Data item
type FRDItem struct {
	ID           string   `json:"id"`
	Term         string   `json:"term"`
	Definition   string   `json:"definition"`
	Reference    string   `json:"reference,omitempty"`
	ReferenceURL string   `json:"reference_url,omitempty"`
	Note         string   `json:"note,omitempty"`
	ReferencedFR []string `json:"referenced_fr,omitempty"`
}

// FRRBase represents the base structure for FRR requirements
type FRRBase struct {
	Base struct {
		Application  string    `json:"application"`
		ID           string    `json:"id"`
		Requirements []FRRItem `json:"requirements"`
	} `json:"base"`
}

// FRRItem represents a FedRAMP Requirements Rule item
type FRRItem struct {
	ID                   string     `json:"id"`
	Statement            string     `json:"statement"`
	Affects              []string   `json:"affects,omitempty"`
	PrimaryKeyWord       string     `json:"primary_key_word,omitempty"`
	AppliedImpactLevels  []string   `json:"applied_impact_levels,omitempty"`
	ReferencedRules      []string   `json:"referenced_rules,omitempty"`
	ReferencedFR         []string   `json:"referenced_fr,omitempty"`
	Frequency            string     `json:"frequency,omitempty"`
	FrequencyComment     string     `json:"frequency_comment,omitempty"`
	Condition            *Condition `json:"condition,omitempty"`
	AssessmentObjectives []string   `json:"assessment_objectives,omitempty"`
}

// Condition represents a conditional requirement
type Condition struct {
	AppliesWhen string `json:"applies_when"`
	Statement   string `json:"statement"`
}

// FRAItem represents a FedRAMP Requirements Assistance item
type FRAItem struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ReferencedFR []string `json:"referenced_fr,omitempty"`
	Examples     []string `json:"examples,omitempty"`
}

// KSIItem represents a Key Security Indicator
type KSIItem struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Indicator    string           `json:"indicator"`
	Requirements []KSIRequirement `json:"requirements"`
}

// KSIRequirement represents a specific KSI requirement
type KSIRequirement struct {
	ID                  string    `json:"id"`
	Statement           string    `json:"statement"`
	AppliedImpactLevels []string  `json:"applied_impact_levels"`
	Controls            []Control `json:"controls,omitempty"`
}

// Control represents a NIST control reference
type Control struct {
	ControlID string `json:"control_id"`
	Title     string `json:"title"`
}

// ParseFRMR parses a FRMR JSON document
func ParseFRMR(r io.Reader) (*FRMRDocument, error) {
	var doc FRMRDocument
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode FRMR document: %w", err)
	}
	return &doc, nil
}

// GetCurrentRelease returns the current release information
func (d *FRMRDocument) GetCurrentRelease() (*Release, error) {
	for _, release := range d.Info.Releases {
		if release.ID == d.Info.CurrentRelease {
			return &release, nil
		}
	}
	return nil, fmt.Errorf("current release %s not found", d.Info.CurrentRelease)
}

// GetKSIByID returns a specific KSI by its ID
func (d *FRMRDocument) GetKSIByID(id string) (*KSIItem, error) {
	for _, ksi := range d.KSI {
		if ksi.ID == id {
			return &ksi, nil
		}
	}
	return nil, fmt.Errorf("KSI %s not found", id)
}

// ValidateKSI validates if all KSI requirements are met
func (ksi *KSIItem) ValidateKSI(evidence map[string]bool) ValidationResult {
	result := ValidationResult{
		KSIID:        ksi.ID,
		KSIName:      ksi.Name,
		Requirements: make([]RequirementResult, 0),
	}

	for _, req := range ksi.Requirements {
		reqResult := RequirementResult{
			ID:        req.ID,
			Statement: req.Statement,
			Met:       evidence[req.ID],
		}
		result.Requirements = append(result.Requirements, reqResult)
		if !reqResult.Met {
			result.UnmetCount++
		}
	}

	result.TotalCount = len(ksi.Requirements)
	result.MetCount = result.TotalCount - result.UnmetCount
	result.FullyMet = result.UnmetCount == 0

	return result
}

// ValidationResult represents the result of KSI validation
type ValidationResult struct {
	KSIID        string
	KSIName      string
	FullyMet     bool
	TotalCount   int
	MetCount     int
	UnmetCount   int
	Requirements []RequirementResult
}

// RequirementResult represents the result of a single requirement
type RequirementResult struct {
	ID        string
	Statement string
	Met       bool
}

// GetEffectiveDate returns the effective date for a specific context
func (r *Release) GetEffectiveDate(context string) (*time.Time, error) {
	var dateStr string
	
	switch context {
	case "20x":
		if r.Effective.TwentyX != nil && r.Effective.TwentyX.Timeline.Pilot != nil {
			dateStr = r.Effective.TwentyX.Timeline.Pilot.StartDate
		}
	case "r5":
		if r.Effective.R5 != nil && r.Effective.R5.Timeline.Pilot != nil {
			dateStr = r.Effective.R5.Timeline.Pilot.StartDate
		}
	default:
		return nil, fmt.Errorf("unknown context: %s", context)
	}

	if dateStr == "" {
		return nil, fmt.Errorf("no effective date for context: %s", context)
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date %s: %w", dateStr, err)
	}

	return &t, nil
} 