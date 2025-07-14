package fedramp

import (
	"encoding/json"
	"fmt"
	"time"
)

// KeySecurityMetric represents a single security metric for continuous reporting
type KeySecurityMetric struct {
	MetricID          string                 `json:"metric_id"`
	MetricName        string                 `json:"metric_name"`
	Category          string                 `json:"category"`
	Description       string                 `json:"description"`
	Value             interface{}            `json:"value"`
	Unit              string                 `json:"unit"`
	Threshold         interface{}            `json:"threshold,omitempty"`
	Status            MetricStatus           `json:"status"`
	LastUpdated       time.Time              `json:"last_updated"`
	CollectionMethod  string                 `json:"collection_method"`
	Frequency         string                 `json:"frequency"`
	RelatedControls   []string               `json:"related_controls"`
	Evidence          []MetricEvidence       `json:"evidence,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// MetricStatus represents the status of a security metric
type MetricStatus string

const (
	MetricStatusGreen  MetricStatus = "green"  // Within acceptable thresholds
	MetricStatusYellow MetricStatus = "yellow" // Warning - approaching threshold
	MetricStatusRed    MetricStatus = "red"    // Critical - threshold exceeded
	MetricStatusGray   MetricStatus = "gray"   // Unknown or not applicable
)

// MetricEvidence represents supporting evidence for a metric
type MetricEvidence struct {
	EvidenceType string    `json:"evidence_type"`
	Description  string    `json:"description"`
	Source       string    `json:"source"`
	Timestamp    time.Time `json:"timestamp"`
	Reference    string    `json:"reference,omitempty"`
}

// ContinuousReport represents a continuous monitoring report
type ContinuousReport struct {
	ReportID          string              `json:"report_id"`
	ServiceOfferingID string              `json:"service_offering_id"`
	ReportingPeriod   ReportingPeriod     `json:"reporting_period"`
	GeneratedAt       time.Time           `json:"generated_at"`
	Metrics           []KeySecurityMetric `json:"metrics"`
	Summary           ReportSummary       `json:"summary"`
	Incidents         []SecurityIncident  `json:"incidents,omitempty"`
	Changes           []string            `json:"changes,omitempty"` // References to SCNs
	Attestation       ReportAttestation   `json:"attestation"`
}

// ReportingPeriod defines the time period for the report
type ReportingPeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Type      string    `json:"type"` // daily, weekly, monthly, quarterly
}

// ReportSummary provides high-level summary information
type ReportSummary struct {
	TotalMetrics     int                        `json:"total_metrics"`
	MetricsByStatus  map[MetricStatus]int       `json:"metrics_by_status"`
	CriticalFindings int                        `json:"critical_findings"`
	TrendAnalysis    map[string]interface{}     `json:"trend_analysis,omitempty"`
	KeyInsights      []string                   `json:"key_insights,omitempty"`
}

// SecurityIncident represents a security incident for reporting
type SecurityIncident struct {
	IncidentID      string    `json:"incident_id"`
	Title           string    `json:"title"`
	Severity        string    `json:"severity"`
	Status          string    `json:"status"`
	DetectedAt      time.Time `json:"detected_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	AffectedSystems []string  `json:"affected_systems"`
	Response        string    `json:"response"`
}

// ReportAttestation represents the attestation for the report
type ReportAttestation struct {
	AttestorName     string    `json:"attestor_name"`
	AttestorTitle    string    `json:"attestor_title"`
	AttestorEmail    string    `json:"attestor_email"`
	AttestationDate  time.Time `json:"attestation_date"`
	Statement        string    `json:"statement"`
	DigitalSignature string    `json:"digital_signature,omitempty"`
}

// CRSManager handles Continuous Reporting Standard operations
type CRSManager struct {
	reports map[string]*ContinuousReport
	metrics map[string]*KeySecurityMetric
}

// NewCRSManager creates a new CRS manager
func NewCRSManager() *CRSManager {
	return &CRSManager{
		reports: make(map[string]*ContinuousReport),
		metrics: make(map[string]*KeySecurityMetric),
	}
}

// CreateReport creates a new continuous monitoring report
func (mgr *CRSManager) CreateReport(serviceID string, period ReportingPeriod) *ContinuousReport {
	reportID := fmt.Sprintf("%s-%s-%s", serviceID, period.Type, period.EndDate.Format("2006-01-02"))
	
	report := &ContinuousReport{
		ReportID:          reportID,
		ServiceOfferingID: serviceID,
		ReportingPeriod:   period,
		GeneratedAt:       time.Now(),
		Metrics:           make([]KeySecurityMetric, 0),
		Summary: ReportSummary{
			MetricsByStatus: make(map[MetricStatus]int),
		},
		Incidents: make([]SecurityIncident, 0),
		Changes:   make([]string, 0),
	}
	
	mgr.reports[reportID] = report
	return report
}

// AddMetric adds a metric to a report
func (mgr *CRSManager) AddMetric(reportID string, metric KeySecurityMetric) error {
	report, exists := mgr.reports[reportID]
	if !exists {
		return fmt.Errorf("report %s not found", reportID)
	}
	
	metric.LastUpdated = time.Now()
	report.Metrics = append(report.Metrics, metric)
	mgr.updateReportSummary(report)
	
	return nil
}

// updateReportSummary updates the summary statistics for a report
func (mgr *CRSManager) updateReportSummary(report *ContinuousReport) {
	summary := &report.Summary
	summary.TotalMetrics = len(report.Metrics)
	summary.MetricsByStatus = make(map[MetricStatus]int)
	summary.CriticalFindings = 0
	
	for _, metric := range report.Metrics {
		summary.MetricsByStatus[metric.Status]++
		if metric.Status == MetricStatusRed {
			summary.CriticalFindings++
		}
	}
}

// GenerateStandardMetrics creates standard FedRAMP continuous monitoring metrics
func (mgr *CRSManager) GenerateStandardMetrics() []KeySecurityMetric {
	now := time.Now()
	
	return []KeySecurityMetric{
		{
			MetricID:         "vuln-scan-coverage",
			MetricName:       "Vulnerability Scan Coverage",
			Category:         "Vulnerability Management",
			Description:      "Percentage of systems covered by vulnerability scanning",
			Value:            95.5,
			Unit:             "percentage",
			Threshold:        90.0,
			Status:           MetricStatusGreen,
			LastUpdated:      now,
			CollectionMethod: "automated",
			Frequency:        "daily",
			RelatedControls:  []string{"RA-5", "RA-5(1)"},
		},
		{
			MetricID:         "patch-compliance",
			MetricName:       "Security Patch Compliance",
			Category:         "Configuration Management",
			Description:      "Percentage of systems with current security patches",
			Value:            88.2,
			Unit:             "percentage",
			Threshold:        95.0,
			Status:           MetricStatusYellow,
			LastUpdated:      now,
			CollectionMethod: "automated",
			Frequency:        "daily",
			RelatedControls:  []string{"SI-2", "CM-6"},
		},
		{
			MetricID:         "failed-logins",
			MetricName:       "Failed Login Attempts",
			Category:         "Access Control",
			Description:      "Number of failed login attempts in the last 24 hours",
			Value:            23,
			Unit:             "count",
			Threshold:        100,
			Status:           MetricStatusGreen,
			LastUpdated:      now,
			CollectionMethod: "automated",
			Frequency:        "real-time",
			RelatedControls:  []string{"AC-7", "AU-2"},
		},
		{
			MetricID:         "backup-success-rate",
			MetricName:       "Backup Success Rate",
			Category:         "Contingency Planning",
			Description:      "Percentage of successful automated backups",
			Value:            99.8,
			Unit:             "percentage",
			Threshold:        99.0,
			Status:           MetricStatusGreen,
			LastUpdated:      now,
			CollectionMethod: "automated",
			Frequency:        "daily",
			RelatedControls:  []string{"CP-9", "CP-10"},
		},
		{
			MetricID:         "encryption-coverage",
			MetricName:       "Data Encryption Coverage",
			Category:         "System and Communications Protection",
			Description:      "Percentage of data encrypted at rest and in transit",
			Value:            100.0,
			Unit:             "percentage",
			Threshold:        100.0,
			Status:           MetricStatusGreen,
			LastUpdated:      now,
			CollectionMethod: "automated",
			Frequency:        "daily",
			RelatedControls:  []string{"SC-8", "SC-28"},
		},
		{
			MetricID:         "mfa-coverage",
			MetricName:       "Multi-Factor Authentication Coverage",
			Category:         "Identification and Authentication",
			Description:      "Percentage of privileged accounts using MFA",
			Value:            100.0,
			Unit:             "percentage",
			Threshold:        100.0,
			Status:           MetricStatusGreen,
			LastUpdated:      now,
			CollectionMethod: "automated",
			Frequency:        "daily",
			RelatedControls:  []string{"IA-2(1)", "IA-2(2)"},
		},
	}
}

// ExportReport exports a continuous monitoring report as JSON
func (mgr *CRSManager) ExportReport(reportID string) ([]byte, error) {
	report, exists := mgr.reports[reportID]
	if !exists {
		return nil, fmt.Errorf("report %s not found", reportID)
	}
	
	return json.MarshalIndent(report, "", "  ")
}

// ValidateReport validates a continuous monitoring report
func (mgr *CRSManager) ValidateReport(reportID string) error {
	report, exists := mgr.reports[reportID]
	if !exists {
		return fmt.Errorf("report %s not found", reportID)
	}
	
	// Validate required fields
	if report.ServiceOfferingID == "" {
		return fmt.Errorf("service offering ID is required")
	}
	
	if report.Attestation.AttestorName == "" {
		return fmt.Errorf("attestor name is required")
	}
	
	if len(report.Metrics) == 0 {
		return fmt.Errorf("at least one metric is required")
	}
	
	// Validate each metric
	for i, metric := range report.Metrics {
		if metric.MetricID == "" {
			return fmt.Errorf("metric %d: metric ID is required", i)
		}
		if metric.MetricName == "" {
			return fmt.Errorf("metric %d: metric name is required", i)
		}
		if metric.Value == nil {
			return fmt.Errorf("metric %d: metric value is required", i)
		}
	}
	
	return nil
}

// GetMetricTrends analyzes trends for a specific metric across multiple reports
func (mgr *CRSManager) GetMetricTrends(serviceID, metricID string, days int) map[string]interface{} {
	trends := map[string]interface{}{
		"metric_id": metricID,
		"service_id": serviceID,
		"period_days": days,
		"data_points": make([]map[string]interface{}, 0),
		"trend_direction": "stable",
		"average_value": 0.0,
	}
	
	// This would be implemented to analyze historical data
	// For now, return placeholder structure
	
	return trends
}

// GenerateDashboardData creates data for a continuous monitoring dashboard
func (mgr *CRSManager) GenerateDashboardData(serviceID string) map[string]interface{} {
	dashboard := map[string]interface{}{
		"service_id": serviceID,
		"last_updated": time.Now(),
		"overall_status": "green",
		"metrics_summary": map[string]int{
			"total": 0,
			"green": 0,
			"yellow": 0,
			"red": 0,
		},
		"recent_incidents": make([]SecurityIncident, 0),
		"key_metrics": make([]KeySecurityMetric, 0),
	}
	
	// Aggregate data from recent reports
	for _, report := range mgr.reports {
		if report.ServiceOfferingID == serviceID {
			// Add metrics to dashboard
			for _, metric := range report.Metrics {
				dashboard["key_metrics"] = append(dashboard["key_metrics"].([]KeySecurityMetric), metric)
			}
			
			// Add incidents
			for _, incident := range report.Incidents {
				dashboard["recent_incidents"] = append(dashboard["recent_incidents"].([]SecurityIncident), incident)
			}
		}
	}
	
	return dashboard
} 

// KeySecurityMetrics represents the key security metrics for FedRAMP
type KeySecurityMetrics struct {
	VulnerabilityScanning   VulnerabilityMetric   `json:"vulnerability_scanning"`
	SecurityIncidents       IncidentMetric        `json:"security_incidents"`
	ConfigurationManagement ConfigurationMetric   `json:"configuration_management"`
	AccessControl          AccessControlMetric    `json:"access_control"`
	SystemAvailability     AvailabilityMetric    `json:"system_availability"`
	PatchManagement        PatchMetric           `json:"patch_management"`
}

// VulnerabilityMetric tracks vulnerability scanning metrics
type VulnerabilityMetric struct {
	ScansCompleted      int       `json:"scans_completed"`
	CriticalFindings    int       `json:"critical_findings"`
	HighFindings        int       `json:"high_findings"`
	MediumFindings      int       `json:"medium_findings"`
	LowFindings         int       `json:"low_findings"`
	LastScanDate        time.Time `json:"last_scan_date"`
	RemediationRate     float64   `json:"remediation_rate"`
}

// IncidentMetric tracks security incident metrics
type IncidentMetric struct {
	TotalIncidents      int     `json:"total_incidents"`
	OpenIncidents       int     `json:"open_incidents"`
	ClosedIncidents     int     `json:"closed_incidents"`
	AverageResolutionTime float64 `json:"average_resolution_time_hours"`
	IncidentsByCategory map[string]int `json:"incidents_by_category"`
}

// ConfigurationMetric tracks configuration management metrics
type ConfigurationMetric struct {
	BaselineCompliance  float64 `json:"baseline_compliance_percentage"`
	UnauthorizedChanges int     `json:"unauthorized_changes"`
	ConfigurationDrift  int     `json:"configuration_drift_instances"`
	LastBaselineReview  time.Time `json:"last_baseline_review"`
}

// AccessControlMetric tracks access control metrics
type AccessControlMetric struct {
	ActiveUsers         int     `json:"active_users"`
	PrivilegedUsers     int     `json:"privileged_users"`
	FailedLoginAttempts int     `json:"failed_login_attempts"`
	AccountLockouts     int     `json:"account_lockouts"`
	MFAAdoptionRate     float64 `json:"mfa_adoption_rate"`
}

// AvailabilityMetric tracks system availability metrics
type AvailabilityMetric struct {
	UptimePercentage    float64 `json:"uptime_percentage"`
	PlannedDowntime     float64 `json:"planned_downtime_hours"`
	UnplannedDowntime   float64 `json:"unplanned_downtime_hours"`
	MTTRHours           float64 `json:"mttr_hours"`
}

// PatchMetric tracks patch management metrics
type PatchMetric struct {
	PatchesAvailable    int     `json:"patches_available"`
	PatchesApplied      int     `json:"patches_applied"`
	CriticalPatches     int     `json:"critical_patches"`
	PatchComplianceRate float64 `json:"patch_compliance_rate"`
	AveragePatchAge     float64 `json:"average_patch_age_days"`
} 

// ContinuousReportingStandard represents the FedRAMP continuous reporting standard
type ContinuousReportingStandard struct {
	CSOId           string             `json:"cso_id"`
	ReportingPeriod string             `json:"reporting_period"`
	GeneratedAt     time.Time          `json:"generated_at"`
	Metrics         KeySecurityMetrics `json:"metrics"`
	ComplianceScore float64            `json:"compliance_score"`
	Status          string             `json:"status"`
} 