package storage

import "errors"

type BannerStat struct {
	BannerId int
	Views    int
	Clicks   int
}

type BannersRepository interface {
	AddToSlot(bannerId, slotId int) error
	DropFromSlot(bannerId, slotId int) error
	Click(bannerId, slotId, groupId int) error
	Show(bannerId, slotId, groupId int) error
	GetStats(slotId, groupId int) ([]BannerStat, error)
	GetAllBanners(slotId int) ([]int, error)
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

var (
	ErrConnFailed    = errors.New("failed to connect to database")
	ErrLinkExists    = errors.New("relation already exists")
	ErrOperationFail = errors.New("whoops, something went wrong during database operation")
	ErrEmptyResult   = errors.New("whoops, something went wrong, we have no idea what to show")
)
