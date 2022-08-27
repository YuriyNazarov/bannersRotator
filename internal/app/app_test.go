package app_test

// Different pkg because mock imports struct from app

import (
	"errors"
	"testing"

	internalapp "github.com/YuriyNazarov/bannersRotator/internal/app"
	"github.com/YuriyNazarov/bannersRotator/internal/app/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetBanner(t *testing.T) {
	t.Run("Successful banner selection", func(t *testing.T) {
		logger := mocks.Logger{}
		bannersRepo := mocks.BannersRepository{}
		statsRepo := mocks.StatsRepository{}
		selector := mocks.BannersSelector{}
		queue := mocks.StatsOutput{}
		app := internalapp.New(&logger, &bannersRepo, &statsRepo, &selector, &queue)
		stats := []internalapp.BannerStat{
			{BannerID: 1, Views: 1000, Clicks: 0},
			{BannerID: 2, Views: 1000, Clicks: 0},
			{BannerID: 3, Views: 1000, Clicks: 0},
			{BannerID: 4, Views: 1000, Clicks: 0},
		}
		statsRepo.On("GetStats", 1, 1).Return(stats, nil)
		selector.On("SelectBanner", stats).Return(1, nil)
		statsRepo.On("Show", 1, 1, 1).Return(nil)
		queue.On(
			"View",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return()
		logger.On("Error", mock.Anything).Return()
		id, err := app.GetBanner(1, 1)
		require.NoError(t, err)
		require.Contains(t, []int{1, 2, 3, 4}, id)
		statsRepo.AssertCalled(t, "GetStats", 1, 1)
		selector.AssertCalled(t, "SelectBanner", stats)
		statsRepo.AssertCalled(t, "Show", 1, 1, 1)
		queue.AssertCalled(t, "View", 1, 1, 1, mock.Anything)
		logger.AssertNotCalled(t, "Error", mock.Anything)
	})

	t.Run("Random banner selection", func(t *testing.T) {
		logger := mocks.Logger{}
		bannersRepo := mocks.BannersRepository{}
		statsRepo := mocks.StatsRepository{}
		selector := mocks.BannersSelector{}
		queue := mocks.StatsOutput{}
		app := internalapp.New(&logger, &bannersRepo, &statsRepo, &selector, &queue)
		statsRepo.On("GetStats", 1, 1).Return([]internalapp.BannerStat{}, errors.New("failed"))
		bannersRepo.On("GetAllBanners", 1).Return([]int{2}, nil)
		statsRepo.On("Show", 2, 1, 1).Return(nil)
		queue.On(
			"View",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return()
		id, err := app.GetBanner(1, 1)
		require.NoError(t, err)
		require.Equal(t, 2, id)
		bannersRepo.AssertCalled(t, "GetAllBanners", 1)
	})

	t.Run("Error on banner selection", func(t *testing.T) {
		logger := mocks.Logger{}
		bannersRepo := mocks.BannersRepository{}
		statsRepo := mocks.StatsRepository{}
		selector := mocks.BannersSelector{}
		queue := mocks.StatsOutput{}
		app := internalapp.New(&logger, &bannersRepo, &statsRepo, &selector, &queue)
		statsRepo.On("GetStats", 1, 1).Return([]internalapp.BannerStat{}, errors.New("failed"))
		bannersRepo.On("GetAllBanners", 1).Return([]int{}, errors.New("failed"))
		_, err := app.GetBanner(1, 1)
		require.NotEmpty(t, err)
		bannersRepo.AssertCalled(t, "GetAllBanners", 1)
	})
}

func TestAddBanner(t *testing.T) {
	logger := mocks.Logger{}
	bannersRepo := mocks.BannersRepository{}
	statsRepo := mocks.StatsRepository{}
	selector := mocks.BannersSelector{}
	queue := mocks.StatsOutput{}
	app := internalapp.New(&logger, &bannersRepo, &statsRepo, &selector, &queue)
	bannersRepo.On("AddToSlot", 1, 1).Return(nil)
	_ = app.AddBanner(1, 1)
	bannersRepo.AssertCalled(t, "AddToSlot", 1, 1)
}

func TestDeleteBanner(t *testing.T) {
	logger := mocks.Logger{}
	bannersRepo := mocks.BannersRepository{}
	statsRepo := mocks.StatsRepository{}
	selector := mocks.BannersSelector{}
	queue := mocks.StatsOutput{}
	app := internalapp.New(&logger, &bannersRepo, &statsRepo, &selector, &queue)
	bannersRepo.On("DropFromSlot", 1, 1).Return(nil)
	_ = app.DeleteBanner(1, 1)
	bannersRepo.AssertCalled(t, "DropFromSlot", 1, 1)
}

func TestRegisterClick(t *testing.T) {
	logger := mocks.Logger{}
	bannersRepo := mocks.BannersRepository{}
	statsRepo := mocks.StatsRepository{}
	selector := mocks.BannersSelector{}
	queue := mocks.StatsOutput{}
	app := internalapp.New(&logger, &bannersRepo, &statsRepo, &selector, &queue)
	queue.On(
		"Click",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.Anything,
	).Return()
	statsRepo.On("Click", 1, 1, 1).Return(nil)
	_ = app.RegisterClick(1, 1, 1)
	queue.AssertCalled(t, "Click", 1, 1, 1, mock.Anything)
	statsRepo.AssertCalled(t, "Click", 1, 1, 1)
}
