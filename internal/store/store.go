// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"
	"time"

	"log-aggregation/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// LogSource
	CreateLogSource(s *model.LogSource) error
	GetLogSource(id string) (*model.LogSource, error)
	GetLogSourceByName(name string) (*model.LogSource, error)
	ListLogSources() []*model.LogSource
	UpdateLogSource(s *model.LogSource) error
	DeleteLogSource(id string) error

	// LogEntry
	CreateLogEntry(e *model.LogEntry) error
	BatchCreateLogEntries(entries []*model.LogEntry) error
	GetLogEntry(id string) (*model.LogEntry, error)
	ListLogEntries() []*model.LogEntry
	DeleteLogEntry(id string) error
	DeleteLogEntriesBySourceID(sourceID string) error

	// Index
	CreateIndex(i *model.Index) error
	GetIndex(id string) (*model.Index, error)
	ListIndexes() []*model.Index
	UpdateIndex(i *model.Index) error
	DeleteIndex(id string) error

	// Query
	CreateQuery(q *model.Query) error
	GetQuery(id string) (*model.Query, error)
	ListQueries() []*model.Query
	UpdateQuery(q *model.Query) error
	DeleteQuery(id string) error

	// AlertRule
	CreateAlertRule(r *model.AlertRule) error
	GetAlertRule(id string) (*model.AlertRule, error)
	ListAlertRules() []*model.AlertRule
	UpdateAlertRule(r *model.AlertRule) error
	DeleteAlertRule(id string) error
	ListActiveAlertRulesBySourceID(sourceID string) []*model.AlertRule

	// Alert
	CreateAlert(a *model.Alert) error
	GetAlert(id string) (*model.Alert, error)
	ListAlerts() []*model.Alert
	UpdateAlert(a *model.Alert) error
	DeleteAlert(id string) error

	// Tag
	CreateTag(t *model.Tag) error
	GetTag(id string) (*model.Tag, error)
	GetTagByName(name string) (*model.Tag, error)
	ListTags() []*model.Tag
	UpdateTag(t *model.Tag) error
	DeleteTag(id string) error

	// Notification
	CreateNotification(n *model.Notification) error
	GetNotification(id string) (*model.Notification, error)
	ListNotifications() []*model.Notification
	UpdateNotification(n *model.Notification) error
	DeleteNotification(id string) error
	CountUnreadNotifications() int

	// RetentionPolicy
	CreateRetentionPolicy(p *model.RetentionPolicy) error
	GetRetentionPolicy(id string) (*model.RetentionPolicy, error)
	GetRetentionPolicyBySourceID(sourceID string) (*model.RetentionPolicy, error)
	ListRetentionPolicies() []*model.RetentionPolicy
	UpdateRetentionPolicy(p *model.RetentionPolicy) error
	DeleteRetentionPolicy(id string) error

	// Stats helpers
	CountLogSources() int
	CountLogEntries() int
	CountLogEntriesBySourceID(sourceID string) int
	CountAlerts() int
	CountAlertsByStatus(status string) int
	CountLogEntriesByLevel(level string) int
	CountLogEntriesBySourceAndHour(sourceID string, hour time.Time) int
	CountLogEntriesBySourceAndLevelAndHour(sourceID, level string, hour time.Time) int
	CountAlertsBySourceID(sourceID string) int
}
