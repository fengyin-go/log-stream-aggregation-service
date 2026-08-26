package model

import (
	"strings"
	"time"
)

type Query struct {
	ID        string    `json:"id"`
	SourceID  string    `json:"source_id"`
	Expression string   `json:"expression"`
	CreatedBy string    `json:"created_by"`
	ExecutedAt time.Time `json:"executed_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Query) Validate() error {
	q.Expression = strings.TrimSpace(q.Expression)
	q.CreatedBy = strings.TrimSpace(q.CreatedBy)
	if q.Expression == "" {
		return NewValidationError("expression", "查询表达式不能为空")
	}
	if q.CreatedBy == "" {
		return NewValidationError("created_by", "创建人不能为空")
	}
	if q.ExecutedAt.IsZero() {
		q.ExecutedAt = time.Now()
	}
	return nil
}

type QueryFilter struct {
	SourceID string
	CreatedBy string
	Keyword  string
}

func (f QueryFilter) Match(q *Query) bool {
	if f.SourceID != "" && q.SourceID != f.SourceID {
		return false
	}
	if f.CreatedBy != "" && q.CreatedBy != f.CreatedBy {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(q.Expression), k) {
			return false
		}
	}
	return true
}
