package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerLogEntryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/entries", s.createLogEntry)
	mux.HandleFunc("POST /api/entries/batch", s.batchCreateLogEntries)
	mux.HandleFunc("GET /api/entries", s.listLogEntries)
	mux.HandleFunc("GET /api/entries/{id}", s.getLogEntry)
	mux.HandleFunc("DELETE /api/entries/source/{source_id}", s.deleteLogEntriesBySource)
}

type createLogEntryRequest struct {
	SourceID  string   `json:"source_id"`
	Level     string   `json:"level"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	Tags      []string `json:"tags"`
}

func (s *Server) createLogEntry(w http.ResponseWriter, r *http.Request) {
	var req createLogEntryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	entry, err := s.svc.CreateLogEntry(model.LogEntry{SourceID: req.SourceID, Level: req.Level, Message: req.Message, Tags: req.Tags})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, entry)
}

type batchCreateLogEntriesRequest struct {
	Entries []createLogEntryRequest `json:"entries"`
}

func (s *Server) batchCreateLogEntries(w http.ResponseWriter, r *http.Request) {
	var req batchCreateLogEntriesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.LogEntry, 0, len(req.Entries))
	for _, e := range req.Entries {
		inputs = append(inputs, model.LogEntry{SourceID: e.SourceID, Level: e.Level, Message: e.Message, Tags: e.Tags})
	}
	entries, err := s.svc.BatchCreateLogEntries(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, entries)
}

func (s *Server) listLogEntries(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.LogEntryFilter{
		SourceID: r.URL.Query().Get("source_id"),
		Level:    r.URL.Query().Get("level"),
		Keyword:  r.URL.Query().Get("keyword"),
		Tag:      r.URL.Query().Get("tag"),
	}
	items, total, err := s.svc.ListLogEntries(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getLogEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entry, err := s.svc.GetLogEntry(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, entry)
}

func (s *Server) deleteLogEntriesBySource(w http.ResponseWriter, r *http.Request) {
	sourceID := r.PathValue("source_id")
	if err := s.svc.DeleteLogEntriesBySourceID(sourceID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
