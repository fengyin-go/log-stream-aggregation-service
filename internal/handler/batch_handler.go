package handler

import (
	"net/http"

	"log-aggregation/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alerts/batch/acknowledge", s.batchAcknowledgeAlerts)
	mux.HandleFunc("POST /api/alerts/batch/resolve", s.batchResolveAlerts)
	mux.HandleFunc("POST /api/entries/batch/delete", s.batchDeleteLogEntries)
	mux.HandleFunc("POST /api/sources/batch/pause", s.batchPauseSources)
	mux.HandleFunc("POST /api/sources/batch/resume", s.batchResumeSources)
	mux.HandleFunc("POST /api/entries/cleanup", s.cleanupOldEntries)
	mux.HandleFunc("POST /api/alert-rules/{id}/duplicate", s.duplicateAlertRule)
}

type batchAlertIDsRequest struct {
	AlertIDs []string `json:"alert_ids"`
}

func (s *Server) batchAcknowledgeAlerts(w http.ResponseWriter, r *http.Request) {
	var req batchAlertIDsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchAcknowledgeAlerts(req.AlertIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"acknowledged_count": count})
}

func (s *Server) batchResolveAlerts(w http.ResponseWriter, r *http.Request) {
	var req batchAlertIDsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchResolveAlerts(req.AlertIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"resolved_count": count})
}

type batchEntryIDsRequest struct {
	EntryIDs []string `json:"entry_ids"`
}

func (s *Server) batchDeleteLogEntries(w http.ResponseWriter, r *http.Request) {
	var req batchEntryIDsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchDeleteLogEntries(req.EntryIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"deleted_count": count})
}

type batchSourceIDsRequest struct {
	SourceIDs []string `json:"source_ids"`
}

func (s *Server) batchPauseSources(w http.ResponseWriter, r *http.Request) {
	var req batchSourceIDsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchPauseSources(req.SourceIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"paused_count": count})
}

func (s *Server) batchResumeSources(w http.ResponseWriter, r *http.Request) {
	var req batchSourceIDsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchResumeSources(req.SourceIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"resumed_count": count})
}

type cleanupRequest struct {
	Hours int `json:"hours"`
}

func (s *Server) cleanupOldEntries(w http.ResponseWriter, r *http.Request) {
	var req cleanupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.CleanupOldEntries(req.Hours)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"cleaned_count": count})
}

func (s *Server) duplicateAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	newRule, err := s.svc.DuplicateAlertRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, newRule)
}
