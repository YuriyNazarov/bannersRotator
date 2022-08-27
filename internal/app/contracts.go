package app

import (
	"errors"
	"time"
)

type Logger interface {
	Close()
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

type StatsOutput interface {
	Click(bannerID, slotID, groupID int, clickTime time.Time)
	View(bannerID, slotID, groupID int, clickTime time.Time)
}

type BannerStat struct {
	BannerID int
	Views    int
	Clicks   int
}

type BannersRepository interface {
	AddToSlot(bannerID, slotID int) error
	DropFromSlot(bannerID, slotID int) error
	GetAllBanners(slotID int) ([]int, error)
}

type StatsRepository interface {
	Click(bannerID, slotID, groupID int) error
	Show(bannerID, slotID, groupID int) error
	GetStats(slotID, groupID int) ([]BannerStat, error)
}

type BannersSelector interface {
	SelectBanner(stats []BannerStat) (int, error)
}

var ErrNoBanners = errors.New("no banners available for selected slot")
