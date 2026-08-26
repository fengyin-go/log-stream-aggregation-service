package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerAlertRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alert-rules", s.createAlertRule)
	mux.HandleFunc("GET /api/alert-rules", s.listAlertRules)
	mux.HandleFunc("GET /api/alert-rules/{id}", s.getAlertRule)
	mux.HandleFunc("PUT /api/alert-rules/{id}", s.updateAlertRule)
	mux.HandleFunc("DELETE /api/alert-rules/{id}", s.deleteAlertRule)
	mux.HandleFunc("POST /api/alert-rules/{id}/transition", s.transitionAlertRuleStatus)
}

type createAlertRuleRequest struct {
	Name           string `json:"name"`
	SourceID       string `json:"source_id"`
	LevelThreshold string `json:"level_threshold"`
	Keyword        string `json:"keyword"`
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var req createAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ar, err := s.svc.CreateAlertRule(model.AlertRule{Name: req.Name, SourceID: req.SourceID, LevelThreshold: req.LevelThreshold, Keyword: req.Keyword})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, ar)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertRuleFilter{
		Name:     r.URL.Query().Get("name"),
		SourceID: r.URL.Query().Get("source_id"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAlertRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ar, err := s.svc.GetAlertRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ar)
}

type updateAlertRuleRequest struct {
	Name           string `json:"name"`
	LevelThreshold string `json:"level_threshold"`
	Keyword        string `json:"keyword"`
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ar, err := s.svc.UpdateAlertRule(id, model.AlertRule{Name: req.Name, LevelThreshold: req.LevelThreshold, Keyword: req.Keyword})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ar)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlertRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionAlertRuleStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionAlertRuleStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionAlertRuleStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ar, err := s.svc.TransitionAlertRuleStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ar)
}
