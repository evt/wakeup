package server

import (
	"errors"
	"net/http"
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

	s.respond(w, r, users, http.StatusOK)
}
