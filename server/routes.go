package server

import (
	"Goo/handlers"
)

func (s *Server) setupRoutes() {
	handlers.Health(s.mux, s.database)
	handlers.FrontPage(s.mux)

	handlers.NewsletterSignup(s.mux, s.database, nil)
	handlers.NewsletterThanks(s.mux)
}
