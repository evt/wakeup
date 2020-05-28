package server

func (s *Server) routes() {
	// index page
	s.router.HandleFunc("/", s.handleIndex())
}
