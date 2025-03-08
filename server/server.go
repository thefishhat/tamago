package server

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/yohamta/donburi"
)

type Store interface {
	GetWorld() donburi.World
	GetEntry(id uint32) *donburi.Entry
	GetEntries() map[uint32]*donburi.Entry
}

type Server struct {
	store      Store
	httpServer *http.Server
}

type Config struct {
	Addr string
}

func Start(store Store, cfg Config) (server *Server, err error) {
	log := log.New(log.Writer(), "[server] ", log.LstdFlags)

	server = &Server{
		store: store,
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler.HandleFunc("/", handlePanic(server.listArchetypesHandler))
	handler.HandleFunc("/entities", handlePanic(server.listEntitiesHandler))
	handler.HandleFunc("/entities/{id}", handlePanic(server.getEntityHandler))
	handler.HandleFunc("/entities/{id}/components", handlePanic(server.getEntityHandler))
	handler.HandleFunc("/entities/{entity_id}/components/{component_name}",
		handlePanic(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					server.getComponentHandler(w, r)
				} else if r.Method == http.MethodPut {
					server.setComponentHandler(w, r)
				}
			},
		))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.Addr != "" {
		s.Addr = cfg.Addr
	}

	server.httpServer = s

	log.Println("Starting editor server on ", s.Addr)
	go func() {
		err = s.ListenAndServe()
		if err != nil {
			log.Println("[server] serving http: ", err)
		}
	}()

	return server, nil
}

func (s *Server) Stop() {
	log.Println("Stopping editor server")
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}

func handlePanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Printf("panic: %v\n%s", r, debug.Stack())
			}
		}()

		next(w, r)
	}
}
