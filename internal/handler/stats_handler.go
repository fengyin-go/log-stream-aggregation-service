package handler

import (
	"net/http"
	"strconv"

	"log-aggregation/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/global", s.getGlobalStats)
	mux.HandleFunc("GET /api/stats/levels", s.getLevelStats)
	mux.HandleFunc("GET /api/stats/sources", s.getSourceStats)
	mux.HandleFunc("GET /api/stats/trend", s.getTimeTrend)
	mux.HandleFunc("GET /api/stats/alert-status", s.getAlertStatusStats)
	mux.HandleFunc("GET /api/stats/health", s.getSourceHealthStats)
	mux.HandleFunc("GET /api/sources/{id}/export", s.exportLogEntries)
	mux.HandleFunc("GET /api/sources/{id}/analysis", s.analyzeLogSource)
	mux.HandleFunc("GET /api/sources/{id}/error-rate", s.getErrorRateTrend)
}

func (s *Server) getGlobalStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetGlobalStats()
	httpx.OK(w, stats)
}

func (s *Server) getLevelStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetLevelStats()
	httpx.OK(w, stats)
}

func (s *Server) getSourceStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetSourceStats()
	httpx.OK(w, stats)
}

func (s *Server) getTimeTrend(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source_id")
	if sourceID == "" {
		httpx.BadRequest(w, "source_id 不能为空")
		return
	}
	hours, _ := strconv.Atoi(r.URL.Query().Get("hours"))
	stats := s.svc.GetTimeTrend(sourceID, hours)
	httpx.OK(w, stats)
}

func (s *Server) getAlertStatusStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetAlertStatusStats()
	httpx.OK(w, stats)
}

func (s *Server) getSourceHealthStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.ListSourceHealth()
	httpx.OK(w, stats)
}

func (s *Server) exportLogEntries(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := s.svc.ExportLogEntriesSummary(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, data)
}

func (s *Server) analyzeLogSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hours, _ := strconv.Atoi(r.URL.Query().Get("hours"))
	analysis, err := s.svc.AnalyzeLogSource(id, hours)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, analysis)
}

func (s *Server) getErrorRateTrend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hours, _ := strconv.Atoi(r.URL.Query().Get("hours"))
	stats := s.svc.GetErrorRateTrend(id, hours)
	httpx.OK(w, stats)
}
