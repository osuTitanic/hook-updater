package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Server   *Server
}

type Server struct {
	Config  *Config
	Logger  *Logger
	Manager *ReleaseManager
}

func NewServer(config *Config, manager *ReleaseManager) *Server {
	return &Server{
		Logger:  CreateLogger("Server", INFO),
		Config:  config,
		Manager: manager,
	}
}

func (server *Server) Bind() string {
	return fmt.Sprintf("%s:%d", server.Config.Server.Host, server.Config.Server.Port)
}

func (server *Server) Serve() {
	server.Logger.Infof("Listening on %s", server.Bind())

	r := mux.NewRouter()
	r.HandleFunc("/update", server.contextMiddleware(UpdateHandler)).Methods("GET")

	loggedMux := server.loggingMiddleware(r)
	http.ListenAndServe(server.Bind(), loggedMux)
}

func (server *Server) contextMiddleware(handler func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := &Context{
			Response: w,
			Request:  r,
			Server:   server,
		}

		w.Header().Set("Server", "Titanic")
		handler(context)
	}
}

func (server *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		server.Logger.Infof(
			"%s - %s %s (%v)",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}
