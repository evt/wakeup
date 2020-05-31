package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table rooms...")
		_, err := db.Exec(`CREATE TABLE rooms(
			room_number integer not null primary key,
      firstname varchar(64) not null,
      lastname varchar(64) not null,
      call_time char(5) not null,
			created timestamp default current_timestamp
    )`)
		if err != nil {
			return err
		}
		fmt.Println("creating table calls...")
		_, err = db.Exec(`CREATE TABLE calls(
			call_id serial not null primary key,
			room_number integer not null,
			call_status integer not null,
			created timestamp default current_timestamp
    )`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table rooms...")
		_, err := db.Exec(`DROP TABLE rooms`)
		if err != nil {
			return err
		}
		fmt.Println("dropping table calls...")
		_, err = db.Exec(`DROP TABLE calls`)
		return err
	})
}
