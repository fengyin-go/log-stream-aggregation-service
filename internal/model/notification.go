package model

import (
	"strings"
	"time"
)

const (
	NotificationTypeAlert     = "alert"
	NotificationTypeBatch     = "batch"
	NotificationTypeSystem    = "system"
	NotificationStatusUnread  = "unread"
	NotificationStatusRead    = "read"
)

type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Notification) Validate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Message = strings.TrimSpace(n.Message)
	if n.Title == "" {
		return NewValidationError("title", "通知标题不能为空")
	}
	if n.Message == "" {
		return NewValidationError("message", "通知内容不能为空")
	}
	if n.Type == "" {
		n.Type = NotificationTypeSystem
	}
	if n.Type != NotificationTypeAlert && n.Type != NotificationTypeBatch && n.Type != NotificationTypeSystem {
		return NewValidationError("type", "通知类型不合法")
	}
	if n.Status == "" {
		n.Status = NotificationStatusUnread
	}
	if n.Status != NotificationStatusUnread && n.Status != NotificationStatusRead {
		return NewValidationError("status", "通知状态不合法")
	}
	return nil
}

type NotificationFilter struct {
	Type    string
	Status  string
	Keyword string
}

func (f NotificationFilter) Match(n *Notification) bool {
	if f.Type != "" && n.Type != f.Type {
		return false
	}
	if f.Status != "" && n.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(n.Title), k) &&
			!strings.Contains(strings.ToLower(n.Message), k) {
			return false
		}
	}
	return true
}
