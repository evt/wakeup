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
	rooms, err := s.db.FindRooms(callTime, s.config.SchedulerMaxRetryCount)
	if err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Found %d rooms to call for a call time %s and SchedulerMaxRetryCount = %d", len(rooms), callTime, s.config.SchedulerMaxRetryCount)
	// Make sure we have rooms to call
	if len(rooms) == 0 {
		s.respond(w, r, map[string]interface{}{
			"status": "No rooms to call",
		}, http.StatusOK)
		return
	}
	// Call rooms
	var g errgroup.Group
	// Collect calls
	calls := make([]*model.Call, 0, len(rooms))
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
			// Check if call passed
			callPassed, err := s.callPassed(room, status)
			if err != nil {
				log.Printf("[Room %d] callPassed error: %s", room.RoomNumber, err)
				return err
			}
			// Schedule call retry if needed
			needsCallRetry, err := s.needsCallRetry(room, status)
			if err != nil {
				log.Printf("[Room %d] needsCallRetry error: %s", room.RoomNumber, err)
				return err
			}
			if callPassed {
				log.Printf("[Room %d] call passed (status - %d).", room.RoomNumber, status)
			} else {
				if needsCallRetry {
					retryCallTime, err := addRetryPeriod(room.CallTime, s.config.SchedulerRetryPeriod)
					if err != nil {
						log.Printf("[Room %d] addRetryPeriod error: %s", room.RoomNumber, err)
						return err
					}

					log.Printf("[Room %d] room retry count - %d. Scheduling retry call at %s.", room.RoomNumber, room.RetryCount, retryCallTime)

					room.CallTime = retryCallTime
					s.scheduleJob([]*model.Room{room})

					s.db.IncRoomRetryCount(room)
				} else {
					log.Printf("[Room %d] room retry count - %d. No more retries scheduled.", room.RoomNumber, room.RetryCount)
				}
			}
			// Save call in DB
			call := &model.Call{
				RoomNumber: room.RoomNumber,
				CallStatus: status,
			}
			err = s.db.SaveCall(call)
			if err != nil {
				log.Printf("[Room %d] save call error: %s", room.RoomNumber, err)
				return err
			}
			calls = append(calls, call)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		s.error(w, r, err, http.StatusInternalServerError)
		return
	}
	s.respond(w, r, calls, http.StatusOK)
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

	log.Printf("[Room %d] calling 3rd party service API on %s", room.RoomNumber, endpoint)

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
	if room.RetryCount >= s.config.SchedulerMaxRetryCount {
		return false, nil
	}
	if status >= 400 && status < 600 {
		return true, nil
	}
	return false, nil
}

// callPassed returns true if status is 200
func (s *Server) callPassed(room *model.Room, status int) (bool, error) {
	if room == nil {
		return false, errors.New("No room provided")
	}
	if status == 200 {
		return true, nil
	}
	return false, nil
}
