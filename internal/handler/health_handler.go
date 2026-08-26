package handler

import (
	"net/http"

	"log-aggregation/pkg/httpx"
)

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sources/{id}/health", s.checkSourceHealth)
	mux.HandleFunc("GET /api/health/all", s.checkAllSourcesHealth)
	mux.HandleFunc("GET /api/health/system", s.getSystemHealth)
	mux.HandleFunc("GET /api/alert-rule-templates", s.listAlertRuleTemplates)
	mux.HandleFunc("GET /api/alert-rule-templates/{id}", s.getAlertRuleTemplate)
	mux.HandleFunc("POST /api/sources/{id}/alert-rules/from-template", s.createAlertRuleFromTemplate)
}

func (s *Server) checkSourceHealth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.CheckSourceHealth(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) checkAllSourcesHealth(w http.ResponseWriter, r *http.Request) {
	results := s.svc.CheckAllSourcesHealth()
	httpx.OK(w, results)
}

func (s *Server) getSystemHealth(w http.ResponseWriter, r *http.Request) {
	health := s.svc.GetSystemHealth()
	httpx.OK(w, health)
}

func (s *Server) listAlertRuleTemplates(w http.ResponseWriter, r *http.Request) {
	templates := s.svc.ListAlertRuleTemplates()
	httpx.OK(w, templates)
}

func (s *Server) getAlertRuleTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tpl, err := s.svc.GetAlertRuleTemplate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, tpl)
}

type createFromTemplateRequest struct {
	TemplateID string `json:"template_id"`
}

func (s *Server) createAlertRuleFromTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createFromTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateAlertRuleFromTemplate(id, req.TemplateID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}
