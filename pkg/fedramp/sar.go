// Package fedramp provides FedRAMP compliance automation tools
//
// Security Assessment Report (SAR) - WORK IN PROGRESS
// Status: Basic structure implemented, integration pending
// TODO: 
//   - Integration with assessment tools
//   - Evidence collection automation
//   - OSCAL SAR format support
//   - Report generation templates
package fedramp

import (
	"encoding/json"
	"fmt"
	"time"
)

// SecurityAssessmentReport represents a FedRAMP SAR
type SecurityAssessmentReport struct {
	ReportID          string              `json:"report_id"`
	ServiceOfferingID string              `json:"service_offering_id"`
	AssessmentType    string              `json:"assessment_type"` // initial, annual, significant-change
	GeneratedAt       time.Time           `json:"generated_at"`
	AssessmentPeriod  AssessmentPeriod    `json:"assessment_period"`
	ExecutiveSummary  ExecutiveSummary    `json:"executive_summary"`
	AssessmentTeam    AssessmentTeam      `json:"assessment_team"`
	Methodology       SARAssessmentMethod `json:"methodology"`
	ControlFindings   []ControlFinding    `json:"control_findings"`
	RiskSummary       RiskSummary         `json:"risk_summary"`
	Recommendations   []Recommendation    `json:"recommendations"`
	TestCases         []TestCase          `json:"test_cases"`
	Evidence          []SARAssessmentEvidence `json:"evidence"`
	ThreePAOStatement ThreePAOStatement   `json:"3pao_statement"`
}

// AssessmentPeriod defines the timeframe of the assessment
type AssessmentPeriod struct {
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	OnSiteDays      int       `json:"onsite_days"`
	RemoteDays      int       `json:"remote_days"`
	TotalAssessors  int       `json:"total_assessors"`
}

// ExecutiveSummary provides high-level assessment results
type ExecutiveSummary struct {
	OverallRisk       string   `json:"overall_risk"` // Low, Moderate, High
	ComplianceStatus  string   `json:"compliance_status"` // Compliant, Compliant with Findings, Non-Compliant
	KeyFindings       []string `json:"key_findings"`
	CriticalFindings  int      `json:"critical_findings"`
	HighFindings      int      `json:"high_findings"`
	ModerateFindings  int      `json:"moderate_findings"`
	LowFindings       int      `json:"low_findings"`
	RecommendedAction string   `json:"recommended_action"` // ATO, ATO with conditions, Denial
}

// AssessmentTeam documents the assessment team members
type AssessmentTeam struct {
	LeadAssessor    TeamMember   `json:"lead_assessor"`
	TeamMembers     []TeamMember `json:"team_members"`
	ThreePAOName    string       `json:"3pao_name"`
	IndependenceStatement string  `json:"independence_statement"`
}

// TeamMember represents an assessment team member
type TeamMember struct {
	Name            string   `json:"name"`
	Role            string   `json:"role"`
	Qualifications  []string `json:"qualifications"`
	YearsExperience int      `json:"years_experience"`
}

// SARAssessmentMethod describes the assessment methodology
type SARAssessmentMethod struct {
	Framework        string   `json:"framework"` // NIST SP 800-53A
	SamplingApproach string   `json:"sampling_approach"`
	TestingMethods   []string `json:"testing_methods"` // interview, examine, test
	ToolsUsed        []string `json:"tools_used"`
	Limitations      []string `json:"limitations"`
}

// ControlFinding represents findings for a specific control
type ControlFinding struct {
	ControlID        string    `json:"control_id"`
	ControlTitle     string    `json:"control_title"`
	FindingID        string    `json:"finding_id"`
	Severity         string    `json:"severity"` // Critical, High, Moderate, Low
	Status           string    `json:"status"` // Satisfied, Other Than Satisfied, Not Applicable
	Description      string    `json:"description"`
	Evidence         []string  `json:"evidence"`
	RootCause        string    `json:"root_cause"`
	Impact           string    `json:"impact"`
	Likelihood       string    `json:"likelihood"`
	RiskRating       string    `json:"risk_rating"`
	Recommendation   string    `json:"recommendation"`
	CSPResponse      string    `json:"csp_response,omitempty"`
	RemediationPlan  string    `json:"remediation_plan,omitempty"`
	TestDate         time.Time `json:"test_date"`
	Tester           string    `json:"tester"`
}

// TestCase documents specific test procedures
type TestCase struct {
	TestID          string    `json:"test_id"`
	ControlID       string    `json:"control_id"`
	TestObjective   string    `json:"test_objective"`
	TestProcedure   string    `json:"test_procedure"`
	ExpectedResult  string    `json:"expected_result"`
	ActualResult    string    `json:"actual_result"`
	TestEvidence    []string  `json:"test_evidence"`
	PassFail        string    `json:"pass_fail"`
	TestDate        time.Time `json:"test_date"`
	TesterName      string    `json:"tester_name"`
}

// RiskSummary provides overall risk assessment
type RiskSummary struct {
	TotalRisk           string            `json:"total_risk"`
	RiskByCategory      map[string]int    `json:"risk_by_category"`
	TrendAnalysis       string            `json:"trend_analysis"`
	ComparisonPrevious  string            `json:"comparison_to_previous"`
	SystematicIssues    []string          `json:"systematic_issues"`
}

// Recommendation provides actionable recommendations
type Recommendation struct {
	ID          string `json:"id"`
	Priority    string `json:"priority"` // Critical, High, Medium, Low
	Category    string `json:"category"`
	Description string `json:"description"`
	Benefit     string `json:"benefit"`
	Effort      string `json:"effort"` // Low, Medium, High
	Timeline    string `json:"timeline"`
}

// SARAssessmentEvidence documents evidence collected
type SARAssessmentEvidence struct {
	EvidenceID   string    `json:"evidence_id"`
	Type         string    `json:"type"` // screenshot, document, interview, observation
	Description  string    `json:"description"`
	ControlIDs   []string  `json:"control_ids"`
	CollectedBy  string    `json:"collected_by"`
	CollectedAt  time.Time `json:"collected_at"`
	Location     string    `json:"location"` // file path or reference
	Hash         string    `json:"hash,omitempty"`
}

// ThreePAOStatement provides the 3PAO attestation
type ThreePAOStatement struct {
	Statement       string    `json:"statement"`
	SignedBy        string    `json:"signed_by"`
	Title           string    `json:"title"`
	Date            time.Time `json:"date"`
	ThreePAOLogo    string    `json:"3pao_logo,omitempty"`
}

// NewSecurityAssessmentReport creates a new SAR
func NewSecurityAssessmentReport(serviceID, assessmentType string) *SecurityAssessmentReport {
	return &SecurityAssessmentReport{
		ReportID:          fmt.Sprintf("SAR-%s-%d", serviceID, time.Now().Unix()),
		ServiceOfferingID: serviceID,
		AssessmentType:    assessmentType,
		GeneratedAt:       time.Now(),
		ControlFindings:   make([]ControlFinding, 0),
		TestCases:         make([]TestCase, 0),
		Evidence:          make([]SARAssessmentEvidence, 0),
		Recommendations:   make([]Recommendation, 0),
	}
}

// AddControlFinding adds a finding to the SAR
func (sar *SecurityAssessmentReport) AddControlFinding(finding ControlFinding) {
	finding.FindingID = fmt.Sprintf("FIND-%s-%d", finding.ControlID, len(sar.ControlFindings)+1)
	sar.ControlFindings = append(sar.ControlFindings, finding)
	sar.updateSummary()
}

// updateSummary recalculates the executive summary
func (sar *SecurityAssessmentReport) updateSummary() {
	critical, high, moderate, low := 0, 0, 0, 0
	
	for _, finding := range sar.ControlFindings {
		if finding.Status == "Other Than Satisfied" {
			switch finding.Severity {
			case "Critical":
				critical++
			case "High":
				high++
			case "Moderate":
				moderate++
			case "Low":
				low++
			}
		}
	}
	
	sar.ExecutiveSummary.CriticalFindings = critical
	sar.ExecutiveSummary.HighFindings = high
	sar.ExecutiveSummary.ModerateFindings = moderate
	sar.ExecutiveSummary.LowFindings = low
	
	// Determine overall risk
	if critical > 0 {
		sar.ExecutiveSummary.OverallRisk = "High"
		sar.ExecutiveSummary.ComplianceStatus = "Non-Compliant"
	} else if high > 3 {
		sar.ExecutiveSummary.OverallRisk = "High"
		sar.ExecutiveSummary.ComplianceStatus = "Compliant with Findings"
	} else if high > 0 || moderate > 10 {
		sar.ExecutiveSummary.OverallRisk = "Moderate"
		sar.ExecutiveSummary.ComplianceStatus = "Compliant with Findings"
	} else {
		sar.ExecutiveSummary.OverallRisk = "Low"
		sar.ExecutiveSummary.ComplianceStatus = "Compliant"
	}
}

// ToJSON exports the SAR as JSON
func (sar *SecurityAssessmentReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(sar, "", "  ")
}

// GenerateTestCase creates a test case for a control
func GenerateTestCase(controlID, objective string) TestCase {
	return TestCase{
		TestID:        fmt.Sprintf("TEST-%s-%d", controlID, time.Now().UnixNano()),
		ControlID:     controlID,
		TestObjective: objective,
		TestDate:      time.Now(),
	}
} 