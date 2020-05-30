package model

import "time"

// Room is a hotel room to call
type Room struct {
	tableName  struct{}  `pg:"rooms"`
	RoomID     string    `json:"room_id" pg:"room_id,notnull,unique"`
	Firstname  string    `json:"firstname" pg:"firstname,notnull"`
	Lastname   string    `json:"lastname" pg:"lastname,notnull"`
	RoomNumber int       `json:"room_number" pg:"room_number,notnull"`
	CallTime   string    `json:"call_time" pg:"call_time,notnull"`
	Created    time.Time `json:"created" pg:"created"`
}
