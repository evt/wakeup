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
      firstname varchar(16) not null,
      lastname varchar(8) not null,
      wakeup_time char(5) not null
    )`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE trace`)
		return err
	})
}
