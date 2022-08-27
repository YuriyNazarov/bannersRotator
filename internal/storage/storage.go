package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/YuriyNazarov/bannersRotator/internal/app"
	_ "github.com/lib/pq" // Postgres driver
)

type Storage struct {
	db  *sql.DB
	log logger
}

const (
	view = iota
	click
)

func (s *Storage) AddToSlot(bannerID, slotID int) error {
	query := "select banner_id from banners_to_slots where banner_id = $1 and slot_id = $2"
	row := s.db.QueryRow(query, bannerID, slotID)
	var id string
	err := row.Scan(id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrLinkExists
	}
	query = "insert into banners_to_slots (banner_id, slot_id) values ($1, $2)"
	_, err = s.db.Exec(query, bannerID, slotID)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on creating banner_to_slot: %s", err))
		return ErrOperationFail
	}
	return nil
}

func (s *Storage) DropFromSlot(bannerID, slotID int) error {
	query := "delete from banners_to_slots where banner_id = $1 and slot_id = $2"
	_, err := s.db.Exec(query, bannerID, slotID)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on deleting banner_to_slot: %s", err))
		return ErrOperationFail
	}
	return nil
}

func (s *Storage) Click(bannerID, slotID, groupID int) error {
	return s.addAction(bannerID, slotID, groupID, click)
}

func (s *Storage) Show(bannerID, slotID, groupID int) error {
	return s.addAction(bannerID, slotID, groupID, view)
}

func (s *Storage) addAction(bannerID, slotID, groupID, actionID int) error {
	query := "insert into actions (action_type, banner_id, slot_id, dem_group_id) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, actionID, bannerID, slotID, groupID)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on saving action: %s", err))
		return ErrOperationFail
	}
	return nil
}

func (s *Storage) GetStats(slotID, groupID int) ([]app.BannerStat, error) {
	query := `select banners_to_slots.banner_id, coalesce(cnt, 0) as cnt, coalesce(action_type, 0) as action_type
		from banners_to_slots
		left join (
		select count(action_type) as cnt, action_type, banner_id
			from actions
			where slot_id = $1
			and dem_group_id = $2
			group by action_type, banner_id
			) stats on stats.banner_id = banners_to_slots.banner_id
			where slot_id = $1;`
	rows, err := s.db.Query(query, slotID, groupID)
	var (
		bannerID, count, actionID int
		bannerStat                app.BannerStat
		ok                        bool
	)
	statMap := make(map[int]app.BannerStat)

	if err != nil {
		s.log.Error(fmt.Sprintf("err on getting stats: %s", err))
		return []app.BannerStat{}, ErrOperationFail
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&bannerID, &count, &actionID)
		if err != nil {
			s.log.Error(fmt.Sprintf("err on scanning stats: %s", err))
			continue
		}
		bannerStat, ok = statMap[bannerID]
		if !ok {
			bannerStat = app.BannerStat{BannerID: bannerID}
		}
		if actionID == view {
			bannerStat.Views = count
		} else {
			bannerStat.Clicks = count
		}
		statMap[bannerID] = bannerStat
	}
	if rows.Err() != nil {
		s.log.Error(fmt.Sprintf("err after scanning stats: %s", err))
	}
	stats := make([]app.BannerStat, len(statMap))
	i := 0
	for _, v := range statMap {
		stats[i] = v
		i++
	}

	return stats, nil
}

func (s *Storage) GetAllBanners(slotID int) ([]int, error) {
	query := "select banner_id from banners_to_slots where slot_id = $1"
	rows, err := s.db.Query(query, slotID)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on getting banners: %s", err))
		return []int{}, ErrOperationFail
	}
	defer rows.Close()

	var (
		banners []int
		banner  int
	)
	for rows.Next() {
		err = rows.Scan(&banner)
		if err != nil {
			s.log.Error(fmt.Sprintf("err on scanning banners: %s", err))
			continue
		}
		banners = append(banners, banner)
	}
	if rows.Err() != nil {
		s.log.Error(fmt.Sprintf("err on scanning banners: %s", err))
	}
	if len(banners) == 0 {
		return banners, ErrEmptyResult
	}
	return banners, nil
}

func New(l logger, dsn string) *Storage {
	storageInstance := Storage{
		log: l,
	}
	err := storageInstance.Connect(dsn)
	if err != nil {
		l.Error(fmt.Sprintf("error occupied on connecting to DB: %s", err))
		return nil
	}
	return &storageInstance
}

func (s *Storage) Connect(dsn string) error {
	s.log.Debug("connect DB")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed on connecting to db: %s", err))
		return ErrConnFailed
	}
	err = db.Ping()
	if err != nil {
		s.log.Error(fmt.Sprintf("db healthcheck failed: %s", err))
		return ErrConnFailed
	}
	s.db = db
	return nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.log.Error(fmt.Sprintf("err on clolsing db connection: %s", err))
	}
}
