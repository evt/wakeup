package server

func (s *Server) routes() {
	// index page
	s.router.HandleFunc("/", s.handleIndex())

	s.router.HandleFunc("/wakeup", s.handleWakeUp()).Methods("POST")
}
