package model

// CallRoomRequest is a request coming from google cloud scheduler to call rooms
type CallRoomRequest struct {
	callTime string `json:"call_time"`
}
