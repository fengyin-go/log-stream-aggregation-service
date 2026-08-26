package handler

import (
	"net/http"
	"strconv"

	"log-aggregation/pkg/httpx"
)

func (s *Server) registerAdvancedStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/top-sources", s.getTopSources)
	mux.HandleFunc("GET /api/stats/recent-errors", s.getRecentErrors)
	mux.HandleFunc("GET /api/stats/level-distribution", s.getLevelDistribution)
	mux.HandleFunc("GET /api/stats/hourly-peak", s.getHourlyPeak)
}

func (s *Server) getTopSources(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	stats := s.svc.GetTopSources(limit)
	httpx.OK(w, stats)
}

func (s *Server) getRecentErrors(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	stats := s.svc.GetRecentErrors(limit)
	httpx.OK(w, stats)
}

func (s *Server) getLevelDistribution(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source_id")
	stats := s.svc.GetLevelDistribution(sourceID)
	httpx.OK(w, stats)
}

func (s *Server) getHourlyPeak(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source_id")
	hours, _ := strconv.Atoi(r.URL.Query().Get("hours"))
	stats := s.svc.GetHourlyPeak(sourceID, hours)
	httpx.OK(w, stats)
}
