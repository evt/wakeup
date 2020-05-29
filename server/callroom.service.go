package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// CallRoom acepts a request from scheduler and calls all users tied to wake up time
func (s *Server) CallRoom(w http.ResponseWriter, r *http.Request) {
	wakeUpTime := r.URL.Query().Get("wakeup_time")
	if wakeUpTime == "" {
		s.error(w, r, errors.New("No wake up time provided"), http.StatusBadRequest)
		return
	}
	// Find users by wake up time
	users, err := s.db.FindUsers(wakeUpTime)
	if err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Users to wake up by time %s:\n%s", wakeUpTime, spew.Sdump(users))
	s.respond(w, r, users, http.StatusOK)
}
