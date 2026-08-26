package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerRetentionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/retention-policies", s.createRetentionPolicy)
	mux.HandleFunc("GET /api/retention-policies", s.listRetentionPolicies)
	mux.HandleFunc("GET /api/retention-policies/{id}", s.getRetentionPolicy)
	mux.HandleFunc("GET /api/sources/{id}/retention", s.getRetentionPolicyBySource)
	mux.HandleFunc("PUT /api/retention-policies/{id}", s.updateRetentionPolicy)
	mux.HandleFunc("DELETE /api/retention-policies/{id}", s.deleteRetentionPolicy)
	mux.HandleFunc("POST /api/sources/{id}/retention/apply", s.applyRetentionPolicy)
}

type createRetentionPolicyRequest struct {
	SourceID string `json:"source_id"`
	Policy   string `json:"policy"`
}

func (s *Server) createRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	var req createRetentionPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateRetentionPolicy(model.RetentionPolicy{SourceID: req.SourceID, Policy: req.Policy})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listRetentionPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RetentionPolicyFilter{
		SourceID: r.URL.Query().Get("source_id"),
		Policy:   r.URL.Query().Get("policy"),
	}
	items, total, err := s.svc.ListRetentionPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetRetentionPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) getRetentionPolicyBySource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetRetentionPolicyBySourceID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type updateRetentionPolicyRequest struct {
	Policy string `json:"policy"`
}

func (s *Server) updateRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRetentionPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateRetentionPolicy(id, model.RetentionPolicy{Policy: req.Policy})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRetentionPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) applyRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	count, err := s.svc.ApplyRetentionPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"cleaned_count": count})
}
