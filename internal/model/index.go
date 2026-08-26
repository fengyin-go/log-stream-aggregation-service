package model

import (
	"strings"
	"time"
)

type Index struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SourceID  string    `json:"source_id"`
	Field     string    `json:"field"`
	CreatedAt time.Time `json:"created_at"`
}

func (i *Index) Validate() error {
	i.Name = strings.TrimSpace(i.Name)
	i.Field = strings.TrimSpace(i.Field)
	if i.Name == "" {
		return NewValidationError("name", "索引名称不能为空")
	}
	if i.SourceID == "" {
		return NewValidationError("source_id", "日志源 ID 不能为空")
	}
	if i.Field == "" {
		return NewValidationError("field", "索引字段不能为空")
	}
	return nil
}

type IndexFilter struct {
	Name     string
	SourceID string
	Field    string
	Keyword  string
}

func (f IndexFilter) Match(i *Index) bool {
	if f.Name != "" && i.Name != f.Name {
		return false
	}
	if f.SourceID != "" && i.SourceID != f.SourceID {
		return false
	}
	if f.Field != "" && i.Field != f.Field {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(i.Name), k) {
			return false
		}
	}
	return true
}
