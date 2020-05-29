package server

func (s *Server) routes() {
	// index page
	s.router.HandleFunc("/", s.handleIndex())
	// user wants to wake them up at xx:yy
	s.router.HandleFunc("/wakeup", s.handleWakeUp()).Methods("POST")
	// call all users by wake up time (expected to run by google cloud scheduler)
	s.router.HandleFunc("/callrooom", s.handleCallRoom()).Methods("POST")
}
