package server

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

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
	// Make sure we have rooms to call
	if len(rooms) == 0 {
		s.respond(w, r, map[string]interface{}{
			"status": "No rooms to call",
		}, http.StatusOK)
		return
	}
	// Call rooms
	var g errgroup.Group
	for _, room := range rooms {
		room := room // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			log.Printf("[Room %d] calling guest - %s %s", room.RoomNumber, room.Firstname, room.Lastname)
			status, err := call(s.config.CallEndpoint, room)
			if err != nil {
				log.Printf("[Room %d] call error: %s, call status - %d", room.RoomNumber, err, status)
				return err
			}
			log.Printf("[Room %d] call status - %d", room.RoomNumber, status)
			// Schedule call retry if needed
			needsCallRetry, err := s.needsCallRetry(room, status)
			if err != nil {
				log.Printf("[Room %d] needsCallRetry error: %s", room.RoomNumber, err)
				return err
			}
			if needsCallRetry {
				callTime, err := addRetryPeriod(room.CallTime, s.config.SchedulerRetryPeriod)
				if err != nil {
					log.Printf("[Room %d] addRetryPeriod error: %s", room.RoomNumber, err)
					return err
				}
				room.CallTime = callTime
				s.scheduleJob([]*model.Room{room})
				s.db.IncRoomRetryCount(room)
			} else {
				log.Printf("[Room %d] retry count - %d. No more retries scheduled.", room.RoomNumber, room.RetryCount)
			}
			// Save call in DB
			err = s.db.SaveCall(&model.Call{
				RoomNumber: room.RoomNumber,
				CallStatus: status,
			})
			if err != nil {
				log.Printf("[Room %d] save call error: %s", room.RoomNumber, err)
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}
	s.respond(w, r, rooms, http.StatusOK)
}

// call makes a REST call to external call service.
// Success status: 200
// Error status: non-200
func call(endpoint string, room *model.Room) (int, error) {
	if endpoint == "" {
		return 0, errors.New("No call endpoint found in config")
	}
	if room == nil {
		return 0, errors.New("No room provided")
	}
	resp, err := HTTPClient.Get(endpoint)
	if err != nil {
		return 0, errors.Wrap(err, "calll->HTTPClient.Get")
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// needsCallRetry returns true if status is 4xx-5xx
func (s *Server) needsCallRetry(room *model.Room, status int) (bool, error) {
	if room == nil {
		return false, errors.New("No room provided")
	}
	if room.RetryCount >= s.config.SchedulerRetryCount {
		return false, nil
	}
	if status >= 400 && status < 600 {
		return true, nil
	}
	return false, nil
}
