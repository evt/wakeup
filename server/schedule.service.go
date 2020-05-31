package server

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/evt/wakeup/model"
)

// ScheduleCall accepts a list of rooms to call at certain time, saves them in Postgres and creates scheduler job if not exists yet
func (s *Server) ScheduleCall(w http.ResponseWriter, r *http.Request) {
	var rooms []*model.Room
	if err := s.decode(w, r, &rooms); err != nil {
		s.error(w, r, err, http.StatusBadRequest)
		return
	}
	if len(rooms) == 0 {
		s.error(w, r, errors.New("No rooms provided"), http.StatusBadRequest)
		return
	}
	// Validate room details provided
	for _, room := range rooms {
		if room.RoomNumber == 0 {
			s.error(w, r, fmt.Errorf("No room number specified for guest %s %s", room.Firstname, room.Lastname), http.StatusBadRequest)
			return
		}
		if room.Firstname == "" {
			s.error(w, r, fmt.Errorf("No first name provided for room number %d", room.RoomNumber), http.StatusBadRequest)
			return
		}
		if room.Lastname == "" {
			s.error(w, r, fmt.Errorf("No last name specified for room number %d", room.RoomNumber), http.StatusBadRequest)
			return
		}
		if room.CallTime == "" {
			s.error(w, r, fmt.Errorf("No call time specified for room number %d", room.RoomNumber), http.StatusBadRequest)
			return
		}
		if err := validateCallTime(room.CallTime); err != nil {
			s.error(w, r, fmt.Errorf("Call time (%s) has incorrect format (must be xx:yy) for room number %d", room.CallTime, room.RoomNumber), http.StatusBadRequest)
			return
		}
		// Create scheduler job
		callRoomURL := s.config.CallRoomEndpoint + "?call_time=" + room.CallTime
		if err := s.scheduler.CreateJob(room.CallTime, callRoomURL, s.config.SchedulerLocation, s.config.SchedulerTimeZone); err != nil {
			s.error(w, r, err, http.StatusInternalServerError)
			return
		}
	}
	// Save metadata to Postgres
	if err := s.db.AddRooms(rooms); err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, r, rooms, http.StatusOK)
}
