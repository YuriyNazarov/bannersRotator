package storage

import "errors"

type BannerStat struct {
	BannerID int
	Views    int
	Clicks   int
}

type BannersRepository interface {
	AddToSlot(bannerID, slotID int) error
	DropFromSlot(bannerID, slotID int) error
	Click(bannerID, slotID, groupID int) error
	Show(bannerID, slotID, groupID int) error
	GetStats(slotID, groupID int) ([]BannerStat, error)
	GetAllBanners(slotID int) ([]int, error)
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
