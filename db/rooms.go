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
