package amqp

import "time"

type StatsOutput interface {
	Click(bannerId, slotId, groupId int, clickTime time.Time)
	Show(bannerId, slotId, groupId int, clickTime time.Time)
}
