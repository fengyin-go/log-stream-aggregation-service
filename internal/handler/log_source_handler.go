package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerLogSourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sources", s.createLogSource)
	mux.HandleFunc("GET /api/sources", s.listLogSources)
	mux.HandleFunc("GET /api/sources/{id}", s.getLogSource)
	mux.HandleFunc("PUT /api/sources/{id}", s.updateLogSource)
	mux.HandleFunc("DELETE /api/sources/{id}", s.deleteLogSource)
	mux.HandleFunc("POST /api/sources/{id}/transition", s.transitionLogSourceStatus)
}

type createLogSourceRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
}

func (s *Server) createLogSource(w http.ResponseWriter, r *http.Request) {
	var req createLogSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	src, err := s.svc.CreateLogSource(model.LogSource{Name: req.Name, Type: req.Type, Host: req.Host, Path: req.Path})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, src)
}

func (s *Server) listLogSources(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.LogSourceFilter{
		Name:    r.URL.Query().Get("name"),
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListLogSources(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getLogSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	src, err := s.svc.GetLogSource(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, src)
}

type updateLogSourceRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
}

func (s *Server) updateLogSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateLogSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	src, err := s.svc.UpdateLogSource(id, model.LogSource{Name: req.Name, Type: req.Type, Host: req.Host, Path: req.Path})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, src)
}

func (s *Server) deleteLogSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteLogSource(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionLogSourceStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionLogSourceStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionLogSourceStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	src, err := s.svc.TransitionLogSourceStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, src)
}
