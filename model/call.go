package model

// CallRoomRequest is a request coming from google cloud scheduler to call users
type CallRoomRequest struct {
	WakeUpTime string `json:"wakeup_time"`
}
