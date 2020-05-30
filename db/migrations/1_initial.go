package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table rooms...")
		_, err := db.Exec(`CREATE TABLE rooms(
      room_id uuid not null primary key,
      firstname varchar(64) not null,
      lastname varchar(64) not null,
			room_number integer not null unique,
      call_time char(5) not null,
			created timestamp default current_timestamp
    )`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table rooms...")
		_, err := db.Exec(`DROP TABLE rooms`)
		return err
	})
}
