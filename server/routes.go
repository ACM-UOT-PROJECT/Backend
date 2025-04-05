package server

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
)

func ProtectMutation(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Method: ", r.Method)
			if r.Method == http.MethodGet {
				fmt.Println("Method?: ", r.Method)
				next.ServeHTTP(w, r)
				return
			}
			fmt.Println("Method???: ", r.Method)
			// Do something before the request
			authhWall := fuego.AuthWall(roles...)

			authhWall(next).ServeHTTP(w, r)
			// Do something after the request
		})
	}
}

func (s *Server) RegisterRoutes() {
	s.registerAuthSubrouteOn(s.server)
	apiRoute := fuego.Group(s.server, "/api")
	fuego.Get(apiRoute, "/", s.helloWorld)
	// fuego.Use(apiRoute, ProtectMutation(""))
	{
		s.registerAuthSubrouteOn(apiRoute)
		s.registerJudgeSubrouteOn(apiRoute)
	}
}

func (s *Server) helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello World!", nil
}
