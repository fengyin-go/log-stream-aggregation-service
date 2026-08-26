package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerAlertRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alerts", s.createAlert)
	mux.HandleFunc("GET /api/alerts", s.listAlerts)
	mux.HandleFunc("GET /api/alerts/{id}", s.getAlert)
	mux.HandleFunc("PUT /api/alerts/{id}", s.updateAlert)
	mux.HandleFunc("DELETE /api/alerts/{id}", s.deleteAlert)
	mux.HandleFunc("POST /api/alerts/{id}/transition", s.transitionAlertStatus)
}

type createAlertRequest struct {
	RuleID   string `json:"rule_id"`
	SourceID string `json:"source_id"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}

func (s *Server) createAlert(w http.ResponseWriter, r *http.Request) {
	var req createAlertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	alert, err := s.svc.CreateAlert(model.Alert{RuleID: req.RuleID, SourceID: req.SourceID, Level: req.Level, Message: req.Message})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, alert)
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertFilter{
		RuleID:   r.URL.Query().Get("rule_id"),
		SourceID: r.URL.Query().Get("source_id"),
		Level:    r.URL.Query().Get("level"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAlerts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	alert, err := s.svc.GetAlert(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, alert)
}

type updateAlertRequest struct {
	Message string `json:"message"`
	Level   string `json:"level"`
}

func (s *Server) updateAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	alert, err := s.svc.UpdateAlert(id, model.Alert{Message: req.Message, Level: req.Level})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, alert)
}

func (s *Server) deleteAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlert(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionAlertStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionAlertStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionAlertStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	alert, err := s.svc.TransitionAlertStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, alert)
}
