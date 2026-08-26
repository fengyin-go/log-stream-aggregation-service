package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerIndexRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/indexes", s.createIndex)
	mux.HandleFunc("GET /api/indexes", s.listIndexes)
	mux.HandleFunc("GET /api/indexes/{id}", s.getIndex)
	mux.HandleFunc("PUT /api/indexes/{id}", s.updateIndex)
	mux.HandleFunc("DELETE /api/indexes/{id}", s.deleteIndex)
}

type createIndexRequest struct {
	Name     string `json:"name"`
	SourceID string `json:"source_id"`
	Field    string `json:"field"`
}

func (s *Server) createIndex(w http.ResponseWriter, r *http.Request) {
	var req createIndexRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	idx, err := s.svc.CreateIndex(model.Index{Name: req.Name, SourceID: req.SourceID, Field: req.Field})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, idx)
}

func (s *Server) listIndexes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.IndexFilter{
		Name:     r.URL.Query().Get("name"),
		SourceID: r.URL.Query().Get("source_id"),
		Field:    r.URL.Query().Get("field"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListIndexes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getIndex(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idx, err := s.svc.GetIndex(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, idx)
}

type updateIndexRequest struct {
	Name  string `json:"name"`
	Field string `json:"field"`
}

func (s *Server) updateIndex(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateIndexRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	idx, err := s.svc.UpdateIndex(id, model.Index{Name: req.Name, Field: req.Field})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, idx)
}

func (s *Server) deleteIndex(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteIndex(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
