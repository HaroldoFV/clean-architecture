package webserver

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	if s.Handlers[path] == nil {
		s.Handlers[path] = make(map[string]http.HandlerFunc)
	}
	s.Handlers[path][method] = handler
}

func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)

	for path, methodHandlers := range s.Handlers {
		s.Router.Route(path, func(r chi.Router) {
			for method, handler := range methodHandlers {
				switch method {
				case http.MethodGet:
					r.Get("/", handler)
				case http.MethodPost:
					r.Post("/", handler)
				case http.MethodPut:
					r.Put("/", handler)
				case http.MethodDelete:
					r.Put("/", handler)
				}
			}
		})
	}
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
