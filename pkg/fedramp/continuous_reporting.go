package fedramp

import (
	"encoding/json"
	"time"
)

// ContinuousReportingStatus represents the status of continuous reporting
type ContinuousReportingStatus struct {
	ServiceOfferingID   string                     `json:"service_offering_id"`
	LastUpdated         time.Time                  `json:"last_updated"`
	ReportingFrequency  string                     `json:"reporting_frequency"`
	AutomatedKSIs       []string                   `json:"automated_ksis"`
	ManualKSIs          []string                   `json:"manual_ksis"`
	ReportingEndpoint   string                     `json:"reporting_endpoint"`
	APIVersion          string                     `json:"api_version"`
	NextReportDue       time.Time                  `json:"next_report_due"`
	ValidationSchedule  map[string]ValidationSchedule `json:"validation_schedule"`
}

// ValidationSchedule defines when and how a KSI is validated
type ValidationSchedule struct {
	KSIID            string    `json:"ksi_id"`
	Frequency        string    `json:"frequency"` // real-time, daily, weekly, monthly
	LastValidated    time.Time `json:"last_validated"`
	NextValidation   time.Time `json:"next_validation"`
	ValidationMethod string    `json:"validation_method"` // automated, manual, hybrid
	DataSource       string    `json:"data_source"`
}

// ContinuousReportingProposal represents the proposal for continuous reporting
type ContinuousReportingProposal struct {
	ServiceOfferingID   string                          `json:"service_offering_id"`
	ProposalDate        time.Time                       `json:"proposal_date"`
	Architecture        ContinuousReportingArchitecture `json:"architecture"`
	CoveragePercentage  float64                         `json:"coverage_percentage"`
	AutomatedKSIs       []AutomatedKSIProposal          `json:"automated_ksis"`
	ReportingMechanisms []ReportingMechanism            `json:"reporting_mechanisms"`
	Implementation      ImplementationPlan              `json:"implementation_plan"`
}

// ContinuousReportingArchitecture describes the technical architecture
type ContinuousReportingArchitecture struct {
	DataCollection   []DataCollectionMethod `json:"data_collection"`
	Processing       ProcessingCapability   `json:"processing"`
	Storage          StorageCapability      `json:"storage"`
	Reporting        ReportingCapability    `json:"reporting"`
	Security         SecurityMeasures       `json:"security"`
}

// DataCollectionMethod describes how data is collected
type DataCollectionMethod struct {
	Method      string   `json:"method"`
	DataSources []string `json:"data_sources"`
	Frequency   string   `json:"frequency"`
	Automated   bool     `json:"automated"`
}

// ProcessingCapability describes data processing capabilities
type ProcessingCapability struct {
	RealTime     bool     `json:"real_time"`
	BatchSize    int      `json:"batch_size"`
	Technologies []string `json:"technologies"`
}

// StorageCapability describes data storage
type StorageCapability struct {
	RetentionPeriod string `json:"retention_period"`
	Encryption      bool   `json:"encryption"`
	Backup          bool   `json:"backup"`
	Location        string `json:"location"`
}

// ReportingCapability describes reporting capabilities
type ReportingCapability struct {
	Formats      []string `json:"formats"` // JSON, XML, HTML, PDF
	APIs         []string `json:"apis"`
	Dashboard    bool     `json:"dashboard"`
	Alerts       bool     `json:"alerts"`
	CustomReports bool     `json:"custom_reports"`
}

// SecurityMeasures describes security controls for continuous reporting
type SecurityMeasures struct {
	Authentication    string `json:"authentication"`
	Authorization     string `json:"authorization"`
	Encryption        string `json:"encryption"`
	IntegrityChecking bool   `json:"integrity_checking"`
	AuditLogging      bool   `json:"audit_logging"`
}

// AutomatedKSIProposal describes how a specific KSI will be automated
type AutomatedKSIProposal struct {
	KSIID              string   `json:"ksi_id"`
	ValidationMethod   string   `json:"validation_method"`
	DataSources        []string `json:"data_sources"`
	ValidationLogic    string   `json:"validation_logic"`
	UpdateFrequency    string   `json:"update_frequency"`
	ConfidenceLevel    string   `json:"confidence_level"`
	FallbackProcedure  string   `json:"fallback_procedure"`
}

// ReportingMechanism describes how reports will be delivered
type ReportingMechanism struct {
	Type        string `json:"type"` // API, webhook, email, portal
	Format      string `json:"format"`
	Frequency   string `json:"frequency"`
	Recipients  []string `json:"recipients,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
	Authentication string `json:"authentication,omitempty"`
}

// ImplementationPlan describes the implementation timeline
type ImplementationPlan struct {
	CurrentState    string      `json:"current_state"`
	TargetState     string      `json:"target_state"`
	Milestones      []Milestone `json:"milestones"`
	EstimatedEffort string      `json:"estimated_effort"`
	Resources       []string    `json:"resources"`
}

// Milestone represents an implementation milestone
type Milestone struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TargetDate  time.Time `json:"target_date"`
	Deliverables []string `json:"deliverables"`
}

// ContinuousReportingManager manages continuous reporting
type ContinuousReportingManager struct {
	ServiceID  string
	KSIReports map[string]*KSIReport
	Schedule   map[string]ValidationSchedule
}

// NewContinuousReportingManager creates a new manager
func NewContinuousReportingManager(serviceID string) *ContinuousReportingManager {
	return &ContinuousReportingManager{
		ServiceID:  serviceID,
		KSIReports: make(map[string]*KSIReport),
		Schedule:   make(map[string]ValidationSchedule),
	}
}

// GenerateProposal creates a continuous reporting proposal
func (m *ContinuousReportingManager) GenerateProposal() *ContinuousReportingProposal {
	proposal := &ContinuousReportingProposal{
		ServiceOfferingID: m.ServiceID,
		ProposalDate:      time.Now(),
		Architecture: ContinuousReportingArchitecture{
			DataCollection: []DataCollectionMethod{
				{
					Method:      "API Integration",
					DataSources: []string{"Cloud Provider APIs", "SIEM", "Vulnerability Scanner"},
					Frequency:   "real-time",
					Automated:   true,
				},
				{
					Method:      "Log Analysis",
					DataSources: []string{"Application Logs", "System Logs", "Security Logs"},
					Frequency:   "continuous",
					Automated:   true,
				},
			},
			Processing: ProcessingCapability{
				RealTime:     true,
				BatchSize:    1000,
				Technologies: []string{"Stream Processing", "Machine Learning", "Rule Engine"},
			},
			Storage: StorageCapability{
				RetentionPeriod: "1 year",
				Encryption:      true,
				Backup:          true,
				Location:        "FedRAMP Authorized Cloud",
			},
			Reporting: ReportingCapability{
				Formats:       []string{"JSON", "XML", "HTML", "PDF"},
				APIs:          []string{"REST API v1.0", "GraphQL"},
				Dashboard:     true,
				Alerts:        true,
				CustomReports: true,
			},
			Security: SecurityMeasures{
				Authentication:    "OAuth 2.0 with PKCE",
				Authorization:     "RBAC with least privilege",
				Encryption:        "TLS 1.3 in transit, AES-256 at rest",
				IntegrityChecking: true,
				AuditLogging:      true,
			},
		},
		ReportingMechanisms: []ReportingMechanism{
			{
				Type:           "API",
				Format:         "JSON",
				Frequency:      "real-time",
				Endpoint:       "https://api.example.com/fedramp/continuous-reporting",
				Authentication: "Bearer Token",
			},
			{
				Type:      "webhook",
				Format:    "JSON",
				Frequency: "event-driven",
				Endpoint:  "configurable",
			},
		},
	}
	
	// Calculate automated KSIs
	automatedKSIs := []AutomatedKSIProposal{
		{
			KSIID:            "KSI-CNA",
			ValidationMethod: "Cloud API queries",
			DataSources:      []string{"AWS Config", "Azure Policy", "GCP Security Command Center"},
			ValidationLogic:  "Query cloud configuration state and compare against baseline",
			UpdateFrequency:  "hourly",
			ConfidenceLevel:  "high",
			FallbackProcedure: "Manual review if API unavailable",
		},
		{
			KSIID:            "KSI-SC",
			ValidationMethod: "Configuration scanning",
			DataSources:      []string{"Infrastructure as Code", "Cloud APIs", "Certificate Manager"},
			ValidationLogic:  "Scan configurations for encryption, key rotation, patch status",
			UpdateFrequency:  "daily",
			ConfidenceLevel:  "high",
			FallbackProcedure: "Export configuration for manual review",
		},
		{
			KSIID:            "KSI-MLA",
			ValidationMethod: "SIEM integration",
			DataSources:      []string{"SIEM API", "Log aggregator", "Vulnerability scanner"},
			ValidationLogic:  "Query SIEM for log coverage, scan results, and alerts",
			UpdateFrequency:  "real-time",
			ConfidenceLevel:  "high",
			FallbackProcedure: "SIEM report generation",
		},
	}
	
	proposal.AutomatedKSIs = automatedKSIs
	proposal.CoveragePercentage = float64(len(automatedKSIs)) / float64(len(KSIDefinitions)) * 100
	
	// Implementation plan
	proposal.Implementation = ImplementationPlan{
		CurrentState: "Manual reporting with quarterly assessments",
		TargetState:  "Automated continuous reporting for 60%+ of KSIs",
		Milestones: []Milestone{
			{
				Name:        "API Development",
				Description: "Develop and test continuous reporting APIs",
				TargetDate:  time.Now().AddDate(0, 1, 0),
				Deliverables: []string{"API specification", "API implementation", "API documentation"},
			},
			{
				Name:        "Integration Testing",
				Description: "Test integration with FedRAMP systems",
				TargetDate:  time.Now().AddDate(0, 2, 0),
				Deliverables: []string{"Test results", "Integration guide"},
			},
			{
				Name:        "Production Deployment",
				Description: "Deploy continuous reporting to production",
				TargetDate:  time.Now().AddDate(0, 3, 0),
				Deliverables: []string{"Production endpoint", "Monitoring dashboard"},
			},
		},
		EstimatedEffort: "3 months, 2 FTE",
		Resources:       []string{"Development team", "Security team", "3PAO coordination"},
	}
	
	return proposal
}

// GenerateContinuousReport generates a current continuous monitoring report
func (m *ContinuousReportingManager) GenerateContinuousReport() ([]byte, error) {
	report := map[string]interface{}{
		"service_offering_id": m.ServiceID,
		"report_timestamp":    time.Now(),
		"report_type":         "continuous_monitoring",
		"ksi_validations":     make(map[string]interface{}),
		"metrics": map[string]interface{}{
			"automated_percentage": 0.0,
			"validation_frequency": "varies",
			"last_full_validation": time.Now().AddDate(0, 0, -7),
		},
	}
	
	// Add current KSI validation status
	if latestReport, exists := m.KSIReports["latest"]; exists {
		for ksiID, validation := range latestReport.Validations {
			report["ksi_validations"].(map[string]interface{})[ksiID] = map[string]interface{}{
				"status":         validation.Status,
				"last_validated": validation.LastValidated,
				"automated":      validation.AutomatedCheck,
				"evidence_count": len(validation.Evidence),
			}
		}
	}
	
	return json.MarshalIndent(report, "", "  ")
}

// ScheduleValidation schedules a KSI for validation
func (m *ContinuousReportingManager) ScheduleValidation(ksiID string, frequency string, method string) {
	nextValidation := calculateNextValidation(frequency)
	
	m.Schedule[ksiID] = ValidationSchedule{
		KSIID:            ksiID,
		Frequency:        frequency,
		LastValidated:    time.Now(),
		NextValidation:   nextValidation,
		ValidationMethod: method,
		DataSource:       "automated",
	}
}

// calculateNextValidation calculates the next validation time based on frequency
func calculateNextValidation(frequency string) time.Time {
	now := time.Now()
	switch frequency {
	case "real-time":
		return now.Add(1 * time.Minute)
	case "hourly":
		return now.Add(1 * time.Hour)
	case "daily":
		return now.Add(24 * time.Hour)
	case "weekly":
		return now.AddDate(0, 0, 7)
	case "monthly":
		return now.AddDate(0, 1, 0)
	default:
		return now.AddDate(0, 0, 1) // Default to daily
	}
} 