package model

// User is a usere to wake up
type User struct {
	tableName  struct{} `pg:"users"`
	UserID     string   `json:"user_id" pg:"user_id,notnull,unique"`
	Firstname  string   `json:"firstname" pg:"firstname,notnull"`
	Lastname   string   `json:"lastname" pg:"lastname,notnull"`
	Phone      string   `json:"phone" pg:"phone"`
	RoomNumber int      `json:"room_number" pg:"room_number,notnull"`
	WakeUpTime string   `json:"wakeup_time" pg:"wakeup_time,notnull"`
}
