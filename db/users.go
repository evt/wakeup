package db

import "github.com/evt/wakeup/model"

// AddUsers saves a list of users in Postgres
func (db *PgDB) AddUsers(users []*model.User) error {
	return db.Insert(&users)
}

// FindUsers returns user list by wake up time
func (db *PgDB) FindUsers(wakeUpTime string) ([]*model.User, error) {
	var users []*model.User
	err := db.Model(&users).Where("wakeup_time = ?", wakeUpTime).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}
