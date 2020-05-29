package server

import (
	"errors"
	"net/http"

	"github.com/evt/wakeup/model"
)

// CallRoom acepts a request from scheduler and calls all users tied to wake up time
func (s *Server) CallRoom(w http.ResponseWriter, r *http.Request) {
	var payload model.CallRoomRequest
	if err := s.decode(w, r, &payload); err != nil {
		s.error(w, r, err, http.StatusBadRequest)
		return
	}
	if payload.WakeUpTime == "" {
		s.error(w, r, errors.New("No wake up time provided"), http.StatusBadRequest)
		return
	}

	s.respond(w, r, ok, http.StatusOK)
}
