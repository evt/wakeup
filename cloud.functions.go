package wakeup

import (
	"context"
	"log"
	"net/http"

	"github.com/evt/wakeup/config"
	"github.com/evt/wakeup/db"
	"github.com/evt/wakeup/db/migrations"
	"github.com/evt/wakeup/scheduler"
	"github.com/evt/wakeup/server"
)

var s *server.Server

func init() {
	ctx := context.Background()

	// config
	cfg := config.Get()

	pgDB, err := db.Dial(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Run Postgres migrations
	if err := migrations.Run(pgDB); err != nil {
		log.Fatal(err)
	}

	// create google cloud scheduler client
	sch, err := scheduler.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// create new server instance
	s = server.Init(ctx, cfg, pgDB, sch)
}

// ScheduleCall
func ScheduleCall(w http.ResponseWriter, r *http.Request) {
	s.ScheduleCall(w, r)
}

// CallRoom
func CallRoom(w http.ResponseWriter, r *http.Request) {
	s.CallRoom(w, r)
}
