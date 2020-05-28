package wakeup

import (
	"context"
	"log"
	"net/http"

	"github.com/evt/wakeup/config"
	"github.com/evt/wakeup/db"
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

	// create new server instance
	s = server.Init(ctx, cfg, pgDB)
}

// WakeUp
func WakeUp(w http.ResponseWriter, r *http.Request) {
	s.WakeUp(w, r)
}
