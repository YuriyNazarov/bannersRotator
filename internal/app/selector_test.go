package app

import (
	"github.com/YuriyNazarov/bannersRotator/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBestBanner(t *testing.T) {
	bannerStats := []storage.BannerStat{
		{BannerId: 1, Views: 10, Clicks: 0},
		{BannerId: 2, Views: 10, Clicks: 5},
		{BannerId: 3, Views: 10, Clicks: 10},
		{BannerId: 4, Views: 10, Clicks: 0},
	}

	expectedID := 3
	bannerID, err := selectBanner(bannerStats)
	require.NoError(t, err)
	require.Equal(t, expectedID, bannerID)
}

func TestShowAllBanners(t *testing.T) {
	bannerStats := []storage.BannerStat{
		{BannerId: 1, Views: 99999, Clicks: 1},
		{BannerId: 2, Views: 99999, Clicks: 1},
		{BannerId: 3, Views: 99999, Clicks: 1},
		{BannerId: 4, Views: 99999, Clicks: 1},
	}

	var (
		bannerID int
		banners  []int
		err      error
	)
	for i := 0; i < 100; i++ {
		bannerID, err = selectBanner(bannerStats)
		require.NoError(t, err)
		banners = append(banners, bannerID)
	}

	require.Contains(t, banners, 1)
	require.Contains(t, banners, 2)
	require.Contains(t, banners, 3)
	require.Contains(t, banners, 4)
}

func TestNoBanners(t *testing.T) {
	_, err := selectBanner([]storage.BannerStat{})
	require.Equal(t, ErrNoBanners, err)
}
