package server

import (
	d "backend/database"
	"github.com/go-fuego/fuego"
)

func (s *Server) registerStateSubrouteOn(parentRoute *fuego.Server) *fuego.Server {
	stateRoute := fuego.Group(parentRoute, "/state")
	{
		fuego.Get(stateRoute, "/", s.GetState)
		fuego.Post(stateRoute, "", s.PostUser)
	}

	return stateRoute
}

func (s *Server) GetState(ctx fuego.ContextNoBody) (d.State, error) {
	s.logger.Info("Getting state...")

	return s.db.GetState()
}
