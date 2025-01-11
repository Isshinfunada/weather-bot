package entity

import "time"

type NotificationHistory struct {
	ID               int
	UserID           int
	NotificationTime time.Time
	IsNotifyTrigger  bool
	WeatherCodes     []string
	WeatherData      []byte
	CreatedAt        time.Time
}
