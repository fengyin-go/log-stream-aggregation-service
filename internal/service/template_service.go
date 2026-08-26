package service

import (
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

type AlertRuleTemplate struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	LevelThreshold string `json:"level_threshold"`
	Keyword        string `json:"keyword"`
	Description    string `json:"description"`
}

var defaultTemplates = []AlertRuleTemplate{
	{ID: "tpl-oom", Name: "内存溢出检测", LevelThreshold: model.LogLevelError, Keyword: "OutOfMemory", Description: "检测 JVM 或系统 OOM 错误"},
	{ID: "tpl-panic", Name: "Panic 检测", LevelThreshold: model.LogLevelFatal, Keyword: "panic", Description: "检测 Go panic 或系统崩溃"},
	{ID: "tpl-slow", Name: "慢查询检测", LevelThreshold: model.LogLevelWarn, Keyword: "slow", Description: "检测慢查询日志"},
	{ID: "tpl-404", Name: "404 异常检测", LevelThreshold: model.LogLevelWarn, Keyword: "404", Description: "检测大量 404 请求"},
	{ID: "tpl-timeout", Name: "超时检测", LevelThreshold: model.LogLevelError, Keyword: "timeout", Description: "检测超时错误"},
}

func (s *Service) ListAlertRuleTemplates() []AlertRuleTemplate {
	result := make([]AlertRuleTemplate, len(defaultTemplates))
	copy(result, defaultTemplates)
	return result
}

func (s *Service) GetAlertRuleTemplate(templateID string) (*AlertRuleTemplate, error) {
	for _, t := range defaultTemplates {
		if t.ID == templateID {
			return &AlertRuleTemplate{
				ID:             t.ID,
				Name:           t.Name,
				LevelThreshold: t.LevelThreshold,
				Keyword:        t.Keyword,
				Description:    t.Description,
			}, nil
		}
	}
	return nil, model.NewValidationError("template_id", "模板不存在")
}

func (s *Service) CreateAlertRuleFromTemplate(sourceID, templateID string) (*model.AlertRule, error) {
	tpl, err := s.GetAlertRuleTemplate(templateID)
	if err != nil {
		return nil, err
	}
	if _, err := s.store.GetLogSource(sourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	rule := &model.AlertRule{
		ID:             idgen.Hex(),
		Name:           tpl.Name,
		SourceID:       sourceID,
		LevelThreshold: tpl.LevelThreshold,
		Keyword:        tpl.Keyword,
		Status:         model.AlertRuleStatusActive,
		CreatedAt:      time.Now(),
	}
	if err := s.store.CreateAlertRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}
