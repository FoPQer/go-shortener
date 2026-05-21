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

type Action string

const (
	ActionShorten Action = "shorten"
	ActionFollow  Action = "follow"
)

type AuditEvent struct {
	Action Action
	UserID utils.UserID
	URL    string
}

type AuditDTO struct {
	Timestamp int64        `json:"ts"`
	Action    Action       `json:"action"`
	UserID    utils.UserID `json:"user_id,omitempty"`
	URL       string       `json:"url,omitempty"`
}

type Publisher interface {
	AddSubscriber(a Auditor)
	Publish(event AuditEvent)
}

type AuditBus struct {
	mu          sync.RWMutex
	subscribers []chan AuditEvent
	buffer      int
}

func NewAuditBus(buffer int) *AuditBus {
	if buffer < 1 {
		buffer = 1
	}

	return &AuditBus{buffer: buffer}
}

func (b *AuditBus) AddSubscriber(a Auditor) {
	ch := make(chan AuditEvent, b.buffer)

	b.mu.Lock()
	b.subscribers = append(b.subscribers, ch)
	b.mu.Unlock()

	go a.Subscribe(ch)
}

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

type Auditor interface {
	Subscribe(events <-chan AuditEvent)
}

type AuditFile struct {
	ID       int
	FilePath string
}

func NewAuditFile(id int, filePath string) *AuditFile {
	return &AuditFile{
		ID:       id,
		FilePath: filePath,
	}
}

func (s *AuditFile) Subscribe(events <-chan AuditEvent) {
	for event := range events {
		if err := s.WriteAudit(event); err != nil {
			logger.GetSugar().Errorf("AuditFile write error, event: %+v, file path: %s, err: %v", event, s.FilePath, err)
			continue
		}
		logger.GetSugar().Infof("AuditFile processed event, event: %+v, file path: %s", event, s.FilePath)
	}
}

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

type AuditURL struct {
	ID     int
	URL    string
	Events chan AuditEvent
}

func NewAuditURL(id int, url string) *AuditURL {
	return &AuditURL{
		ID:     id,
		URL:    url,
		Events: make(chan AuditEvent, 100),
	}
}

func (s *AuditURL) Subscribe(events <-chan AuditEvent) {
	for event := range events {
		logger.GetSugar().Infof("AuditURL received event, event: %+v, id: %v, url: %s", event, s.ID, s.URL)
	}
}

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

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
