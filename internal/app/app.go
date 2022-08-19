package app

import (
	"fmt"
	"github.com/YuriyNazarov/bannersRotator/internal/amqp"
	"github.com/YuriyNazarov/bannersRotator/internal/storage"
	"math/rand"
	"time"
)

type App struct {
	logger    Logger
	storage   storage.BannersRepository
	Analytics amqp.StatsOutput
}

func (a *App) GetBanner(slotId, groupId int) (int, error) {
	stats, err := a.storage.GetStats(slotId, groupId)
	if err != nil {
		return a.randomBanner(slotId, groupId)
	}
	bannerId, err := selectBanner(stats)
	if err != nil {
		return a.randomBanner(slotId, groupId)
	}

	err = a.storage.Show(bannerId, slotId, groupId)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to save action: %s", err))
	}
	a.Analytics.Show(bannerId, slotId, groupId, time.Now())
	return bannerId, nil
}

func (a *App) randomBanner(slotId, groupId int) (int, error) {
	banners, err := a.storage.GetAllBanners(slotId)
	if err != nil {
		return 0, err
	}
	bannerId := banners[rand.Intn(len(banners))]

	err = a.storage.Show(bannerId, slotId, groupId)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to save action: %s", err))
	}
	a.Analytics.Show(bannerId, slotId, groupId, time.Now())
	return bannerId, nil
}

func (a *App) AddBanner(bannerId, slotId int) error {
	return a.storage.AddToSlot(bannerId, slotId)
}

func (a *App) DeleteBanner(bannerId, slotId int) error {
	return a.storage.DropFromSlot(bannerId, slotId)
}

func (a *App) RegisterClick(bannerId, slotId, groupId int) error {
	a.Analytics.Click(bannerId, slotId, groupId, time.Now())
	return a.storage.Click(bannerId, slotId, groupId)
}

func New(l Logger, s storage.BannersRepository, q amqp.StatsOutput) *App {
	return &App{
		logger:    l,
		storage:   s,
		Analytics: q,
	}
}
