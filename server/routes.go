package server

// routes lists routes for our HTTP server
func (s *Server) routes() {
	// index page
	s.router.HandleFunc("/", s.handleIndex())
	// schedule room calls
	s.router.HandleFunc("/schedule", s.handleCallSchedule()).Methods("POST")
	// call rooms
	s.router.HandleFunc("/call", s.handleCallRoom()).Methods("GET")
}
