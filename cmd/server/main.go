package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocomply/fedramp/pkg/api"
	"github.com/gocomply/fedramp/pkg/database"
	"github.com/gocomply/fedramp/pkg/monitor"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Parse command line flags
	var (
		port         = flag.String("port", "8080", "API server port")
		dbHost       = flag.String("db-host", "localhost", "Database host")
		dbPort       = flag.Int("db-port", 5432, "Database port")
		dbUser       = flag.String("db-user", "fedramp", "Database user")
		dbPassword   = flag.String("db-password", "", "Database password")
		dbName       = flag.String("db-name", "fedramp", "Database name")
		enableAuth   = flag.Bool("enable-auth", false, "Enable authentication")
		enableDash   = flag.Bool("enable-dashboard", true, "Enable web dashboard")
		logLevel     = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Configure logging
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})

	log.Info("Starting FedRAMP R5 Balance & 20x Server")

	// Initialize database
	dbConfig := &database.Config{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		Database: *dbName,
		SSLMode:  "disable",
	}

	// If no password provided, try environment variable
	if dbConfig.Password == "" {
		dbConfig.Password = os.Getenv("DB_PASSWORD")
	}

	// For development, use SQLite if no database configured
	var db *database.DB
	if dbConfig.Host == "localhost" && dbConfig.Password == "" {
		log.Info("Using in-memory database for development")
		// TODO: Implement SQLite support for development
	} else {
		db, err = database.NewDB(dbConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
	}

	// Initialize API server
	apiConfig := &api.Config{
		Port:            *port,
		DatabaseURL:     fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database),
		EnableAuth:      *enableAuth,
		EnableMetrics:   true,
		EnableDashboard: *enableDash,
	}

	server := api.NewServer(apiConfig)

	// Initialize continuous monitoring
	monitorConfig := &monitor.Config{
		CheckInterval:   5 * time.Minute,
		MetricInterval:  1 * time.Minute,
		AlertThreshold:  95.0,
		EnabledChecks:   []string{"ksi", "vulnerability", "configuration", "access"},
		NotificationURL: os.Getenv("NOTIFICATION_WEBHOOK_URL"),
	}

	continuousMonitor := monitor.NewContinuousMonitor(db, monitorConfig)
	
	// Start monitoring
	if err := continuousMonitor.Start(); err != nil {
		log.Fatalf("Failed to start continuous monitoring: %v", err)
	}

	// Start API server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	log.Infof("Server started on port %s", *port)
	if *enableDash {
		log.Infof("Dashboard available at http://localhost:%s/dashboard/", *port)
	}
	log.Info("API available at http://localhost:" + *port + "/api/v1/health")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down server...")
	
	// Graceful shutdown
	if err := continuousMonitor.Stop(); err != nil {
		log.Errorf("Error stopping continuous monitor: %v", err)
	}

	log.Info("Server shutdown complete")
}

// Development mode helpers
func isDevelopment() bool {
	return os.Getenv("FEDRAMP_ENV") == "development" || os.Getenv("FEDRAMP_ENV") == ""
}

func setupDevelopmentData(db *database.DB) {
	log.Info("Setting up development data...")
	
	// Create sample CSOs
	sampleCSOs := []string{"CSO-001", "CSO-002", "CSO-003"}
	for _, csoID := range sampleCSOs {
		// Create KSI evidence
		evidence := map[string]bool{
			"KSI-CED-01": true,
			"KSI-CED-02": true,
			"KSI-CMT-01": true,
			"KSI-CMT-02": true,
			"KSI-CMT-03": true,
			"KSI-CMT-04": true,
			"KSI-CMT-05": false, // One failing for demo
		}
		
		for ksiID, status := range evidence {
			db.SaveKSIEvidence(csoID, ksiID, map[string]interface{}{
				"status": status,
				"evidence": "Development test data",
			})
		}
	}
	
	log.Info("Development data setup complete")
} 