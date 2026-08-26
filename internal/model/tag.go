package model

import (
	"strings"
	"time"
)

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tag) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return NewValidationError("name", "标签名称不能为空")
	}
	return nil
}

type TagFilter struct {
	Name    string
	Keyword string
}

func (f TagFilter) Match(t *Tag) bool {
	if f.Name != "" && t.Name != f.Name {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) {
			return false
		}
	}
	return true
}
