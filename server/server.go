package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/evt/wakeup/config"
	"github.com/evt/wakeup/db"
	"github.com/gorilla/mux"
)

// Server is a server
type Server struct {
	context context.Context
	config  *config.Config
	router  *mux.Router
	db      *db.PgDB
}

// Init returns new server instance
func Init(ctx context.Context, config *config.Config, db *db.PgDB) *Server {
	router := mux.NewRouter()
	s := &Server{
		context: ctx,
		config:  config,
		router:  router,
		db:      db,
	}
	s.routes()
	return s
}

// respond responds with JSON response
func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			// TODO
		}
	}
}

// error responds with JSON response
func (s *Server) error(w http.ResponseWriter, r *http.Request, err error, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err != nil {
		err := json.NewEncoder(w).Encode(e(err))
		if err != nil {
			// TODO
		}
	}
}

// decode decodes incoming JSON request
func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// ServeHTTP makes our server http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
