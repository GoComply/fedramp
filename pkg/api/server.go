package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocomply/fedramp/pkg/fedramp"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Server represents the FedRAMP API server
type Server struct {
	router *mux.Router
	config *Config
}

// Config holds server configuration
type Config struct {
	Port            string
	DatabaseURL     string
	EnableAuth      bool
	EnableMetrics   bool
	EnableDashboard bool
}

// NewServer creates a new API server instance
func NewServer(config *Config) *Server {
	s := &Server{
		router: mux.NewRouter(),
		config: config,
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures all API endpoints
func (s *Server) setupRoutes() {
	// API versioning
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", s.healthCheck).Methods("GET")

	// KSI endpoints
	api.HandleFunc("/ksi/validate", s.validateKSI).Methods("POST")
	api.HandleFunc("/ksi/report/{csoId}", s.getKSIReport).Methods("GET")
	api.HandleFunc("/ksi/evidence/{csoId}", s.submitEvidence).Methods("POST")
	api.HandleFunc("/ksi/continuous/{csoId}", s.continuousMonitoring).Methods("GET")

	// SCN endpoints
	api.HandleFunc("/scn", s.createSCN).Methods("POST")
	api.HandleFunc("/scn/{csoId}", s.listSCNs).Methods("GET")
	api.HandleFunc("/scn/{csoId}/{scnId}", s.getSCN).Methods("GET")
	api.HandleFunc("/scn/{csoId}/{scnId}/approve", s.approveSCN).Methods("POST")

	// CRS endpoints
	api.HandleFunc("/crs/report", s.createCRSReport).Methods("POST")
	api.HandleFunc("/crs/metrics/{csoId}", s.getMetrics).Methods("GET")
	api.HandleFunc("/crs/dashboard/{csoId}", s.getDashboard).Methods("GET")

	// MAS endpoints
	api.HandleFunc("/mas/assessment", s.createAssessment).Methods("POST")
	api.HandleFunc("/mas/assessment/{assessmentId}", s.getAssessment).Methods("GET")
	api.HandleFunc("/mas/findings", s.submitFindings).Methods("POST")

	// SSAD endpoints
	api.HandleFunc("/ssad/package", s.createPackage).Methods("POST")
	api.HandleFunc("/ssad/package/{packageId}", s.getPackage).Methods("GET")
	api.HandleFunc("/ssad/repository", s.listPackages).Methods("GET")

	// FRMR endpoints
	api.HandleFunc("/frmr/validate", s.validateFRMR).Methods("POST")
	api.HandleFunc("/frmr/transform", s.transformDocument).Methods("POST")

	// Dashboard and UI
	if s.config.EnableDashboard {
		s.router.PathPrefix("/dashboard/").Handler(http.StripPrefix("/dashboard/", http.FileServer(http.Dir("./web/dashboard"))))
	}

	// Metrics endpoint
	if s.config.EnableMetrics {
		api.HandleFunc("/metrics", s.getPrometheusMetrics).Methods("GET")
	}
}

// Start begins serving HTTP requests
func (s *Server) Start() error {
	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)

	// Add middleware
	handler = loggingMiddleware(handler)
	if s.config.EnableAuth {
		handler = authMiddleware(handler)
	}

	log.Infof("Starting FedRAMP API server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, handler)
}

// Health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"services": map[string]string{
			"ksi":  "operational",
			"scn":  "operational",
			"crs":  "operational",
			"mas":  "operational",
			"ssad": "operational",
			"frmr": "operational",
		},
	}
	respondJSON(w, http.StatusOK, health)
}

// KSI Endpoints

func (s *Server) validateKSI(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CSOId    string                 `json:"csoId"`
		Evidence map[string]interface{} `json:"evidence"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Perform KSI validation
	validation := fedramp.NewKSIValidation(req.CSOId)
	// TODO: Apply evidence to validation
	
	respondJSON(w, http.StatusOK, validation)
}

func (s *Server) getKSIReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	csoId := vars["csoId"]

	report := fedramp.GenerateKSIReport(csoId, time.Now())
	respondJSON(w, http.StatusOK, report)
}

func (s *Server) submitEvidence(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	csoId := vars["csoId"]

	var evidence map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&evidence); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid evidence format")
		return
	}

	// TODO: Store evidence in database
	log.Infof("Evidence submitted for CSO %s", csoId)

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "accepted",
		"csoId":  csoId,
	})
}

func (s *Server) continuousMonitoring(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	csoId := vars["csoId"]

	// TODO: Implement real-time monitoring data
	monitoring := map[string]interface{}{
		"csoId":      csoId,
		"status":     "monitoring",
		"lastCheck":  time.Now().UTC(),
		"compliance": 98.5,
		"alerts":     []string{},
	}

	respondJSON(w, http.StatusOK, monitoring)
}

// SCN Endpoints

func (s *Server) createSCN(w http.ResponseWriter, r *http.Request) {
	var scn fedramp.SignificantChangeNotification
	if err := json.NewDecoder(r.Body).Decode(&scn); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid SCN data")
		return
	}

	// Set metadata
	scn.CreatedAt = time.Now()
	scn.Status = "pending"

	// TODO: Store in database
	respondJSON(w, http.StatusCreated, scn)
}

func (s *Server) listSCNs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	csoId := vars["csoId"]

	// TODO: Fetch from database
	scns := []fedramp.SignificantChangeNotification{}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"csoId": csoId,
		"scns":  scns,
		"total": len(scns),
	})
}

func (s *Server) getSCN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scnId := vars["scnId"]

	// TODO: Fetch from database
	respondJSON(w, http.StatusNotFound, map[string]string{
		"error": fmt.Sprintf("SCN %s not found", scnId),
	})
}

func (s *Server) approveSCN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scnId := vars["scnId"]

	var approval struct {
		ApprovedBy string `json:"approvedBy"`
		Comments   string `json:"comments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&approval); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid approval data")
		return
	}

	// TODO: Update SCN status in database
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "approved",
		"scnId":  scnId,
	})
}

// CRS Endpoints

func (s *Server) createCRSReport(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Extract CSO ID from request
	csoId, _ := req["cso_id"].(string)
	reportingPeriod, _ := req["reporting_period"].(string)

	// Create metrics with proper structure
	metrics := fedramp.KeySecurityMetrics{
		VulnerabilityScanning: fedramp.VulnerabilityMetric{
			ScansCompleted:  12,
			CriticalFindings: 0,
			HighFindings:    2,
			MediumFindings:  5,
			LowFindings:     10,
			LastScanDate:    time.Now(),
			RemediationRate: 95.5,
		},
		SecurityIncidents: fedramp.IncidentMetric{
			TotalIncidents:        3,
			OpenIncidents:         0,
			ClosedIncidents:       3,
			AverageResolutionTime: 4.2,
			IncidentsByCategory:   map[string]int{"phishing": 2, "malware": 1},
		},
		ConfigurationManagement: fedramp.ConfigurationMetric{
			BaselineCompliance:  98.5,
			UnauthorizedChanges: 0,
			ConfigurationDrift:  2,
			LastBaselineReview:  time.Now(),
		},
		AccessControl: fedramp.AccessControlMetric{
			ActiveUsers:         150,
			PrivilegedUsers:     12,
			FailedLoginAttempts: 8,
			AccountLockouts:     1,
			MFAAdoptionRate:     100.0,
		},
		SystemAvailability: fedramp.AvailabilityMetric{
			UptimePercentage:  99.95,
			PlannedDowntime:   2.0,
			UnplannedDowntime: 0.5,
			MTTRHours:         0.25,
		},
		PatchManagement: fedramp.PatchMetric{
			PatchesAvailable:    15,
			PatchesApplied:      15,
			CriticalPatches:     2,
			PatchComplianceRate: 100.0,
			AveragePatchAge:     3.5,
		},
	}

	// Create CRS report
	crs := &fedramp.ContinuousReportingStandard{
		CSOId:           csoId,
		ReportingPeriod: reportingPeriod,
		GeneratedAt:     time.Now(),
		Metrics:         metrics,
		ComplianceScore: 98.5,
		Status:          "compliant",
	}

	// TODO: Store report
	respondJSON(w, http.StatusCreated, crs)
}

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	csoId := vars["csoId"]

	// TODO: Fetch from metrics database
	metrics := fedramp.KeySecurityMetrics{
		VulnerabilityScanning: fedramp.VulnerabilityMetric{
			ScansCompleted:  52,
			CriticalFindings: 0,
			HighFindings:    1,
			MediumFindings:  8,
			LowFindings:     23,
			LastScanDate:    time.Now(),
			RemediationRate: 98.5,
		},
		SecurityIncidents: fedramp.IncidentMetric{
			TotalIncidents:        12,
			OpenIncidents:         0,
			ClosedIncidents:       12,
			AverageResolutionTime: 3.5,
			IncidentsByCategory:   map[string]int{"phishing": 8, "malware": 2, "other": 2},
		},
		ConfigurationManagement: fedramp.ConfigurationMetric{
			BaselineCompliance:  99.2,
			UnauthorizedChanges: 0,
			ConfigurationDrift:  1,
			LastBaselineReview:  time.Now().AddDate(0, 0, -7),
		},
		AccessControl: fedramp.AccessControlMetric{
			ActiveUsers:         250,
			PrivilegedUsers:     15,
			FailedLoginAttempts: 12,
			AccountLockouts:     2,
			MFAAdoptionRate:     95.8,
		},
		SystemAvailability: fedramp.AvailabilityMetric{
			UptimePercentage:  99.99,
			PlannedDowntime:   4.0,
			UnplannedDowntime: 0.1,
			MTTRHours:         0.15,
		},
		PatchManagement: fedramp.PatchMetric{
			PatchesAvailable:    8,
			PatchesApplied:      8,
			CriticalPatches:     0,
			PatchComplianceRate: 100.0,
			AveragePatchAge:     2.1,
		},
	}

	// Wrap in response with CSO ID
	response := map[string]interface{}{
		"cso_id":  csoId,
		"metrics": metrics,
		"timestamp": time.Now(),
	}

	respondJSON(w, http.StatusOK, response)
}

func (s *Server) getDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	csoId := vars["csoId"]

	// TODO: Aggregate dashboard data
	dashboard := map[string]interface{}{
		"csoId":           csoId,
		"complianceScore": 98.5,
		"activeAlerts":    0,
		"lastAssessment":  time.Now().AddDate(0, -1, 0),
		"nextAssessment":  time.Now().AddDate(0, 11, 0),
	}

	respondJSON(w, http.StatusOK, dashboard)
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Middleware

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Infof("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement proper authentication
		token := r.Header.Get("Authorization")
		if token == "" {
			respondError(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MAS, SSAD, and FRMR endpoints would follow similar patterns...
// Truncating for brevity, but would include all documented endpoints

// MAS Endpoints

func (s *Server) createAssessment(w http.ResponseWriter, r *http.Request) {
	var assessment map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&assessment); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid assessment data")
		return
	}

	assessment["id"] = fmt.Sprintf("MAS-%d", time.Now().Unix())
	assessment["createdAt"] = time.Now()
	assessment["status"] = "in-progress"

	respondJSON(w, http.StatusCreated, assessment)
}

func (s *Server) getAssessment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assessmentId := vars["assessmentId"]

	// Mock response
	assessment := map[string]interface{}{
		"id":     assessmentId,
		"status": "completed",
		"type":   "initial",
		"csoId":  "CSO-001",
	}

	respondJSON(w, http.StatusOK, assessment)
}

func (s *Server) submitFindings(w http.ResponseWriter, r *http.Request) {
	var findings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&findings); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid findings data")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "findings recorded",
		"id":     fmt.Sprintf("FIND-%d", time.Now().Unix()),
	})
}

// SSAD Endpoints

func (s *Server) createPackage(w http.ResponseWriter, r *http.Request) {
	var pkg map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid package data")
		return
	}

	pkg["id"] = fmt.Sprintf("PKG-%d", time.Now().Unix())
	pkg["createdAt"] = time.Now()
	pkg["status"] = "active"

	respondJSON(w, http.StatusCreated, pkg)
}

func (s *Server) getPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageId := vars["packageId"]

	pkg := map[string]interface{}{
		"id":     packageId,
		"status": "active",
		"type":   "ssp",
		"csoId":  "CSO-001",
	}

	respondJSON(w, http.StatusOK, pkg)
}

func (s *Server) listPackages(w http.ResponseWriter, r *http.Request) {
	packages := []map[string]interface{}{
		{
			"id":     "PKG-001",
			"name":   "SSP Package",
			"status": "active",
			"csoId":  "CSO-001",
		},
		{
			"id":     "PKG-002",
			"name":   "SAR Package",
			"status": "active",
			"csoId":  "CSO-002",
		},
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"packages": packages,
		"total":    len(packages),
	})
}

// FRMR Endpoints

func (s *Server) validateFRMR(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid validation request")
		return
	}

	validation := map[string]interface{}{
		"valid":   true,
		"errors":  []string{},
		"warnings": []string{},
		"timestamp": time.Now(),
	}

	respondJSON(w, http.StatusOK, validation)
}

func (s *Server) transformDocument(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid transform request")
		return
	}

	transformed := map[string]interface{}{
		"status": "transformed",
		"format": "oscal",
		"timestamp": time.Now(),
	}

	respondJSON(w, http.StatusOK, transformed)
}

// Metrics Endpoint

func (s *Server) getPrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	// Simple Prometheus metrics format
	metrics := `# HELP fedramp_compliance_score Current compliance score
# TYPE fedramp_compliance_score gauge
fedramp_compliance_score 98.5

# HELP fedramp_active_csos Number of active CSOs
# TYPE fedramp_active_csos gauge
fedramp_active_csos 12

# HELP fedramp_open_alerts Number of open alerts
# TYPE fedramp_open_alerts gauge
fedramp_open_alerts 2

# HELP fedramp_ksi_compliant Number of compliant KSIs
# TYPE fedramp_ksi_compliant gauge
fedramp_ksi_compliant 11
`
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
} 