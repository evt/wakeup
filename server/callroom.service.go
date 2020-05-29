package server

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/evt/wakeup/model"
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
	// Make sure we have users to call
	if len(users) == 0 {
		s.respond(w, r, map[string]interface{}{
			"status": "No users to call",
		}, http.StatusOK)
	}
	// Call users
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user *model.User) {
			defer wg.Done()
			log.Printf("Calling user %s %s staying in room number %d", user.Firstname, user.Lastname, user.RoomNumber)
		}(user)
	}
	wg.Wait()

	s.respond(w, r, users, http.StatusOK)
}
