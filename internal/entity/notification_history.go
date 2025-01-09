package entity

import "time"

type NotificationHistory struct {
	ID               int
	UserID           int
	NotificationTime time.Time
	WeatherData      []byte
	IsNotifyTrigger  bool
	CreatedAt        time.Time
}
