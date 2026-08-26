package handler

import (
	"net/http"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sources/{id}/export/csv", s.exportLogEntriesCSV)
	mux.HandleFunc("GET /api/alerts/export/csv", s.exportAlertsCSV)
	mux.HandleFunc("GET /api/sources/export/json", s.exportSourcesJSON)
}

func (s *Server) exportLogEntriesCSV(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.ExportLogEntriesCSV(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+result.Filename)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result.Body))
}

func (s *Server) exportAlertsCSV(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source_id")
	result, err := s.svc.ExportAlertsCSV(sourceID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+result.Filename)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result.Body))
}

func (s *Server) exportSourcesJSON(w http.ResponseWriter, r *http.Request) {
	body, err := s.svc.ExportSourcesJSON()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=sources.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}
