package server

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/evt/wakeup/model"
)

// CallRoom acepts a request from scheduler and calls all rooms tied to wake up time
func (s *Server) CallRoom(w http.ResponseWriter, r *http.Request) {
	callTime := r.URL.Query().Get("call_time")
	if callTime == "" {
		s.error(w, r, errors.New("No call time provided"), http.StatusBadRequest)
		return
	}
	// Find rooms by wake up time
	rooms, err := s.db.FindRooms(callTime)
	if err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Rooms to wake up by time %s:\n%s", callTime, spew.Sdump(rooms))
	// Make sure we have rooms to call
	if len(rooms) == 0 {
		s.respond(w, r, map[string]interface{}{
			"status": "No rooms to call",
		}, http.StatusOK)
		return
	}
	// Call rooms
	var wg sync.WaitGroup
	for _, room := range rooms {
		wg.Add(1)
		go func(room *model.Room) {
			defer wg.Done()
			log.Printf("Calling room %s %s staying in room number %d", room.Firstname, room.Lastname, room.RoomNumber)
		}(room)
	}
	wg.Wait()

	s.respond(w, r, rooms, http.StatusOK)
}
