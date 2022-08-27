package app

import (
	"fmt"
	"math/rand"
	"time"
)

type App struct {
	logger    Logger
	storage   BannersRepository
	stats     StatsRepository
	selector  BannersSelector
	Analytics StatsOutput
}

func (a *App) GetBanner(slotID, groupID int) (int, error) {
	stats, err := a.stats.GetStats(slotID, groupID)
	if err != nil {
		return a.randomBanner(slotID, groupID)
	}
	bannerID, err := a.selector.SelectBanner(stats)
	if err != nil {
		return a.randomBanner(slotID, groupID)
	}

	err = a.stats.Show(bannerID, slotID, groupID)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to save action: %s", err))
	}
	a.Analytics.View(bannerID, slotID, groupID, time.Now())
	return bannerID, nil
}

func (a *App) randomBanner(slotID, groupID int) (int, error) {
	banners, err := a.storage.GetAllBanners(slotID)
	if err != nil {
		return 0, err
	}
	bannerID := banners[rand.Intn(len(banners))] //nolint:gosec

	err = a.stats.Show(bannerID, slotID, groupID)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to save action: %s", err))
	}
	a.Analytics.View(bannerID, slotID, groupID, time.Now())
	return bannerID, nil
}

func (a *App) AddBanner(bannerID, slotID int) error {
	return a.storage.AddToSlot(bannerID, slotID)
}

func (a *App) DeleteBanner(bannerID, slotID int) error {
	return a.storage.DropFromSlot(bannerID, slotID)
}

func (a *App) RegisterClick(bannerID, slotID, groupID int) error {
	a.Analytics.Click(bannerID, slotID, groupID, time.Now())
	return a.stats.Click(bannerID, slotID, groupID)
}

func New(l Logger, storage BannersRepository, stats StatsRepository, selector BannersSelector, q StatsOutput) *App {
	return &App{
		logger:    l,
		storage:   storage,
		stats:     stats,
		selector:  selector,
		Analytics: q,
	}
}
