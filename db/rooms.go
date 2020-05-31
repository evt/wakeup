package db

import (
	"github.com/evt/wakeup/model"
	"github.com/go-pg/pg/v9"
)

// AddRooms saves a list of rooms in Postgres
func (db *PgDB) AddRooms(rooms []*model.Room) error {
	_, err := db.Model(&rooms).
		OnConflict("(room_number) DO UPDATE").
		Insert()
	return err
}

// FindRooms returns room list by call time
func (db *PgDB) FindRooms(callTime string) ([]*model.Room, error) {
	var rooms []*model.Room
	err := db.Model(&rooms).Where("call_time = ?", callTime).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil
}

// IncRoomRetryCount increases call try count by 1
func (db *PgDB) IncRoomRetryCount(room *model.Room) error {
	room.RetryCount++
	_, err := db.Model(room).Set("retry_count = ?retry_count").Where("room_number = ?room_number").Update()
	return err
}
