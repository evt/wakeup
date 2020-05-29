package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table users...")
		_, err := db.Exec(`CREATE TABLE users(
      user_id uuid not null primary key,
      firstname varchar(64) not null,
      lastname varchar(64) not null,
			phone varchar(32) default null,
			room_number integer not null,
      wakeup_time char(5) not null
    )`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
}
