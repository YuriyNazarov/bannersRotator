package amqp

import "time"

type StatsOutput interface {
	Click(bannerID, slotID, groupID int, clickTime time.Time)
	Show(bannerID, slotID, groupID int, clickTime time.Time)
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

type statsMessage struct {
	BannerID   int       `json:"bannerId"`
	SlotID     int       `json:"slotId"`
	GroupID    int       `json:"groupId"`
	Timestamp  time.Time `json:"timestamp"`
	ActionType string    `json:"actionType"`
}
