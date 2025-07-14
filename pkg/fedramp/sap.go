// Package fedramp provides FedRAMP compliance automation tools
//
// Security Assessment Plan (SAP) - WORK IN PROGRESS
// Status: Basic structure implemented, integration pending
// TODO:
//   - Test case library for all controls
//   - Sampling methodology calculator
//   - Assessment schedule optimization
//   - Integration with MAS requirements
package fedramp

import (
	"encoding/json"
	"fmt"
	"time"
)

// SecurityAssessmentPlan represents a FedRAMP SAP
type SecurityAssessmentPlan struct {
	PlanID            string             `json:"plan_id"`
	ServiceOfferingID string             `json:"service_offering_id"`
	Version           string             `json:"version"`
	CreatedAt         time.Time          `json:"created_at"`
	AssessmentType    string             `json:"assessment_type"` // initial, annual, significant-change
	Scope             SAPAssessmentScope `json:"scope"`
	Schedule          AssessmentSchedule `json:"schedule"`
	TeamComposition   AssessmentTeam     `json:"team_composition"`
	Methodology       TestMethodology    `json:"methodology"`
	ControlSelection  ControlSelection   `json:"control_selection"`
	TestProcedures    []TestProcedure    `json:"test_procedures"`
	RulesOfEngagement RulesOfEngagement  `json:"rules_of_engagement"`
	Deliverables      []Deliverable      `json:"deliverables"`
}

// SAPAssessmentScope defines what will be assessed
type SAPAssessmentScope struct {
	SystemBoundary    SystemBoundary   `json:"system_boundary"`
	Locations         []string         `json:"locations"`
	DataCenters       []DataCenter     `json:"data_centers"`
	Components        []string         `json:"components"`
	ExcludedItems     []string         `json:"excluded_items"`
	SpecialConditions []string         `json:"special_conditions"`
}

// SystemBoundary defines the authorization boundary
type SystemBoundary struct {
	Description      string   `json:"description"`
	DiagramReference string   `json:"diagram_reference"`
	IPRanges         []string `json:"ip_ranges"`
	DNSDomains       []string `json:"dns_domains"`
}

// DataCenter represents a data center location
type DataCenter struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Type     string `json:"type"` // Primary, Secondary, DR
	Provider string `json:"provider"`
}

// AssessmentSchedule defines the assessment timeline
type AssessmentSchedule struct {
	StartDate         time.Time       `json:"start_date"`
	EndDate           time.Time       `json:"end_date"`
	KeyMilestones     []KeyMilestone  `json:"key_milestones"`
	BlackoutDates     []DateRange     `json:"blackout_dates"`
	AssessmentPhases  []Phase         `json:"assessment_phases"`
}

// KeyMilestone represents important dates
type KeyMilestone struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

// DateRange represents a date range
type DateRange struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Reason string    `json:"reason"`
}

// Phase represents an assessment phase
type Phase struct {
	Name        string    `json:"name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Activities  []string  `json:"activities"`
	Deliverable string    `json:"deliverable"`
}

// TestMethodology describes how testing will be performed
type TestMethodology struct {
	Framework          string              `json:"framework"` // NIST SP 800-53A Rev 5
	TestingLevels      map[string]string   `json:"testing_levels"` // Control: Level
	SamplingStrategy   SamplingStrategy    `json:"sampling_strategy"`
	TestingTools       []TestingTool       `json:"testing_tools"`
	EvidenceCollection EvidenceRequirements `json:"evidence_collection"`
}

// SamplingStrategy defines how samples will be selected
type SamplingStrategy struct {
	Approach     string            `json:"approach"` // Statistical, Judgmental, Fixed
	SampleSizes  map[string]int    `json:"sample_sizes"`
	Confidence   string            `json:"confidence_level"`
	Rationale    string            `json:"rationale"`
}

// TestingTool represents tools to be used
type TestingTool struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Purpose  string   `json:"purpose"`
	License  string   `json:"license"`
	Approved bool     `json:"approved_by_csp"`
}

// EvidenceRequirements defines evidence collection requirements
type EvidenceRequirements struct {
	Screenshots      bool     `json:"screenshots_required"`
	Logs             bool     `json:"logs_required"`
	Configurations   bool     `json:"configurations_required"`
	Interviews       bool     `json:"interviews_required"`
	RetentionPeriod  string   `json:"retention_period"`
	HandlingGuidance string   `json:"handling_guidance"`
}

// ControlSelection defines which controls will be tested
type ControlSelection struct {
	Baseline          string              `json:"baseline"` // Low, Moderate, High
	TotalControls     int                 `json:"total_controls"`
	SelectedControls  []SelectedControl   `json:"selected_controls"`
	SelectionRationale string             `json:"selection_rationale"`
	RiskBasedApproach  bool               `json:"risk_based_approach"`
}

// SelectedControl represents a control selected for testing
type SelectedControl struct {
	ControlID      string   `json:"control_id"`
	TestDepth      string   `json:"test_depth"` // Basic, Focused, Comprehensive
	TestMethods    []string `json:"test_methods"` // Examine, Interview, Test
	Justification  string   `json:"justification,omitempty"`
}

// TestProcedure defines specific test procedures
type TestProcedure struct {
	ProcedureID    string   `json:"procedure_id"`
	ControlID      string   `json:"control_id"`
	Objective      string   `json:"objective"`
	TestSteps      []string `json:"test_steps"`
	ExpectedResult string   `json:"expected_result"`
	TestData       string   `json:"test_data"`
	Prerequisites  []string `json:"prerequisites"`
	Duration       string   `json:"estimated_duration"`
}

// RulesOfEngagement defines assessment rules
type RulesOfEngagement struct {
	AuthorizedActivities   []string          `json:"authorized_activities"`
	ProhibitedActivities   []string          `json:"prohibited_activities"`
	CommunicationProtocol  CommunicationPlan `json:"communication_protocol"`
	EscalationProcedure    []EscalationStep  `json:"escalation_procedure"`
	IncidentResponse       string            `json:"incident_response"`
	DataHandling           DataHandlingRules `json:"data_handling"`
}

// CommunicationPlan defines how communication will occur
type CommunicationPlan struct {
	PrimaryPOC      Contact  `json:"primary_poc"`
	SecondaryPOC    Contact  `json:"secondary_poc"`
	StatusReports   string   `json:"status_report_frequency"`
	MeetingSchedule string   `json:"meeting_schedule"`
	Channels        []string `json:"communication_channels"`
}

// Contact represents a point of contact
type Contact struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// EscalationStep defines escalation procedures
type EscalationStep struct {
	Level       int    `json:"level"`
	Trigger     string `json:"trigger"`
	ContactRole string `json:"contact_role"`
	Timeframe   string `json:"timeframe"`
}

// DataHandlingRules defines how data will be handled
type DataHandlingRules struct {
	Classification    string   `json:"classification"`
	Storage           string   `json:"storage_requirements"`
	Transmission      string   `json:"transmission_requirements"`
	Retention         string   `json:"retention_period"`
	Destruction       string   `json:"destruction_method"`
	AccessRestriction []string `json:"access_restrictions"`
}

// Deliverable represents an assessment deliverable
type Deliverable struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Format      string    `json:"format"`
	Recipients  []string  `json:"recipients"`
}

// NewSecurityAssessmentPlan creates a new SAP
func NewSecurityAssessmentPlan(serviceID, assessmentType string) *SecurityAssessmentPlan {
	return &SecurityAssessmentPlan{
		PlanID:            fmt.Sprintf("SAP-%s-%d", serviceID, time.Now().Unix()),
		ServiceOfferingID: serviceID,
		Version:           "1.0",
		CreatedAt:         time.Now(),
		AssessmentType:    assessmentType,
		TestProcedures:    make([]TestProcedure, 0),
		Deliverables:      make([]Deliverable, 0),
	}
}

// AddTestProcedure adds a test procedure to the SAP
func (sap *SecurityAssessmentPlan) AddTestProcedure(procedure TestProcedure) {
	procedure.ProcedureID = fmt.Sprintf("TP-%s-%d", procedure.ControlID, len(sap.TestProcedures)+1)
	sap.TestProcedures = append(sap.TestProcedures, procedure)
}

// ToJSON exports the SAP as JSON
func (sap *SecurityAssessmentPlan) ToJSON() ([]byte, error) {
	return json.MarshalIndent(sap, "", "  ")
}

// GenerateTestProcedures creates standard test procedures for controls
func GenerateTestProcedures(controlIDs []string) []TestProcedure {
	procedures := make([]TestProcedure, 0)
	
	for _, controlID := range controlIDs {
		procedure := TestProcedure{
			ControlID: controlID,
			Objective: fmt.Sprintf("Verify implementation and effectiveness of control %s", controlID),
			TestSteps: []string{
				"Review control implementation documentation",
				"Interview control owners and operators",
				"Examine evidence of control operation",
				"Test control effectiveness through sampling",
				"Document findings and exceptions",
			},
			ExpectedResult: "Control is implemented as designed and operating effectively",
			Duration:       "2-4 hours",
		}
		procedures = append(procedures, procedure)
	}
	
	return procedures
} 