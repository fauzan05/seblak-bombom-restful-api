package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type NotificationRepository struct {
	Repository[entity.Category]
	Log *logrus.Logger
}

func NewNotificationRepository(log *logrus.Logger) *NotificationRepository {
	return &NotificationRepository{
		Log: log,
	}
}
