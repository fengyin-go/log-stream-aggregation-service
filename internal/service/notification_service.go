package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateNotification(input model.Notification) (*model.Notification, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	n := &model.Notification{
		ID:        idgen.Hex(),
		Type:      input.Type,
		Title:     input.Title,
		Message:   input.Message,
		Status:    model.NotificationStatusUnread,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) GetNotification(id string) (*model.Notification, error) {
	return s.store.GetNotification(id)
}

func (s *Service) ListNotifications(filter model.NotificationFilter, page, size int) ([]*model.Notification, int, error) {
	all := s.store.ListNotifications()
	matched := make([]*model.Notification, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Notification{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) MarkNotificationAsRead(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.Status = model.NotificationStatusRead
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) MarkAllNotificationsAsRead() (int, error) {
	all := s.store.ListNotifications()
	count := 0
	for _, n := range all {
		if n.Status == model.NotificationStatusUnread {
			n.Status = model.NotificationStatusRead
			if err := s.store.UpdateNotification(n); err == nil {
				count++
			}
		}
	}
	return count, nil
}

func (s *Service) DeleteNotification(id string) error {
	if _, err := s.store.GetNotification(id); err != nil {
		return err
	}
	return s.store.DeleteNotification(id)
}

func (s *Service) GetUnreadNotificationCount() int {
	return s.store.CountUnreadNotifications()
}
