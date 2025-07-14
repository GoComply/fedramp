// Package fedramp provides FedRAMP compliance automation tools
//
// Plan of Action and Milestones (POA&M) - WORK IN PROGRESS
// Status: Basic structure implemented, integration pending
// TODO:
//   - Integration with ConMon findings
//   - Automated risk scoring
//   - FedRAMP POA&M template generation
//   - Deviation request handling
package fedramp

import (
	"encoding/json"
	"fmt"
	"time"
)

// PlanOfActionMilestones represents a FedRAMP POA&M
type PlanOfActionMilestones struct {
	DocumentID        string       `json:"document_id"`
	ServiceOfferingID string       `json:"service_offering_id"`
	GeneratedAt       time.Time    `json:"generated_at"`
	LastUpdated       time.Time    `json:"last_updated"`
	POAMItems         []POAMItem   `json:"poam_items"`
	Summary           POAMSummary  `json:"summary"`
	RiskAdjustment    RiskAdjustment `json:"risk_adjustment"`
}

// POAMItem represents an individual POA&M entry
type POAMItem struct {
	ItemID              string    `json:"item_id"`
	FindingID           string    `json:"finding_id"`
	ControlID           string    `json:"control_id"`
	Weakness            string    `json:"weakness"`
	Severity            string    `json:"severity"` // Critical, High, Moderate, Low
	RawRisk             string    `json:"raw_risk"`
	Status              string    `json:"status"` // Open, Ongoing, Risk Accepted, Completed, Cancelled
	ResponsibleParty    string    `json:"responsible_party"`
	Resources           string    `json:"resources"`
	MilestoneDates      []POAMMilestone `json:"milestone_dates"`
	IdentifiedDate      time.Time `json:"identified_date"`
	PlannedCompletion   time.Time `json:"planned_completion"`
	ActualCompletion    *time.Time `json:"actual_completion,omitempty"`
	Comments            string    `json:"comments"`
	RemediationPlan     string    `json:"remediation_plan"`
	MitigatingFactors   string    `json:"mitigating_factors,omitempty"`
	ResidualRisk        string    `json:"residual_risk"`
	Source              string    `json:"source"` // SAR, ConMon, Incident, Scan
	VendorDependency    bool      `json:"vendor_dependency"`
	FalsePositive       bool      `json:"false_positive"`
	OperationalRequirement bool   `json:"operational_requirement"`
}

// POAMMilestone represents a milestone in the POA&M
type POAMMilestone struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// POAMSummary provides summary statistics
type POAMSummary struct {
	TotalItems        int            `json:"total_items"`
	OpenItems         int            `json:"open_items"`
	OverdueItems      int            `json:"overdue_items"`
	ItemsBySeverity   map[string]int `json:"items_by_severity"`
	ItemsByStatus     map[string]int `json:"items_by_status"`
	AverageAge        float64        `json:"average_age_days"`
	OldestItem        time.Time      `json:"oldest_item_date"`
	CompletionRate    float64        `json:"completion_rate"`
	ProjectedClosure  time.Time      `json:"projected_closure"`
}

// RiskAdjustment tracks risk acceptance and mitigation
type RiskAdjustment struct {
	AcceptedRisks     []RiskAcceptance `json:"accepted_risks"`
	MitigatedRisks    []RiskMitigation `json:"mitigated_risks"`
	TotalRiskScore    float64          `json:"total_risk_score"`
	AdjustedRiskScore float64          `json:"adjusted_risk_score"`
}

// RiskAcceptance documents accepted risks
type RiskAcceptance struct {
	ItemID           string    `json:"item_id"`
	AcceptanceDate   time.Time `json:"acceptance_date"`
	AcceptedBy       string    `json:"accepted_by"`
	Justification    string    `json:"justification"`
	ReviewDate       time.Time `json:"review_date"`
	ExpirationDate   time.Time `json:"expiration_date"`
}

// RiskMitigation documents risk mitigation measures
type RiskMitigation struct {
	ItemID              string    `json:"item_id"`
	MitigationStrategy  string    `json:"mitigation_strategy"`
	ImplementationDate  time.Time `json:"implementation_date"`
	Effectiveness       string    `json:"effectiveness"` // High, Medium, Low
	ResidualRisk        string    `json:"residual_risk"`
}

// NewPOAM creates a new POA&M document
func NewPOAM(serviceID string) *PlanOfActionMilestones {
	return &PlanOfActionMilestones{
		DocumentID:        fmt.Sprintf("POAM-%s-%d", serviceID, time.Now().Unix()),
		ServiceOfferingID: serviceID,
		GeneratedAt:       time.Now(),
		LastUpdated:       time.Now(),
		POAMItems:         make([]POAMItem, 0),
		Summary:           POAMSummary{ItemsBySeverity: make(map[string]int), ItemsByStatus: make(map[string]int)},
	}
}

// AddItem adds a new POA&M item
func (poam *PlanOfActionMilestones) AddItem(item POAMItem) {
	item.ItemID = fmt.Sprintf("POAM-%d", len(poam.POAMItems)+1)
	poam.POAMItems = append(poam.POAMItems, item)
	poam.LastUpdated = time.Now()
	poam.updateSummary()
}

// UpdateItem updates an existing POA&M item
func (poam *PlanOfActionMilestones) UpdateItem(itemID string, updates map[string]interface{}) error {
	for i, item := range poam.POAMItems {
		if item.ItemID == itemID {
			// Apply updates (simplified for brevity)
			if status, ok := updates["status"].(string); ok {
				poam.POAMItems[i].Status = status
			}
			if completion, ok := updates["actual_completion"].(time.Time); ok {
				poam.POAMItems[i].ActualCompletion = &completion
			}
			poam.LastUpdated = time.Now()
			poam.updateSummary()
			return nil
		}
	}
	return fmt.Errorf("POA&M item %s not found", itemID)
}

// updateSummary recalculates summary statistics
func (poam *PlanOfActionMilestones) updateSummary() {
	summary := &poam.Summary
	summary.TotalItems = len(poam.POAMItems)
	summary.OpenItems = 0
	summary.OverdueItems = 0
	summary.ItemsBySeverity = make(map[string]int)
	summary.ItemsByStatus = make(map[string]int)
	
	now := time.Now()
	totalAge := 0.0
	completed := 0
	
	for _, item := range poam.POAMItems {
		// Count by status
		summary.ItemsByStatus[item.Status]++
		
		// Count by severity
		summary.ItemsBySeverity[item.Severity]++
		
		// Count open items
		if item.Status == "Open" || item.Status == "Ongoing" {
			summary.OpenItems++
			
			// Check if overdue
			if item.PlannedCompletion.Before(now) {
				summary.OverdueItems++
			}
		}
		
		// Track completed
		if item.Status == "Completed" {
			completed++
		}
		
		// Calculate age
		age := now.Sub(item.IdentifiedDate).Hours() / 24
		totalAge += age
		
		// Track oldest
		if summary.OldestItem.IsZero() || item.IdentifiedDate.Before(summary.OldestItem) {
			summary.OldestItem = item.IdentifiedDate
		}
	}
	
	// Calculate averages
	if summary.TotalItems > 0 {
		summary.AverageAge = totalAge / float64(summary.TotalItems)
		summary.CompletionRate = float64(completed) / float64(summary.TotalItems) * 100
	}
}

// GetOverdueItems returns all overdue POA&M items
func (poam *PlanOfActionMilestones) GetOverdueItems() []POAMItem {
	var overdue []POAMItem
	now := time.Now()
	
	for _, item := range poam.POAMItems {
		if (item.Status == "Open" || item.Status == "Ongoing") && item.PlannedCompletion.Before(now) {
			overdue = append(overdue, item)
		}
	}
	
	return overdue
}

// ToJSON exports the POA&M as JSON
func (poam *PlanOfActionMilestones) ToJSON() ([]byte, error) {
	return json.MarshalIndent(poam, "", "  ")
}

// GeneratePOAMFromFindings creates POA&M items from SAR findings
func GeneratePOAMFromFindings(findings []ControlFinding) []POAMItem {
	items := make([]POAMItem, 0)
	
	for _, finding := range findings {
		if finding.Status == "Other Than Satisfied" {
			item := POAMItem{
				FindingID:         finding.FindingID,
				ControlID:         finding.ControlID,
				Weakness:          finding.Description,
				Severity:          finding.Severity,
				RawRisk:           finding.RiskRating,
				Status:            "Open",
				IdentifiedDate:    finding.TestDate,
				PlannedCompletion: finding.TestDate.AddDate(0, 3, 0), // Default 3 months
				RemediationPlan:   finding.RemediationPlan,
				Source:            "SAR",
			}
			items = append(items, item)
		}
	}
	
	return items
} 