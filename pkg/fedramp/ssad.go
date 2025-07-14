package fedramp

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// SSADPackage represents a standardized authorization data package
type SSADPackage struct {
	PackageID         string                 `json:"package_id"`
	ServiceOfferingID string                 `json:"service_offering_id"`
	Version           string                 `json:"version"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Status            string                 `json:"status"` // draft, final, archived
	Metadata          SSADMetadata           `json:"metadata"`
	Components        SSADComponents         `json:"components"`
	AccessControl     SSADAccessControl      `json:"access_control"`
	DistributionLog   []SSADDistribution     `json:"distribution_log"`
	IntegrityCheck    SSADIntegrity          `json:"integrity_check"`
}

// SSADMetadata contains package metadata
type SSADMetadata struct {
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	AuthorizationType string            `json:"authorization_type"` // JAB, Agency, FedRAMP Tailored
	ImpactLevel       string            `json:"impact_level"`
	AuthorizationDate time.Time         `json:"authorization_date"`
	ExpirationDate    time.Time         `json:"expiration_date"`
	CSPName           string            `json:"csp_name"`
	PackageFormat     string            `json:"package_format"` // OSCAL, Legacy, Hybrid
	Tags              []string          `json:"tags"`
	Keywords          []string          `json:"keywords"`
	CustomMetadata    map[string]string `json:"custom_metadata,omitempty"`
}

// SSADComponents represents the components of an authorization package
type SSADComponents struct {
	SSP              *SSADDocument   `json:"ssp"`
	SAP              *SSADDocument   `json:"sap"`
	SAR              *SSADDocument   `json:"sar"`
	POAM             *SSADDocument   `json:"poam"`
	ConMon           []SSADDocument  `json:"continuous_monitoring"`
	SCNs             []SSADDocument  `json:"scns"`
	IncidentReports  []SSADDocument  `json:"incident_reports"`
	Attachments      []SSADDocument  `json:"attachments"`
	KSIReports       []SSADDocument  `json:"ksi_reports,omitempty"` // For 20x
}

// SSADDocument represents a document in the package
type SSADDocument struct {
	DocumentID      string            `json:"document_id"`
	Title           string            `json:"title"`
	Type            string            `json:"type"`
	Format          string            `json:"format"` // JSON, XML, PDF, DOCX
	Version         string            `json:"version"`
	CreatedDate     time.Time         `json:"created_date"`
	LastModified    time.Time         `json:"last_modified"`
	Author          string            `json:"author"`
	Size            int64             `json:"size"`
	Hash            string            `json:"hash"`
	Location        string            `json:"location"` // URL or path
	AccessLevel     string            `json:"access_level"` // public, restricted, confidential
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// SSADAccessControl defines access permissions
type SSADAccessControl struct {
	Owner            string              `json:"owner"`
	SharedWith       []SSADSharee        `json:"shared_with"`
	PublicAccess     bool                `json:"public_access"`
	AccessRestrictions []string          `json:"access_restrictions"`
	DataClassification string            `json:"data_classification"`
	ExportControl    bool                `json:"export_control"`
}

// SSADSharee represents an entity with access
type SSADSharee struct {
	EntityID         string    `json:"entity_id"`
	EntityType       string    `json:"entity_type"` // agency, 3pao, public
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	AccessLevel      string    `json:"access_level"` // read, write, admin
	SharedDate       time.Time `json:"shared_date"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
	AccessConditions []string  `json:"access_conditions,omitempty"`
}

// SSADDistribution logs package distribution
type SSADDistribution struct {
	DistributionID   string    `json:"distribution_id"`
	RecipientID      string    `json:"recipient_id"`
	RecipientName    string    `json:"recipient_name"`
	DistributionDate time.Time `json:"distribution_date"`
	Method           string    `json:"method"` // api, portal, email
	Purpose          string    `json:"purpose"`
	Acknowledgment   bool      `json:"acknowledgment"`
	AckDate          *time.Time `json:"ack_date,omitempty"`
}

// SSADIntegrity ensures package integrity
type SSADIntegrity struct {
	PackageHash      string    `json:"package_hash"`
	HashAlgorithm    string    `json:"hash_algorithm"`
	SignedBy         string    `json:"signed_by"`
	SignatureDate    time.Time `json:"signature_date"`
	VerificationKey  string    `json:"verification_key"`
	BlockchainRef    string    `json:"blockchain_ref,omitempty"`
}

// SSADRepository manages authorization packages
type SSADRepository struct {
	RepositoryID string                    `json:"repository_id"`
	Name         string                    `json:"name"`
	Type         string                    `json:"type"` // central, agency, csp
	Packages     map[string]*SSADPackage   `json:"packages"`
	Index        SSADIndex                 `json:"index"`
	APIEndpoint  string                    `json:"api_endpoint"`
}

// SSADIndex provides searchable index
type SSADIndex struct {
	ByCSP          map[string][]string `json:"by_csp"`
	ByImpactLevel  map[string][]string `json:"by_impact_level"`
	ByAuthType     map[string][]string `json:"by_auth_type"`
	ByTag          map[string][]string `json:"by_tag"`
	LastUpdated    time.Time           `json:"last_updated"`
}

// NewSSADPackage creates a new SSAD package
func NewSSADPackage(serviceID string, metadata SSADMetadata) *SSADPackage {
	now := time.Now()
	return &SSADPackage{
		PackageID:         fmt.Sprintf("SSAD-%s-%s", serviceID, now.Format("20060102-150405")),
		ServiceOfferingID: serviceID,
		Version:           "1.0.0",
		CreatedAt:         now,
		UpdatedAt:         now,
		Status:           "draft",
		Metadata:         metadata,
		Components: SSADComponents{
			ConMon:          make([]SSADDocument, 0),
			SCNs:            make([]SSADDocument, 0),
			IncidentReports: make([]SSADDocument, 0),
			Attachments:     make([]SSADDocument, 0),
			KSIReports:      make([]SSADDocument, 0),
		},
		AccessControl: SSADAccessControl{
			SharedWith: make([]SSADSharee, 0),
		},
		DistributionLog: make([]SSADDistribution, 0),
	}
}

// AddDocument adds a document to the package
func (p *SSADPackage) AddDocument(docType string, doc SSADDocument) error {
	// Calculate document hash
	doc.Hash = calculateHash(doc)
	
	switch docType {
	case "ssp":
		p.Components.SSP = &doc
	case "sap":
		p.Components.SAP = &doc
	case "sar":
		p.Components.SAR = &doc
	case "poam":
		p.Components.POAM = &doc
	case "conmon":
		p.Components.ConMon = append(p.Components.ConMon, doc)
	case "scn":
		p.Components.SCNs = append(p.Components.SCNs, doc)
	case "incident":
		p.Components.IncidentReports = append(p.Components.IncidentReports, doc)
	case "ksi":
		p.Components.KSIReports = append(p.Components.KSIReports, doc)
	case "attachment":
		p.Components.Attachments = append(p.Components.Attachments, doc)
	default:
		return fmt.Errorf("unknown document type: %s", docType)
	}
	
	p.UpdatedAt = time.Now()
	return nil
}

// ShareWith adds sharing permissions
func (p *SSADPackage) ShareWith(sharee SSADSharee) {
	sharee.SharedDate = time.Now()
	p.AccessControl.SharedWith = append(p.AccessControl.SharedWith, sharee)
	p.UpdatedAt = time.Now()
}

// LogDistribution logs a distribution event
func (p *SSADPackage) LogDistribution(recipient, purpose, method string) {
	dist := SSADDistribution{
		DistributionID:   fmt.Sprintf("DIST-%s-%s", p.PackageID, time.Now().Format("20060102-150405")),
		RecipientID:      recipient,
		RecipientName:    recipient, // In practice, would look up name
		DistributionDate: time.Now(),
		Method:          method,
		Purpose:         purpose,
		Acknowledgment:  false,
	}
	
	p.DistributionLog = append(p.DistributionLog, dist)
}

// Finalize prepares the package for distribution
func (p *SSADPackage) Finalize(signedBy string) error {
	if p.Status != "draft" {
		return fmt.Errorf("can only finalize draft packages")
	}
	
	// Validate completeness
	if p.Components.SSP == nil || p.Components.SAR == nil {
		return fmt.Errorf("package must contain at least SSP and SAR")
	}
	
	// Calculate package integrity
	p.IntegrityCheck = SSADIntegrity{
		PackageHash:     p.calculatePackageHash(),
		HashAlgorithm:   "SHA-256",
		SignedBy:        signedBy,
		SignatureDate:   time.Now(),
		VerificationKey: fmt.Sprintf("verify-%s", p.PackageID),
	}
	
	p.Status = "final"
	p.UpdatedAt = time.Now()
	
	return nil
}

// calculatePackageHash calculates hash of entire package
func (p *SSADPackage) calculatePackageHash() string {
	// In practice, this would hash all components
	data, _ := json.Marshal(p.Components)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// calculateHash calculates hash for a document
func calculateHash(doc SSADDocument) string {
	data := fmt.Sprintf("%s-%s-%s-%d", doc.DocumentID, doc.Title, doc.Version, doc.Size)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// NewSSADRepository creates a new repository
func NewSSADRepository(name, repoType string) *SSADRepository {
	return &SSADRepository{
		RepositoryID: fmt.Sprintf("REPO-%s-%s", name, time.Now().Format("20060102")),
		Name:        name,
		Type:        repoType,
		Packages:    make(map[string]*SSADPackage),
		Index: SSADIndex{
			ByCSP:         make(map[string][]string),
			ByImpactLevel: make(map[string][]string),
			ByAuthType:    make(map[string][]string),
			ByTag:         make(map[string][]string),
		},
		APIEndpoint: fmt.Sprintf("https://api.fedramp.gov/ssad/%s", name),
	}
}

// AddPackage adds a package to the repository
func (r *SSADRepository) AddPackage(pkg *SSADPackage) error {
	if pkg.Status != "final" {
		return fmt.Errorf("only finalized packages can be added to repository")
	}
	
	r.Packages[pkg.PackageID] = pkg
	r.updateIndex(pkg)
	
	return nil
}

// updateIndex updates the searchable index
func (r *SSADRepository) updateIndex(pkg *SSADPackage) {
	// Index by CSP
	r.Index.ByCSP[pkg.Metadata.CSPName] = append(r.Index.ByCSP[pkg.Metadata.CSPName], pkg.PackageID)
	
	// Index by impact level
	r.Index.ByImpactLevel[pkg.Metadata.ImpactLevel] = append(r.Index.ByImpactLevel[pkg.Metadata.ImpactLevel], pkg.PackageID)
	
	// Index by authorization type
	r.Index.ByAuthType[pkg.Metadata.AuthorizationType] = append(r.Index.ByAuthType[pkg.Metadata.AuthorizationType], pkg.PackageID)
	
	// Index by tags
	for _, tag := range pkg.Metadata.Tags {
		r.Index.ByTag[tag] = append(r.Index.ByTag[tag], pkg.PackageID)
	}
	
	r.Index.LastUpdated = time.Now()
}

// Search searches packages in the repository
func (r *SSADRepository) Search(criteria map[string]string) []*SSADPackage {
	results := make([]*SSADPackage, 0)
	
	// Simple search implementation
	for _, pkg := range r.Packages {
		match := true
		
		if csp, ok := criteria["csp"]; ok && pkg.Metadata.CSPName != csp {
			match = false
		}
		
		if level, ok := criteria["impact_level"]; ok && pkg.Metadata.ImpactLevel != level {
			match = false
		}
		
		if authType, ok := criteria["auth_type"]; ok && pkg.Metadata.AuthorizationType != authType {
			match = false
		}
		
		if match {
			results = append(results, pkg)
		}
	}
	
	return results
}

// ExportPackage exports a package for distribution
func (r *SSADRepository) ExportPackage(packageID string) ([]byte, error) {
	pkg, exists := r.Packages[packageID]
	if !exists {
		return nil, fmt.Errorf("package not found: %s", packageID)
	}
	
	return json.MarshalIndent(pkg, "", "  ")
} 