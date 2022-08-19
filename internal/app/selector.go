package app

import (
	"math"
	"math/rand"

	"github.com/YuriyNazarov/bannersRotator/internal/storage"
)

func selectBanner(stats []storage.BannerStat) (int, error) {
	if len(stats) == 0 {
		return 0, ErrNoBanners
	}

	var (
		bannerIds  []int
		totalViews int
		maxIncome  float64 = -1
	)
	for _, bannerStat := range stats {
		totalViews += bannerStat.Views
	}
	// создать список не показанных
	for _, bannerStat := range stats {
		if bannerStat.Views == 0 {
			bannerIds = append(bannerIds, bannerStat.BannerID)
		}
	}
	// если есть новый баннер - покажем его. инициализация
	if len(bannerIds) > 0 {
		// случайный баннер из новых
		return bannerIds[rand.Intn(len(bannerIds))], nil //nolint:gosec
	}

	// новых нет. создать список баннеров с макс доходом
	for _, bannerStat := range stats {
		bannerIncome := (float64(bannerStat.Clicks) / float64(bannerStat.Views)) +
			math.Sqrt((2.0*math.Log(float64(totalViews)))/float64(bannerStat.Views))
		if math.IsNaN(bannerIncome) {
			bannerIncome = -1
		}

		if bannerIncome > maxIncome {
			maxIncome = bannerIncome
			bannerIds = bannerIds[:0]
			bannerIds = append(bannerIds, bannerStat.BannerID)
		} else if bannerIncome == maxIncome {
			bannerIds = append(bannerIds, bannerStat.BannerID)
		}
	}

	// случайный баннер из лучших
	return bannerIds[rand.Intn(len(bannerIds))], nil //nolint:gosec
}
