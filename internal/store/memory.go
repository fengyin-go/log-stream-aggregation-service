package store

import (
	"sync"

	"log-aggregation/internal/model"
)

type MemoryStore struct {
	mu                sync.RWMutex
	logSources        map[string]*model.LogSource
	logEntries        map[string]*model.LogEntry
	indexes           map[string]*model.Index
	queries           map[string]*model.Query
	alertRules        map[string]*model.AlertRule
	alerts            map[string]*model.Alert
	tags              map[string]*model.Tag
	notifications     map[string]*model.Notification
	retentionPolicies map[string]*model.RetentionPolicy
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		logSources:        make(map[string]*model.LogSource),
		logEntries:        make(map[string]*model.LogEntry),
		indexes:           make(map[string]*model.Index),
		queries:           make(map[string]*model.Query),
		alertRules:        make(map[string]*model.AlertRule),
		alerts:            make(map[string]*model.Alert),
		tags:              make(map[string]*model.Tag),
		notifications:     make(map[string]*model.Notification),
		retentionPolicies: make(map[string]*model.RetentionPolicy),
	}
}

var _ Store = (*MemoryStore)(nil)
