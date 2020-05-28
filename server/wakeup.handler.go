package server

import (
	"net/http"
)

// handleWakeUp
func (s *Server) handleWakeUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.WakeUp(w, r)
	}
}
