package events

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/FoPQer/go-shortener/internal/utils"
)

// Action describes the type of user interaction captured by audit logging.
type Action string

const (
	// ActionShorten indicates short URL creation.
	ActionShorten Action = "shorten"
	// ActionFollow indicates opening a short URL and following redirect.
	ActionFollow Action = "follow"
)

// AuditEvent is an in-process event published by handlers and services.
type AuditEvent struct {
	Action Action
	UserID utils.UserID
	URL    string
}

// AuditDTO is a serialized audit record used for file and HTTP transport.
type AuditDTO struct {
	Timestamp int64        `json:"ts"`
	Action    Action       `json:"action"`
	UserID    utils.UserID `json:"user_id,omitempty"`
	URL       string       `json:"url,omitempty"`
}

// Publisher broadcasts audit events to registered subscribers.
type Publisher interface {
	AddSubscriber(a Auditor)
	Publish(event AuditEvent)
}

// AuditBus is an in-memory pub/sub bus for audit events.
type AuditBus struct {
	mu          sync.RWMutex
	subscribers []chan AuditEvent
	buffer      int
}

// NewAuditBus creates an AuditBus with a per-subscriber channel buffer.
//
// Buffer values less than 1 are normalized to 1.
func NewAuditBus(buffer int) *AuditBus {
	if buffer < 1 {
		buffer = 1
	}

	return &AuditBus{buffer: buffer}
}

// AddSubscriber registers an auditor and starts asynchronous event consumption.
func (b *AuditBus) AddSubscriber(a Auditor) {
	ch := make(chan AuditEvent, b.buffer)

	b.mu.Lock()
	b.subscribers = append(b.subscribers, ch)
	b.mu.Unlock()

	go a.Subscribe(ch)
}

// Publish sends an event to all subscribers in a non-blocking manner.
//
// Events are dropped when a subscriber buffer is full.
func (b *AuditBus) Publish(event AuditEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
			logger.GetSugar().Warnf("Audit event dropped due to full subscriber buffer, event: %+v", event)
		}
	}
}

// Close closes all subscriber channels, causing their goroutines to finish.
func (b *AuditBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil
}

// Auditor consumes audit events from a subscription channel.
type Auditor interface {
	Subscribe(events <-chan AuditEvent)
}

// AuditFile writes audit events to a local file as JSON lines.
type AuditFile struct {
	ID       int
	FilePath string
}

// NewAuditFile creates a file-backed auditor.
func NewAuditFile(id int, filePath string) *AuditFile {
	return &AuditFile{
		ID:       id,
		FilePath: filePath,
	}
}

// Subscribe consumes events and persists each one to file storage.
func (s *AuditFile) Subscribe(events <-chan AuditEvent) {
	for event := range events {
		if err := s.WriteAudit(event); err != nil {
			logger.GetSugar().Errorf("AuditFile write error, event: %+v, file path: %s, err: %v", event, s.FilePath, err)
			continue
		}
		logger.GetSugar().Infof("AuditFile processed event, event: %+v, file path: %s", event, s.FilePath)
	}
}

// WriteAudit appends a serialized audit event to the configured file.
func (s *AuditFile) WriteAudit(event AuditEvent) error {
	auditDTO := AuditDTO{
		Timestamp: getCurrentTimestamp(),
		Action:    event.Action,
		UserID:    event.UserID,
		URL:       event.URL,
	}
	file, err := os.OpenFile(s.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	json.NewEncoder(file).Encode(auditDTO)

	return nil
}

// AuditURL sends audit events to a remote HTTP endpoint.
type AuditURL struct {
	ID     int
	URL    string
	Events chan AuditEvent
}

// NewAuditURL creates an HTTP-backed auditor with an internal event buffer.
func NewAuditURL(id int, url string) *AuditURL {
	return &AuditURL{
		ID:     id,
		URL:    url,
		Events: make(chan AuditEvent, 100),
	}
}

// Subscribe consumes events from a subscription channel.
//
// Current implementation logs received events.
func (s *AuditURL) Subscribe(events <-chan AuditEvent) {
	for event := range events {
		logger.GetSugar().Infof("AuditURL received event, event: %+v, id: %v, url: %s", event, s.ID, s.URL)
	}
}

// SendAudit serializes and sends an audit event to the configured remote endpoint.
func (s *AuditURL) SendAudit(event AuditEvent) {
	auditDTO := AuditDTO{
		Timestamp: getCurrentTimestamp(),
		Action:    event.Action,
		UserID:    event.UserID,
		URL:       event.URL,
	}
	select {
	case s.Events <- event:
		client := &http.Client{}
		clientURL := service.GetAuditURL()
		auditDTOBytes, err := json.Marshal(auditDTO)
		if err != nil {
			logger.GetSugar().Errorf("Error marshaling auditDTO, event: %+v, id: %v, url: %s, err: %v", event, s.ID, s.URL, err)
			return
		}
		req, err := http.NewRequest(http.MethodPost, clientURL, bytes.NewBuffer(auditDTOBytes))
		if err != nil {
			logger.GetSugar().Errorf("Error creating request for AuditURL, event: %+v, id: %v, url: %s, err: %v", event, s.ID, s.URL, err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			logger.GetSugar().Errorf("Error sending request for AuditURL, event: %+v, id: %v, url: %s, err: %v", event, s.ID, s.URL, err)
			return
		}
		defer resp.Body.Close()
	default:
		logger.GetSugar().Warnf("Audit event dropped due to full channel buffer, event: %+v, id: %v, url: %s", event, s.ID, s.URL)
	}
}

// getCurrentTimestamp returns the current Unix timestamp in seconds.
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
