package fedramp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// SCNType represents the type of significant change
type SCNType string

const (
	SCNAdaptive      SCNType = "adaptive"
	SCNTransformative SCNType = "transformative"
	SCNImpactChange   SCNType = "impact-change"
)

// SignificantChangeNotification represents an SCN as defined in RFC-0007
type SignificantChangeNotification struct {
	// Required fields for all SCNs
	ServiceOfferingID string    `json:"service_offering_id"`
	ThreePAOName      string    `json:"3pao_name,omitempty"`
	ChangeType        string    `json:"change_type"`
	RelatedPOAM       string    `json:"related_poam,omitempty"`
	ShortDescription  string    `json:"short_description"`
	ReasonForChange   string    `json:"reason_for_change"`
	ComponentsAffected []string `json:"components_affected"`
	ControlsAffected  []string `json:"controls_affected"`
	ImpactAnalysis    string    `json:"impact_analysis"`
	ApproverName      string    `json:"approver_name"`
	ApproverTitle     string    `json:"approver_title"`
	
	// SCN Type classification
	SCNType           SCNType   `json:"scn_type"`
	
	// Adaptive change specific fields
	DateOfChange      *time.Time `json:"date_of_change,omitempty"`
	VerificationSteps string     `json:"verification_steps,omitempty"`
	NewRisks          string     `json:"new_risks,omitempty"`
	
	// Transformative change specific fields (before)
	PlannedChangeDate *time.Time `json:"planned_change_date,omitempty"`
	RollbackPlan      string     `json:"rollback_plan,omitempty"`
	OptInRisk         string     `json:"opt_in_risk,omitempty"`
	HowToOptIn        string     `json:"how_to_opt_in,omitempty"`
	AssessmentPlan    string     `json:"assessment_plan,omitempty"`
	
	// Transformative change specific fields (after)
	AssessmentReport  string     `json:"assessment_report,omitempty"`
	
	// Metadata
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Status            string     `json:"status"` // draft, submitted, approved, etc.
}

// NewSCN creates a new Significant Change Notification
func NewSCN(serviceID, changeType, description, reason string) *SignificantChangeNotification {
	now := time.Now()
	return &SignificantChangeNotification{
		ServiceOfferingID: serviceID,
		ChangeType:        changeType,
		ShortDescription:  description,
		ReasonForChange:   reason,
		CreatedAt:         now,
		UpdatedAt:         now,
		Status:           "draft",
		ComponentsAffected: make([]string, 0),
		ControlsAffected:   make([]string, 0),
	}
}

// ClassifySCNType automatically determines the SCN type based on change characteristics
func (scn *SignificantChangeNotification) ClassifySCNType() error {
	// This is a simplified classification logic
	// In practice, this would be more sophisticated
	
	if scn.ChangeType == "impact-level-change" {
		scn.SCNType = SCNImpactChange
		return nil
	}
	
	// Check for transformative indicators in change type
	if scn.ChangeType == "new functionality" || scn.ChangeType == "major component" || 
	   scn.ChangeType == "architecture change" || scn.ChangeType == "new service" {
		scn.SCNType = SCNTransformative
		return nil
	}
	
	// Check for transformative indicators in description
	transformativeKeywords := []string{"new functionality", "major component", "architecture change", "new service", "adding", "introducing"}
	for _, keyword := range transformativeKeywords {
		if containsIgnoreCase(scn.ShortDescription, keyword) || containsIgnoreCase(scn.ReasonForChange, keyword) {
			scn.SCNType = SCNTransformative
			return nil
		}
	}
	
	// Default to adaptive
	scn.SCNType = SCNAdaptive
	return nil
}

// ValidateForSubmission checks if the SCN has all required fields for its type
func (scn *SignificantChangeNotification) ValidateForSubmission() error {
	// Common required fields
	if scn.ServiceOfferingID == "" {
		return fmt.Errorf("service offering ID is required")
	}
	if scn.ShortDescription == "" {
		return fmt.Errorf("short description is required")
	}
	if scn.ReasonForChange == "" {
		return fmt.Errorf("reason for change is required")
	}
	if scn.ApproverName == "" {
		return fmt.Errorf("approver name is required")
	}
	if scn.ApproverTitle == "" {
		return fmt.Errorf("approver title is required")
	}
	
	// Type-specific validation
	switch scn.SCNType {
	case SCNAdaptive:
		if scn.DateOfChange == nil {
			return fmt.Errorf("date of change is required for adaptive changes")
		}
		if scn.VerificationSteps == "" {
			return fmt.Errorf("verification steps are required for adaptive changes")
		}
		
	case SCNTransformative:
		if scn.PlannedChangeDate == nil {
			return fmt.Errorf("planned change date is required for transformative changes")
		}
		if scn.ThreePAOName == "" {
			return fmt.Errorf("3PAO name is required for transformative changes")
		}
		if scn.RollbackPlan == "" {
			return fmt.Errorf("rollback plan is required for transformative changes")
		}
		if scn.AssessmentPlan == "" {
			return fmt.Errorf("assessment plan is required for transformative changes")
		}
		
	case SCNImpactChange:
		return fmt.Errorf("impact categorization changes require reauthorization, not SCN")
	}
	
	return nil
}

// ToJSON exports the SCN as JSON for machine-readable format requirement
func (scn *SignificantChangeNotification) ToJSON() ([]byte, error) {
	return json.MarshalIndent(scn, "", "  ")
}

// FromJSON imports SCN from JSON
func (scn *SignificantChangeNotification) FromJSON(data []byte) error {
	return json.Unmarshal(data, scn)
}

// AddAffectedControl adds a control ID to the list of affected controls
func (scn *SignificantChangeNotification) AddAffectedControl(controlID string) {
	if !containsString(scn.ControlsAffected, controlID) {
		scn.ControlsAffected = append(scn.ControlsAffected, controlID)
		scn.UpdatedAt = time.Now()
	}
}

// AddAffectedComponent adds a component to the list of affected components
func (scn *SignificantChangeNotification) AddAffectedComponent(component string) {
	if !containsString(scn.ComponentsAffected, component) {
		scn.ComponentsAffected = append(scn.ComponentsAffected, component)
		scn.UpdatedAt = time.Now()
	}
}

// SetStatus updates the SCN status
func (scn *SignificantChangeNotification) SetStatus(status string) {
	scn.Status = status
	scn.UpdatedAt = time.Now()
}

// Helper functions
func contains(text, substring string) bool {
	return len(text) >= len(substring) && text[:len(substring)] == substring
}

func containsIgnoreCase(text, substring string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substring))
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SCNManager handles SCN operations
type SCNManager struct {
	notifications map[string]*SignificantChangeNotification
}

// NewSCNManager creates a new SCN manager
func NewSCNManager() *SCNManager {
	return &SCNManager{
		notifications: make(map[string]*SignificantChangeNotification),
	}
}

// AddNotification adds an SCN to the manager
func (mgr *SCNManager) AddNotification(id string, scn *SignificantChangeNotification) {
	mgr.notifications[id] = scn
}

// GetNotification retrieves an SCN by ID
func (mgr *SCNManager) GetNotification(id string) (*SignificantChangeNotification, bool) {
	scn, exists := mgr.notifications[id]
	return scn, exists
}

// ListNotificationsByService returns all SCNs for a service
func (mgr *SCNManager) ListNotificationsByService(serviceID string) []*SignificantChangeNotification {
	var result []*SignificantChangeNotification
	for _, scn := range mgr.notifications {
		if scn.ServiceOfferingID == serviceID {
			result = append(result, scn)
		}
	}
	return result
}

// GenerateSCNReport creates a summary report of all SCNs
func (mgr *SCNManager) GenerateSCNReport() map[string]interface{} {
	report := map[string]interface{}{
		"total_notifications": len(mgr.notifications),
		"by_type": map[string]int{
			"adaptive":      0,
			"transformative": 0,
			"impact_change": 0,
		},
		"by_status": map[string]int{},
	}
	
	for _, scn := range mgr.notifications {
		// Count by type
		switch scn.SCNType {
		case SCNAdaptive:
			report["by_type"].(map[string]int)["adaptive"]++
		case SCNTransformative:
			report["by_type"].(map[string]int)["transformative"]++
		case SCNImpactChange:
			report["by_type"].(map[string]int)["impact_change"]++
		}
		
		// Count by status
		statusCounts := report["by_status"].(map[string]int)
		statusCounts[scn.Status]++
	}
	
	return report
} 