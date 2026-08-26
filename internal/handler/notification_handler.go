package handler

import (
	"net/http"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/httpx"
)

func (s *Server) registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/notifications", s.createNotification)
	mux.HandleFunc("GET /api/notifications", s.listNotifications)
	mux.HandleFunc("GET /api/notifications/{id}", s.getNotification)
	mux.HandleFunc("POST /api/notifications/{id}/read", s.markNotificationAsRead)
	mux.HandleFunc("POST /api/notifications/read-all", s.markAllNotificationsAsRead)
	mux.HandleFunc("DELETE /api/notifications/{id}", s.deleteNotification)
	mux.HandleFunc("GET /api/notifications/unread-count", s.getUnreadCount)
}

type createNotificationRequest struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (s *Server) createNotification(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNotification(model.Notification{Type: req.Type, Title: req.Title, Message: req.Message})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNotifications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NotificationFilter{
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListNotifications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNotification(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) markNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.MarkNotificationAsRead(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) markAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	count, err := s.svc.MarkAllNotificationsAsRead()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"marked_count": count})
}

func (s *Server) deleteNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNotification(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getUnreadCount(w http.ResponseWriter, r *http.Request) {
	count := s.svc.GetUnreadNotificationCount()
	httpx.OK(w, map[string]interface{}{"unread_count": count})
}
