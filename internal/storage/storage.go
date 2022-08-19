package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
)

type Storage struct {
	db  *sql.DB
	log Logger
}

const (
	view = iota
	click
)

func (s *Storage) AddToSlot(bannerId, slotId int) error {
	query := "select banner_id from banners_to_slots where banner_id = $1 and slot_id = $2"
	row := s.db.QueryRow(query, bannerId, slotId)
	var id string
	err := row.Scan(id)
	if err != sql.ErrNoRows {
		return ErrLinkExists
	}
	query = "insert into banners_to_slots (banner_id, slot_id) values ($1, $2)"
	_, err = s.db.Exec(query, bannerId, slotId)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on creating banner_to_slot: %s", err))
		return ErrOperationFail
	}
	return nil
}

func (s *Storage) DropFromSlot(bannerId, slotId int) error {
	query := "delete from banners_to_slots where banner_id = $1 and slot_id = $2"
	_, err := s.db.Exec(query, bannerId, slotId)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on deleting banner_to_slot: %s", err))
		return ErrOperationFail
	}
	return nil
}

func (s *Storage) Click(bannerId, slotId, groupId int) error {
	return s.addAction(bannerId, slotId, groupId, click)
}

func (s *Storage) Show(bannerId, slotId, groupId int) error {
	return s.addAction(bannerId, slotId, groupId, view)
}

func (s *Storage) addAction(bannerId, slotId, groupId, actionId int) error {
	query := "insert into actions (action_type, banner_id, slot_id, dem_group_id) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, actionId, bannerId, slotId, groupId)
	if err != nil {
		s.log.Error(fmt.Sprintf("err on saving action: %s", err))
		return ErrOperationFail
	}
	return nil
}

func (s *Storage) GetStats(slotId, groupId int) ([]BannerStat, error) {
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
	rows, err := s.db.Query(query, slotId, groupId)
	var (
		bannerId, count, actionId int
		bannerStat                BannerStat
		ok                        bool
	)
	statMap := make(map[int]BannerStat)

	if err != nil {
		s.log.Error(fmt.Sprintf("err on getting stats: %s", err))
		return []BannerStat{}, ErrOperationFail
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&bannerId, &count, &actionId)
		if err != nil {
			s.log.Error(fmt.Sprintf("err on scanning stats: %s", err))
			continue
		}
		bannerStat, ok = statMap[bannerId]
		if ok {
			if actionId == view {
				bannerStat.Views = count
			} else {
				bannerStat.Clicks = count
			}
		} else {
			bannerStat = BannerStat{BannerId: bannerId}
			if actionId == view {
				bannerStat.Views = count
			} else {
				bannerStat.Clicks = count
			}
		}
		statMap[bannerId] = bannerStat
	}
	if rows.Err() != nil {
		s.log.Error(fmt.Sprintf("err after scanning stats: %s", err))
	}
	stats := make([]BannerStat, len(statMap))
	i := 0
	for _, v := range statMap {
		stats[i] = v
		i++
	}

	return stats, nil
}

func (s *Storage) GetAllBanners(slotId int) ([]int, error) {
	query := "select banner_id from banners_to_slots where slot_id = $1"
	rows, err := s.db.Query(query, slotId)
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

func New(l Logger, dsn string) *Storage {
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
	err := s.db.Close()
	if err != nil {
		s.log.Error(fmt.Sprintf("err on clolsing db connection: %s", err))
	}
}
