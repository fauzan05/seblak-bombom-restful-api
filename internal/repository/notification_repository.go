package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type NotificationRepository struct {
	Repository[entity.Notification]
	Log *logrus.Logger
}

func NewNotificationRepository(log *logrus.Logger) *NotificationRepository {
	return &NotificationRepository{
		Log: log,
	}
}
