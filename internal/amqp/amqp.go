package amqp

import "time"

type StatsOutput interface {
	Click(bannerId, slotId, groupId int, clickTime time.Time)
	Show(bannerId, slotId, groupId int, clickTime time.Time)
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

type statsMessage struct {
	BannerId   int       `json:"bannerId"`
	SlotId     int       `json:"slotId"`
	GroupId    int       `json:"groupId"`
	Timestamp  time.Time `json:"timestamp"`
	ActionType string    `json:"actionType"`
}
