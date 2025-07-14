package fedramp

import (
	"encoding/json"
	"fmt"
	"time"
)

// KSIValidationStatus represents the validation status of a KSI
type KSIValidationStatus string

const (
	KSIStatusTrue    KSIValidationStatus = "True"
	KSIStatusFalse   KSIValidationStatus = "False"
	KSIStatusPartial KSIValidationStatus = "Partial"
)

// KSIEvidence represents supporting evidence for a KSI validation
type KSIEvidence struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Reference   string    `json:"reference"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
}

// KSIValidation represents a single KSI validation result
type KSIValidation struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Category         string              `json:"category"`
	Status           KSIValidationStatus `json:"status"`
	Evidence         []KSIEvidence       `json:"evidence"`
	AutomatedCheck   bool                `json:"automated_check"`
	LastValidated    time.Time           `json:"last_validated"`
	ValidationMethod string              `json:"validation_method"`
	RelatedControls  []string            `json:"related_controls"`
	ThreePAOAttested bool                `json:"3pao_attested"`
	Notes            string              `json:"notes,omitempty"`
}

// KSIReport represents a complete KSI validation report
type KSIReport struct {
	ServiceOfferingID string                    `json:"service_offering_id"`
	ReportID          string                    `json:"report_id"`
	GeneratedAt       time.Time                 `json:"generated_at"`
	Validations       map[string]*KSIValidation `json:"validations"`
	Summary           KSISummary                `json:"summary"`
	DataSchema        string                    `json:"data_schema"`
	Version           string                    `json:"version"`
}

// KSISummary provides summary statistics for KSI validations
type KSISummary struct {
	TotalKSIs          int                            `json:"total_ksis"`
	ValidationsByStatus map[KSIValidationStatus]int    `json:"validations_by_status"`
	AutomatedCount     int                            `json:"automated_count"`
	ManualCount        int                            `json:"manual_count"`
	ComplianceScore    float64                        `json:"compliance_score"`
}

// KSIDefinitions for FedRAMP 20x Phase One (Release 25.05C)
var KSIDefinitions = map[string]KSIDefinition{
	"KSI-CED": {
		ID:       "KSI-CED",
		Name:     "Cybersecurity Education",
		Category: "Training",
		Description: "A secure cloud service provider will continuously educate their employees on cybersecurity measures, testing them regularly to ensure their knowledge is satisfactory.",
		ValidationPoints: []string{
			"KSI-CED-01: Ensure all employees receive security awareness training",
			"KSI-CED-02: Require role-specific training for high risk roles, including at least roles with privileged access",
		},
		RelatedControls: []string{"AT-2", "AT-3", "AT-6"},
	},
	"KSI-CMT": {
		ID:       "KSI-CMT",
		Name:     "Change Management",
		Category: "Change Control",
		Description: "A secure cloud service provider will ensure that all system changes are properly documented and configuration baselines are updated accordingly.",
		ValidationPoints: []string{
			"KSI-CMT-01: Log and monitor system modifications",
			"KSI-CMT-02: Execute changes through redeployment of version controlled immutable resources rather than direct modification wherever possible",
			"KSI-CMT-03: Implement automated testing and validation of changes prior to deployment",
			"KSI-CMT-04: Have a documented change management procedure",
			"KSI-CMT-05: Evaluate the risk and potential impact of any change",
		},
		RelatedControls: []string{"CM-3", "CM-4", "CM-5", "CM-6", "CM-7", "CM-10", "CM-11"},
	},
	"KSI-CNA": {
		ID:       "KSI-CNA",
		Name:     "Cloud Native Architecture",
		Category: "Architecture",
		Description: "A secure cloud service offering will use cloud native architecture and design principles to enforce and enhance the Confidentiality, Integrity and Availability of the system.",
		ValidationPoints: []string{
			"KSI-CNA-01: Configure ALL information resources to limit inbound and outbound traffic",
			"KSI-CNA-02: Design systems to minimize the attack surface and minimize lateral movement if compromised",
			"KSI-CNA-03: Use logical networking and related capabilities to enforce traffic flow controls",
			"KSI-CNA-04: Use immutable infrastructure with strictly defined functionality and privileges by default",
			"KSI-CNA-05: Have denial of service protection",
			"KSI-CNA-06: Design systems for high availability and rapid recovery",
			"KSI-CNA-07: Ensure cloud-native information resources are implemented based on host provider's best practices and documented guidance",
		},
		RelatedControls: []string{"AC-4", "CA-9", "CP-2", "CP-10", "SC-5", "SC-7", "SC-7(5)", "SC-7(8)", "SC-32", "SC-36", "SC-39", "SI-4"},
	},
	"KSI-IAM": {
		ID:       "KSI-IAM",
		Name:     "Identity and Access Management",
		Category: "Access Control",
		Description: "A secure cloud service offering will protect user data, control access, and apply zero trust principles.",
		ValidationPoints: []string{
			"KSI-IAM-01: Enforce multi-factor authentication (MFA) using methods that are difficult to intercept or impersonate (phishing-resistant MFA) for all user authentication",
			"KSI-IAM-02: Use secure passwordless methods for user authentication and authorization when feasible, otherwise enforce strong passwords with MFA",
			"KSI-IAM-03: Enforce appropriately secure authentication methods for non-user accounts and services",
			"KSI-IAM-04: Use a least-privileged, role and attribute-based, and just-in-time security authorization model for all user and non-user accounts and services",
			"KSI-IAM-05: Apply zero trust design principles",
			"KSI-IAM-06: Automatically disable or otherwise secure accounts with privileged access in response to suspicious activity",
		},
		RelatedControls: []string{"AC-2", "AC-3", "AC-6", "AC-7", "AC-14", "AU-9", "IA-2", "IA-2(1)", "IA-2(2)", "IA-2(8)", "IA-2(12)", "IA-4", "IA-5", "IA-5(1)", "IA-6", "IA-8", "IA-8(1)", "IA-8(2)", "IA-8(4)", "IA-11", "PS-2", "PS-3", "PS-4", "PS-5", "PS-7", "PS-9"},
	},
	"KSI-INR": {
		ID:       "KSI-INR",
		Name:     "Incident Reporting",
		Category: "Incident Management",
		Description: "A secure cloud service offering will document, report, and analyze security incidents to ensure regulatory compliance and continuous security improvement.",
		ValidationPoints: []string{
			"KSI-INR-01: Report incidents according to FedRAMP requirements and cloud service provider policies",
			"KSI-INR-02: Maintain a log of incidents and periodically review past incidents for patterns or vulnerabilities",
			"KSI-INR-03: Generate after action reports and regularly incorporate lessons learned into operations",
		},
		RelatedControls: []string{"IR-4", "IR-5", "IR-6", "IR-7", "IR-8"},
	},
	"KSI-MLA": {
		ID:       "KSI-MLA",
		Name:     "Monitoring, Logging, and Auditing",
		Category: "Monitoring",
		Description: "A secure cloud service offering will monitor, log, and audit all important events, activity, and changes.",
		ValidationPoints: []string{
			"KSI-MLA-01: Operate a Security Information and Event Management (SIEM) or similar system(s) for centralized, tamper-resistent logging of events, activities, and changes",
			"KSI-MLA-02: Regularly review and audit logs",
			"KSI-MLA-03: Rapidly detect and remediate or mitigate vulnerabilities",
			"KSI-MLA-04: Perform authenticated vulnerability scanning on information resources",
			"KSI-MLA-05: Perform Infrastructure as Code and configuration evaluation and testing",
			"KSI-MLA-06: Centrally track and prioritize the mitigation and/or remediation of identified vulnerabilities",
		},
		RelatedControls: []string{"AC-7", "AU-2", "AU-3", "AU-4", "AU-6", "AU-8", "AU-9", "AU-11", "AU-12", "CA-7", "CM-6", "RA-5", "RA-5(2)", "RA-5(11)", "SI-2", "SI-4", "SI-5"},
	},
	"KSI-PIY": {
		ID:       "KSI-PIY",
		Name:     "Policy and Inventory",
		Category: "Governance",
		Description: "A secure cloud service offering will have intentional, organized, universal guidance for how every information resource, including personnel, is secured.",
		ValidationPoints: []string{
			"KSI-PIY-01: Have an up-to-date information resource inventory or code defining all deployed assets, software, and services",
			"KSI-PIY-02: Have policies outlining the security objectives of all information resources",
			"KSI-PIY-03: Maintain a vulnerability disclosure program",
			"KSI-PIY-04: Build security considerations into the Software Development Lifecycle and align with CISA Secure By Design principles",
			"KSI-PIY-05: Document methods used to evaluate information resource implementations",
			"KSI-PIY-06: Have a dedicated staff and budget for security with executive support, commensurate with the size, complexity, scope, and risk of the service offering",
			"KSI-PIY-07: Document risk management decisions for software supply chain security",
		},
		RelatedControls: []string{"AC-1", "AU-1", "CA-1", "CM-1", "CM-8", "CP-1", "IA-1", "IR-1", "MA-1", "MP-1", "PE-1", "PL-1", "PL-2", "PM-1", "PM-3", "PM-11", "PS-1", "RA-1", "SA-1", "SA-2", "SA-3", "SA-5", "SA-8", "SC-1", "SI-1", "SR-1"},
	},
	"KSI-RPL": {
		ID:       "KSI-RPL",
		Name:     "Recovery Planning",
		Category: "Continuity",
		Description: "A secure cloud service offering will define, maintain, and test incident response plan(s) and recovery capabilities to ensure minimal service disruption and data loss during incidents and contingencies.",
		ValidationPoints: []string{
			"KSI-RPL-01: Define Recovery Time Objectives (RTO) and Recovery Point Objectives (RPO)",
			"KSI-RPL-02: Develop and maintain a recovery plan that aligns with the defined recovery objectives",
			"KSI-RPL-03: Perform system backups aligned with recovery objectives",
			"KSI-RPL-04: Regularly test the capability to recover from incidents and contingencies",
		},
		RelatedControls: []string{"CP-2", "CP-4", "CP-9", "CP-10", "IR-4"},
	},
	"KSI-SVC": {
		ID:       "KSI-SVC",
		Name:     "Service Configuration",
		Category: "Configuration",
		Description: "A secure cloud service offering will follow FedRAMP encryption policies, continuously verify information resource integrity, and restrict access to third-party information resources.",
		ValidationPoints: []string{
			"KSI-SVC-01: Harden and review network and system configurations",
			"KSI-SVC-02: Encrypt or otherwise secure network traffic",
			"KSI-SVC-03: Encrypt all federal and sensitive information at rest",
			"KSI-SVC-04: Manage configuration centrally",
			"KSI-SVC-05: Enforce system and information resource integrity through cryptographic means",
			"KSI-SVC-06: Use automated key management systems to manage, protect, and regularly rotate digital keys and certificates",
			"KSI-SVC-07: Use a consistent, risk-informed approach for applying security patches",
		},
		RelatedControls: []string{"CM-2", "CM-6", "IA-5(7)", "IA-7", "SC-8", "SC-8(1)", "SC-13", "SC-28", "SC-28(1)", "SI-2", "SI-3", "SI-6", "SI-7", "SI-7(1)", "SI-7(6)"},
	},
	"KSI-TPR": {
		ID:       "KSI-TPR",
		Name:     "Third-Party Information Resources",
		Category: "Supply Chain",
		Description: "A secure cloud service offering will understand, monitor, and manage supply chain risks from third-party information resources.",
		ValidationPoints: []string{
			"KSI-TPR-01: Identify all third-party information resources",
			"KSI-TPR-02: Regularly confirm that services handling federal information or are likely to impact the confidentiality, integrity, or availability of federal information are FedRAMP authorized and securely configured",
			"KSI-TPR-03: Identify and prioritize mitigation of potential supply chain risks",
			"KSI-TPR-04: Monitor third party software information resources for upstream vulnerabilities, with contractual notification requirements or active monitoring services",
		},
		RelatedControls: []string{"AC-20", "CA-3", "RA-3(1)", "SA-4", "SA-9", "SA-12", "SI-5", "SR-2", "SR-2(1)", "SR-3", "SR-4", "SR-5", "SR-6", "SR-8", "SR-10", "SR-11", "SR-11(2)", "SR-11(3)"},
	},
}

// KSIDefinition defines the structure of a KSI
type KSIDefinition struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	ValidationPoints []string `json:"validation_points"`
	RelatedControls  []string `json:"related_controls"`
}

// KSIRequirement represents an individual KSI requirement (e.g., KSI-CNA-01)
type KSIRequirement struct {
	ID          string    `json:"id"`
	ParentKSI   string    `json:"parent_ksi"`
	Description string    `json:"description"`
	Validated   bool      `json:"validated"`
	Evidence    []KSIEvidence `json:"evidence"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
}

// NewKSIReport creates a new KSI validation report
func NewKSIReport(serviceID string) *KSIReport {
	return &KSIReport{
		ServiceOfferingID: serviceID,
		ReportID:          fmt.Sprintf("KSI-%s-%s", serviceID, time.Now().Format("20060102-150405")),
		GeneratedAt:       time.Now(),
		Validations:       make(map[string]*KSIValidation),
		Version:           "1.0.0",
		DataSchema:        "https://github.com/FedRAMP/docs/blob/main/FRMR.KSI.key-security-indicators.json",
	}
}

// AddValidation adds a KSI validation to the report
func (r *KSIReport) AddValidation(validation *KSIValidation) {
	r.Validations[validation.ID] = validation
	r.updateSummary()
}

// updateSummary recalculates the summary statistics
func (r *KSIReport) updateSummary() {
	r.Summary.TotalKSIs = len(r.Validations)
	r.Summary.ValidationsByStatus = make(map[KSIValidationStatus]int)
	r.Summary.AutomatedCount = 0
	r.Summary.ManualCount = 0
	
	trueCount := 0
	
	for _, v := range r.Validations {
		r.Summary.ValidationsByStatus[v.Status]++
		if v.AutomatedCheck {
			r.Summary.AutomatedCount++
		} else {
			r.Summary.ManualCount++
		}
		if v.Status == KSIStatusTrue {
			trueCount++
		}
	}
	
	if r.Summary.TotalKSIs > 0 {
		r.Summary.ComplianceScore = float64(trueCount) / float64(r.Summary.TotalKSIs) * 100
	}
}

// ToJSON exports the KSI report as JSON
func (r *KSIReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// ValidateKSI performs validation for a specific KSI
func ValidateKSI(ksiID string, evidence []KSIEvidence, automated bool) *KSIValidation {
	def, exists := KSIDefinitions[ksiID]
	if !exists {
		return nil
	}
	
	validation := &KSIValidation{
		ID:               ksiID,
		Name:             def.Name,
		Category:         def.Category,
		Evidence:         evidence,
		AutomatedCheck:   automated,
		LastValidated:    time.Now(),
		RelatedControls:  def.RelatedControls,
		ThreePAOAttested: false,
	}
	
	// Count how many validation points have evidence
	coveredPoints := 0
	evidenceMap := make(map[string]bool)
	
	// Map evidence types to validation points
	for _, e := range evidence {
		evidenceMap[e.Type] = true
	}
	
	// Check coverage of validation points
	for _, point := range def.ValidationPoints {
		// Extract the KSI ID from the validation point (e.g., "KSI-CNA-01")
		pointID := point[:10]
		if evidenceMap[pointID] {
			coveredPoints++
		}
	}
	
	// Determine status based on coverage
	totalPoints := len(def.ValidationPoints)
	if coveredPoints == totalPoints {
		validation.Status = KSIStatusTrue
	} else if coveredPoints >= (totalPoints+1)/2 { // More than half
		validation.Status = KSIStatusPartial
	} else {
		validation.Status = KSIStatusFalse
	}
	
	if automated {
		validation.ValidationMethod = "automated"
	} else {
		validation.ValidationMethod = "manual"
	}
	
	return validation
} 

// NewKSIValidation creates a new basic KSI validation structure for a CSO
func NewKSIValidation(csoID string) *KSIValidation {
	// Return a simple validation structure that matches the existing KSIValidation type
	validation := &KSIValidation{
		ID:               fmt.Sprintf("KSI-VAL-%d", time.Now().Unix()),
		Name:             "KSI Compliance Validation",
		Category:         "Compliance",
		Status:           KSIStatusFalse,
		Evidence:         []KSIEvidence{},
		AutomatedCheck:   false,
		LastValidated:    time.Now(),
		ValidationMethod: "manual",
		RelatedControls:  []string{},
		ThreePAOAttested: false,
		Notes:            fmt.Sprintf("Validation for CSO: %s", csoID),
	}
	
	return validation
}

// GenerateKSIReport generates a complete KSI report for a CSO
func GenerateKSIReport(csoID string, reportDate time.Time) *KSIReport {
	report := NewKSIReport(csoID)
	
	// Add validations for each KSI
	for ksiID, def := range KSIDefinitions {
		// Create mock evidence for demonstration
		evidence := []KSIEvidence{
			{
				Type:        def.ValidationPoints[0][:10], // Use first validation point ID
				Description: "Evidence for " + def.Name,
				Reference:   "DOC-" + ksiID,
				Timestamp:   reportDate,
				Source:      "Automated scan",
			},
		}
		
		validation := ValidateKSI(ksiID, evidence, true)
		if validation != nil {
			report.AddValidation(validation)
		}
	}
	
	return report
}