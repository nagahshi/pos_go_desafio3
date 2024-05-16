package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type handler struct {
	fn     http.HandlerFunc
	method string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]handler
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]handler),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method string, path string, fn http.HandlerFunc) {
	s.Handlers[path] = handler{
		fn:     fn,
		method: method,
	}
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Method(handler.method, path, handler.fn)
	}
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}
