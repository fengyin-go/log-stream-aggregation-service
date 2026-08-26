package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerQueryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/queries", s.createQuery)
	mux.HandleFunc("GET /api/queries", s.listQueries)
	mux.HandleFunc("GET /api/queries/{id}", s.getQuery)
	mux.HandleFunc("PUT /api/queries/{id}", s.updateQuery)
	mux.HandleFunc("DELETE /api/queries/{id}", s.deleteQuery)
}

type createQueryRequest struct {
	SourceID   string `json:"source_id"`
	Expression string `json:"expression"`
	CreatedBy  string `json:"created_by"`
}

func (s *Server) createQuery(w http.ResponseWriter, r *http.Request) {
	var req createQueryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	q, err := s.svc.CreateQuery(model.Query{SourceID: req.SourceID, Expression: req.Expression, CreatedBy: req.CreatedBy})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, q)
}

func (s *Server) listQueries(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.QueryFilter{
		SourceID:  r.URL.Query().Get("source_id"),
		CreatedBy: r.URL.Query().Get("created_by"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListQueries(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getQuery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q, err := s.svc.GetQuery(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, q)
}

type updateQueryRequest struct {
	SourceID   string `json:"source_id"`
	Expression string `json:"expression"`
	CreatedBy  string `json:"created_by"`
}

func (s *Server) updateQuery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateQueryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	q, err := s.svc.UpdateQuery(id, model.Query{SourceID: req.SourceID, Expression: req.Expression, CreatedBy: req.CreatedBy})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, q)
}

func (s *Server) deleteQuery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteQuery(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
