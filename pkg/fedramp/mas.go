package fedramp

import (
	"encoding/json"
	"fmt"
	"time"
)

// MASAssessmentType represents the type of assessment
type MASAssessmentType string

const (
	MASInitial      MASAssessmentType = "initial"
	MASAnnual       MASAssessmentType = "annual"
	MASSignificant  MASAssessmentType = "significant_change"
	MASIncident     MASAssessmentType = "incident_response"
)

// MASRequirement represents a minimum assessment requirement
type MASRequirement struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	AssessmentType   MASAssessmentType `json:"assessment_type"`
	Frequency        string            `json:"frequency"`
	RequiredMethods  []string          `json:"required_methods"`
	MinimumEvidence  []string          `json:"minimum_evidence"`
	RelatedControls  []string          `json:"related_controls"`
	AutomationLevel  string            `json:"automation_level"` // full, partial, manual
}

// MASAssessment represents an assessment following MAS
type MASAssessment struct {
	AssessmentID      string                     `json:"assessment_id"`
	ServiceOfferingID string                     `json:"service_offering_id"`
	AssessmentType    MASAssessmentType          `json:"assessment_type"`
	StartDate         time.Time                  `json:"start_date"`
	EndDate           *time.Time                 `json:"end_date,omitempty"`
	Status            string                     `json:"status"` // planned, in_progress, complete
	ThreePAO          AssessmentOrganization     `json:"3pao"`
	Scope             AssessmentScope            `json:"scope"`
	Methods           []AssessmentMethod         `json:"methods"`
	Findings          []AssessmentFinding        `json:"findings"`
	Evidence          []AssessmentEvidence       `json:"evidence"`
	Attestation       *AssessmentAttestation     `json:"attestation,omitempty"`
}

// AssessmentOrganization represents the 3PAO or assessment org
type AssessmentOrganization struct {
	Name            string   `json:"name"`
	A2LAAccreditation string `json:"a2la_accreditation"`
	LeadAssessor    string   `json:"lead_assessor"`
	TeamMembers     []string `json:"team_members"`
	ContactEmail    string   `json:"contact_email"`
}

// AssessmentScope defines what's being assessed
type AssessmentScope struct {
	FullAssessment      bool     `json:"full_assessment"`
	SampledControls     []string `json:"sampled_controls,omitempty"`
	ExcludedComponents  []string `json:"excluded_components,omitempty"`
	Locations           []string `json:"locations"`
	DataTypes           []string `json:"data_types"`
	UserPopulation      int      `json:"user_population"`
}

// AssessmentMethod describes how assessment is performed
type AssessmentMethod struct {
	MethodType      string    `json:"method_type"` // interview, examine, test
	Description     string    `json:"description"`
	Automated       bool      `json:"automated"`
	ToolsUsed       []string  `json:"tools_used,omitempty"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	ControlsCovered []string  `json:"controls_covered"`
}

// AssessmentFinding represents a finding from the assessment
type AssessmentFinding struct {
	FindingID       string    `json:"finding_id"`
	ControlID       string    `json:"control_id"`
	Severity        string    `json:"severity"` // high, moderate, low
	Status          string    `json:"status"` // open, remediated, risk_accepted
	Description     string    `json:"description"`
	Recommendation  string    `json:"recommendation"`
	Evidence        []string  `json:"evidence_refs"`
	DateIdentified  time.Time `json:"date_identified"`
	DateRemediated  *time.Time `json:"date_remediated,omitempty"`
}

// AssessmentEvidence represents evidence collected
type AssessmentEvidence struct {
	EvidenceID      string    `json:"evidence_id"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	CollectionDate  time.Time `json:"collection_date"`
	CollectedBy     string    `json:"collected_by"`
	Location        string    `json:"location"`
	ControlsCovered []string  `json:"controls_covered"`
	Automated       bool      `json:"automated"`
}

// AssessmentAttestation represents 3PAO attestation
type AssessmentAttestation struct {
	AttestorName    string    `json:"attestor_name"`
	AttestorTitle   string    `json:"attestor_title"`
	Organization    string    `json:"organization"`
	Date            time.Time `json:"date"`
	Statement       string    `json:"statement"`
	Signature       string    `json:"signature"`
}

// MASRequirements defines standard requirements
var MASRequirements = map[MASAssessmentType][]MASRequirement{
	MASInitial: {
		{
			ID:              "MAS-INIT-001",
			Name:            "Full Control Assessment",
			Description:     "Complete assessment of all controls in the baseline",
			AssessmentType:  MASInitial,
			Frequency:       "once",
			RequiredMethods: []string{"examine", "interview", "test"},
			MinimumEvidence: []string{"screenshots", "configurations", "test_results", "policies"},
			RelatedControls: []string{"all"},
			AutomationLevel: "partial",
		},
		{
			ID:              "MAS-INIT-002",
			Name:            "Vulnerability Scanning",
			Description:     "Authenticated vulnerability scanning of all components",
			AssessmentType:  MASInitial,
			Frequency:       "once",
			RequiredMethods: []string{"test"},
			MinimumEvidence: []string{"scan_reports", "remediation_evidence"},
			RelatedControls: []string{"RA-5", "SI-2"},
			AutomationLevel: "full",
		},
		{
			ID:              "MAS-INIT-003",
			Name:            "Penetration Testing",
			Description:     "Application and infrastructure penetration testing",
			AssessmentType:  MASInitial,
			Frequency:       "once",
			RequiredMethods: []string{"test"},
			MinimumEvidence: []string{"pentest_report", "remediation_evidence"},
			RelatedControls: []string{"CA-8"},
			AutomationLevel: "manual",
		},
	},
	MASAnnual: {
		{
			ID:              "MAS-ANN-001",
			Name:            "Annual Assessment",
			Description:     "Assessment of one-third of controls plus all high-risk",
			AssessmentType:  MASAnnual,
			Frequency:       "annual",
			RequiredMethods: []string{"examine", "interview", "test"},
			MinimumEvidence: []string{"control_evidence", "test_results"},
			RelatedControls: []string{"subset"},
			AutomationLevel: "partial",
		},
		{
			ID:              "MAS-ANN-002",
			Name:            "Continuous Monitoring Review",
			Description:     "Review of continuous monitoring data and POA&Ms",
			AssessmentType:  MASAnnual,
			Frequency:       "annual",
			RequiredMethods: []string{"examine"},
			MinimumEvidence: []string{"conmon_reports", "poam_status"},
			RelatedControls: []string{"CA-7", "PM-4"},
			AutomationLevel: "full",
		},
	},
}

// NewMASAssessment creates a new MAS assessment
func NewMASAssessment(serviceID string, assessmentType MASAssessmentType) *MASAssessment {
	return &MASAssessment{
		AssessmentID:      fmt.Sprintf("MAS-%s-%s-%s", serviceID, assessmentType, time.Now().Format("20060102")),
		ServiceOfferingID: serviceID,
		AssessmentType:    assessmentType,
		StartDate:         time.Now(),
		Status:           "planned",
		Methods:          make([]AssessmentMethod, 0),
		Findings:         make([]AssessmentFinding, 0),
		Evidence:         make([]AssessmentEvidence, 0),
	}
}

// AddMethod adds an assessment method
func (a *MASAssessment) AddMethod(method AssessmentMethod) {
	a.Methods = append(a.Methods, method)
}

// AddFinding adds an assessment finding
func (a *MASAssessment) AddFinding(finding AssessmentFinding) {
	a.Findings = append(a.Findings, finding)
}

// AddEvidence adds assessment evidence
func (a *MASAssessment) AddEvidence(evidence AssessmentEvidence) {
	a.Evidence = append(a.Evidence, evidence)
}

// Complete marks the assessment as complete
func (a *MASAssessment) Complete() {
	now := time.Now()
	a.EndDate = &now
	a.Status = "complete"
}

// GetRequirements returns the MAS requirements for this assessment type
func (a *MASAssessment) GetRequirements() []MASRequirement {
	return MASRequirements[a.AssessmentType]
}

// ValidateCompleteness checks if assessment meets MAS requirements
func (a *MASAssessment) ValidateCompleteness() error {
	requirements := a.GetRequirements()
	
	for _, req := range requirements {
		// Check if required methods were used
		methodsUsed := make(map[string]bool)
		for _, method := range a.Methods {
			methodsUsed[method.MethodType] = true
		}
		
		for _, reqMethod := range req.RequiredMethods {
			if !methodsUsed[reqMethod] {
				return fmt.Errorf("missing required assessment method: %s for %s", reqMethod, req.Name)
			}
		}
		
		// Check if minimum evidence types are present
		evidenceTypes := make(map[string]bool)
		for _, evidence := range a.Evidence {
			evidenceTypes[evidence.Type] = true
		}
		
		for _, reqEvidence := range req.MinimumEvidence {
			if !evidenceTypes[reqEvidence] {
				return fmt.Errorf("missing required evidence type: %s for %s", reqEvidence, req.Name)
			}
		}
	}
	
	// Check if assessment has findings (even if none found)
	if a.Status == "complete" && a.Attestation == nil {
		return fmt.Errorf("completed assessment requires attestation")
	}
	
	return nil
}

// GenerateSummary creates an assessment summary
func (a *MASAssessment) GenerateSummary() map[string]interface{} {
	findingsBySeverity := make(map[string]int)
	for _, finding := range a.Findings {
		findingsBySeverity[finding.Severity]++
	}
	
	methodTypes := make(map[string]int)
	for _, method := range a.Methods {
		methodTypes[method.MethodType]++
	}
	
	automatedEvidence := 0
	manualEvidence := 0
	for _, evidence := range a.Evidence {
		if evidence.Automated {
			automatedEvidence++
		} else {
			manualEvidence++
		}
	}
	
	duration := ""
	if a.EndDate != nil {
		duration = a.EndDate.Sub(a.StartDate).String()
	}
	
	return map[string]interface{}{
		"assessment_id":       a.AssessmentID,
		"type":               a.AssessmentType,
		"status":             a.Status,
		"duration":           duration,
		"total_findings":     len(a.Findings),
		"findings_by_severity": findingsBySeverity,
		"methods_used":       methodTypes,
		"evidence_collected": len(a.Evidence),
		"automated_evidence": automatedEvidence,
		"manual_evidence":    manualEvidence,
		"3pao":              a.ThreePAO.Name,
	}
}

// ToJSON exports the assessment as JSON
func (a *MASAssessment) ToJSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
} 