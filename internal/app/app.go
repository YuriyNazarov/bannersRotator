package app

import (
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
		return a.randomBanner(slotId)
	}
	bannerId, err := selectBanner(stats)
	if err != nil {
		return a.randomBanner(slotId)
	}
	return bannerId, nil
}

func (a *App) randomBanner(slotId int) (int, error) {
	banners, err := a.storage.GetAllBanners(slotId)
	if err != nil {
		return 0, err
	}
	return banners[rand.Intn(len(banners))], nil
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
