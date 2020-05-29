package db

import "github.com/evt/wakeup/model"

// AddUsers saves a list of users in Postgres
func (db *PgDB) AddUsers(users []*model.User) error {
	return db.Insert(&users)
}
