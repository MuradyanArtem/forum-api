package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

// NewServer retrurns server object
func NewServer(store store.Store, sessionStore sessions.Store) *server {
	return &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}
}
