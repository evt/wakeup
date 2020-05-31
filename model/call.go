package model

import "time"

// CallRoomRequest is a request coming from google cloud scheduler to call rooms
type CallRoomRequest struct {
	callTime string `json:"call_time"`
}

// Call is a room call that we save in Postgres at the moment room is called to know when it happened and what call status was
type Call struct {
	tableName  struct{}  `pg:"calls"`
	CallID     int       `json:"call_id" pg:"call_id,notnull,pk"`
	RoomNumber int       `json:"room_number" pg:"room_number,notnull"`
	CallStatus int       `json:"call_status" pg:"call_status,notnull"`
	Created    time.Time `json:"created" pg:"created"`
}
