package model

import "time"

// Room is a hotel room to call
type Room struct {
	tableName  struct{}  `pg:"rooms"`
	RoomNumber int       `json:"room_number" pg:"room_number,notnull,pk"`
	Firstname  string    `json:"firstname" pg:"firstname,notnull"`
	Lastname   string    `json:"lastname" pg:"lastname,notnull"`
	CallTime   string    `json:"call_time" pg:"call_time,notnull"`
	RetryCount int       `json:"retry_count" pg:"retry_count,notnull"`
	Created    time.Time `json:"created" pg:"created"`
}
