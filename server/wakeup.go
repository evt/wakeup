package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/evt/wakeup/model"
	"github.com/google/uuid"
)

func (s *Server) WakeUp(w http.ResponseWriter, r *http.Request) {
	var payload []*model.User
	if err := s.decode(w, r, &payload); err != nil {
		s.error(w, r, err, http.StatusBadRequest)
		return
	}
	if len(payload) == 0 {
		s.error(w, r, errors.New("No users provided"), http.StatusBadRequest)
		return
	}
	for _, user := range payload {
		// Validate payload
		if user.Firstname == "" {
			s.error(w, r, errors.New("No first name provided"), http.StatusBadRequest)
			return
		}
		if user.Lastname == "" {
			s.error(w, r, fmt.Errorf("No last name specified for user %s", user.Firstname), http.StatusBadRequest)
			return
		}
		if user.WakeUpTime == "" {
			s.error(w, r, fmt.Errorf("No wake up time specified for user %s %s", user.Firstname, user.Lastname), http.StatusBadRequest)
			return
		}
		userUUID, err := uuid.NewRandom()
		if err != nil {
			log.Printf("uuid.NewRandom error: %s", err)
			s.error(w, r, fmt.Errorf("Can't generate UUID for user %s %s", user.Firstname, user.Lastname), http.StatusBadRequest)
			return
		}
		user.UserID = userUUID.String()
	}
	// Save metadata to Postgres
	if err := s.db.Insert(&payload); err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, r, ok, http.StatusOK)
}
