package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Alert represents a compliance alert
type Alert struct {
	ID          string      `json:"id"`
	Severity    string      `json:"severity"` // critical, high, medium, low
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CSOId       string      `json:"csoId"`
	Timestamp   time.Time   `json:"timestamp"`
	Violations  []Violation `json:"violations,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AlertManager handles alert processing and notifications
type AlertManager struct {
	alerts          chan *Alert
	notificationURL string
	handlers        map[string]AlertHandler
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
}

// AlertHandler interface for custom alert processing
type AlertHandler interface {
	Handle(alert *Alert) error
	Name() string
}

// NewAlertManager creates a new alert manager
func NewAlertManager(notificationURL string) *AlertManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	am := &AlertManager{
		alerts:          make(chan *Alert, 1000),
		notificationURL: notificationURL,
		handlers:        make(map[string]AlertHandler),
		ctx:             ctx,
		cancel:          cancel,
	}

	// Register default handlers
	am.RegisterHandler("webhook", &WebhookHandler{url: notificationURL})
	am.RegisterHandler("email", &EmailHandler{})
	am.RegisterHandler("slack", &SlackHandler{})
	am.RegisterHandler("pagerduty", &PagerDutyHandler{})

	return am
}

// RegisterHandler adds a new alert handler
func (am *AlertManager) RegisterHandler(name string, handler AlertHandler) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.handlers[name] = handler
}

// SendAlert queues an alert for processing
func (am *AlertManager) SendAlert(alert *Alert) {
	alert.ID = fmt.Sprintf("ALERT-%d", time.Now().UnixNano())
	
	select {
	case am.alerts <- alert:
		log.Infof("Alert queued: %s - %s", alert.ID, alert.Title)
	default:
		log.Errorf("Alert queue full, dropping alert: %s", alert.Title)
	}
}

// Start begins processing alerts
func (am *AlertManager) Start() {
	am.wg.Add(1)
	go am.processAlerts()
}

// Stop gracefully stops the alert manager
func (am *AlertManager) Stop() error {
	am.cancel()
	close(am.alerts)
	am.wg.Wait()
	return nil
}

// processAlerts handles alerts from the queue
func (am *AlertManager) processAlerts() {
	defer am.wg.Done()

	for {
		select {
		case alert, ok := <-am.alerts:
			if !ok {
				return
			}
			am.handleAlert(alert)
		case <-am.ctx.Done():
			// Process remaining alerts
			for alert := range am.alerts {
				am.handleAlert(alert)
			}
			return
		}
	}
}

// handleAlert processes a single alert through all handlers
func (am *AlertManager) handleAlert(alert *Alert) {
	log.Infof("Processing alert: %s", alert.ID)

	am.mu.RLock()
	handlers := make(map[string]AlertHandler)
	for k, v := range am.handlers {
		handlers[k] = v
	}
	am.mu.RUnlock()

	var wg sync.WaitGroup
	for name, handler := range handlers {
		wg.Add(1)
		go func(n string, h AlertHandler, a *Alert) {
			defer wg.Done()
			
			if err := h.Handle(a); err != nil {
				log.Errorf("Handler %s failed for alert %s: %v", n, a.ID, err)
			} else {
				log.Debugf("Handler %s processed alert %s successfully", n, a.ID)
			}
		}(name, handler, alert)
	}
	wg.Wait()

	// Store alert in database for audit trail
	am.storeAlert(alert)
}

// storeAlert persists alert to database
func (am *AlertManager) storeAlert(alert *Alert) {
	// TODO: Implement database storage
	log.Debugf("Storing alert %s in database", alert.ID)
}

// Alert Handlers

// WebhookHandler sends alerts via HTTP webhook
type WebhookHandler struct {
	url string
}

func (h *WebhookHandler) Name() string {
	return "Webhook Handler"
}

func (h *WebhookHandler) Handle(alert *Alert) error {
	if h.url == "" {
		return nil
	}

	data, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	resp, err := http.Post(h.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// EmailHandler sends alerts via email
type EmailHandler struct {
	smtpHost string
	smtpPort int
	from     string
	to       []string
}

func (h *EmailHandler) Name() string {
	return "Email Handler"
}

func (h *EmailHandler) Handle(alert *Alert) error {
	// TODO: Implement email sending
	log.Debugf("Would send email for alert %s", alert.ID)
	return nil
}

// SlackHandler sends alerts to Slack
type SlackHandler struct {
	webhookURL string
}

func (h *SlackHandler) Name() string {
	return "Slack Handler"
}

func (h *SlackHandler) Handle(alert *Alert) error {
	if h.webhookURL == "" {
		return nil
	}

	// Format alert for Slack
	slackMsg := map[string]interface{}{
		"text": fmt.Sprintf("*%s Alert*: %s", alert.Severity, alert.Title),
		"attachments": []map[string]interface{}{
			{
				"color": h.getSeverityColor(alert.Severity),
				"fields": []map[string]interface{}{
					{
						"title": "CSO ID",
						"value": alert.CSOId,
						"short": true,
					},
					{
						"title": "Time",
						"value": alert.Timestamp.Format(time.RFC3339),
						"short": true,
					},
					{
						"title": "Description",
						"value": alert.Description,
						"short": false,
					},
				},
			},
		},
	}

	data, err := json.Marshal(slackMsg)
	if err != nil {
		return err
	}

	resp, err := http.Post(h.webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (h *SlackHandler) getSeverityColor(severity string) string {
	switch severity {
	case "critical":
		return "#FF0000"
	case "high":
		return "#FF8C00"
	case "medium":
		return "#FFD700"
	case "low":
		return "#00FF00"
	default:
		return "#808080"
	}
}

// PagerDutyHandler sends alerts to PagerDuty
type PagerDutyHandler struct {
	apiKey      string
	routingKey  string
}

func (h *PagerDutyHandler) Name() string {
	return "PagerDuty Handler"
}

func (h *PagerDutyHandler) Handle(alert *Alert) error {
	if h.apiKey == "" || alert.Severity != "critical" {
		return nil
	}

	// Create PagerDuty event
	event := map[string]interface{}{
		"routing_key":  h.routingKey,
		"event_action": "trigger",
		"dedup_key":    alert.ID,
		"payload": map[string]interface{}{
			"summary":  alert.Title,
			"source":   "fedramp-monitor",
			"severity": h.mapSeverity(alert.Severity),
			"custom_details": map[string]interface{}{
				"cso_id":      alert.CSOId,
				"description": alert.Description,
				"violations":  alert.Violations,
			},
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://events.pagerduty.com/v2/enqueue", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", h.apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (h *PagerDutyHandler) mapSeverity(severity string) string {
	switch severity {
	case "critical":
		return "critical"
	case "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "info"
	}
}

// AlertAggregator groups similar alerts to reduce noise
type AlertAggregator struct {
	window    time.Duration
	threshold int
	alerts    map[string][]*Alert
	mu        sync.Mutex
}

// NewAlertAggregator creates a new alert aggregator
func NewAlertAggregator(window time.Duration, threshold int) *AlertAggregator {
	return &AlertAggregator{
		window:    window,
		threshold: threshold,
		alerts:    make(map[string][]*Alert),
	}
}

// Add adds an alert to the aggregator
func (aa *AlertAggregator) Add(alert *Alert) bool {
	aa.mu.Lock()
	defer aa.mu.Unlock()

	key := fmt.Sprintf("%s-%s-%s", alert.CSOId, alert.Severity, alert.Title)
	
	// Clean old alerts
	cutoff := time.Now().Add(-aa.window)
	filtered := []*Alert{}
	for _, a := range aa.alerts[key] {
		if a.Timestamp.After(cutoff) {
			filtered = append(filtered, a)
		}
	}
	aa.alerts[key] = filtered

	// Add new alert
	aa.alerts[key] = append(aa.alerts[key], alert)

	// Check if threshold exceeded
	return len(aa.alerts[key]) >= aa.threshold
} 