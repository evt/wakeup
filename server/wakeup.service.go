package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/evt/wakeup/model"
	"github.com/evt/wakeup/scheduler"
	"github.com/google/uuid"
)

// WakeUp accepts a list of users to wake up at certain time
func (s *Server) WakeUp(w http.ResponseWriter, r *http.Request) {
	var users []*model.User
	if err := s.decode(w, r, &users); err != nil {
		s.error(w, r, err, http.StatusBadRequest)
		return
	}
	if len(users) == 0 {
		s.error(w, r, errors.New("No users provided"), http.StatusBadRequest)
		return
	}
	for _, user := range users {
		// Validate user details provided
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
		if user.RoomNumber == 0 {
			s.error(w, r, fmt.Errorf("No room number specified for user %s %s", user.Firstname, user.Lastname), http.StatusBadRequest)
			return
		}
		userUUID, err := uuid.NewRandom()
		if err != nil {
			log.Printf("uuid.NewRandom error: %s", err)
			s.error(w, r, fmt.Errorf("Can't generate UUID for user %s %s", user.Firstname, user.Lastname), http.StatusBadRequest)
			return
		}
		user.UserID = userUUID.String()
		// Create scheduler job
		callRoomURL := s.config.CallRoomEndpoint + "?wakeup_time=" + user.WakeUpTime
		if err := scheduler.CreateJob(context.Background(), user.WakeUpTime, callRoomURL); err != nil {
			s.error(w, r, err, http.StatusInternalServerError)
			return
		}
	}
	// Save metadata to Postgres
	if err := s.db.AddUsers(users); err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, r, users, http.StatusOK)
}
