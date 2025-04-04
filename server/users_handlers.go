package server

import (
	"time"

	d "github.com/CollCaz/UniSite/database"
	"github.com/go-fuego/fuego"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) registerAuthSubrouteOn(parentRoute *fuego.Server) *fuego.Server {
	security := fuego.NewSecurity()
	security.ExpiresInterval = 24 * time.Hour
	userRoute := fuego.Group(parentRoute, "/auth")
	fuego.Post(userRoute, "/login", security.LoginHandler(s.LoginUser))
	fuego.Post(userRoute, "", s.PostUser)

	return userRoute
}

type UserToken struct {
	jwt.RegisteredClaims
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func (s *Server) LoginUser(username, password string) (jwt.Claims, error) {
	user, err := s.db.LoginUser(d.LoginUserArgs{
		Email:    username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	token := UserToken{
		RegisteredClaims: jwt.RegisteredClaims{
			//FIXME: the fields
			Issuer:   "CONENT_API_CHANGE_LATER",
			Subject:  user.Username,
			Audience: jwt.ClaimStrings{},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * time.Hour),
			},
			NotBefore: &jwt.NumericDate{},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ID: "",
		},
		Username: user.Username,
		Roles:    user.Roles,
	}

	return token, nil
}

func (s *Server) PostUser(c fuego.ContextWithBody[d.RegisterUserArgs]) (d.User, error) {
	body, err := c.Body()
	if err != nil {
		s.logger.Error(err.Error())
		return d.User{}, err
	}

	return s.db.RegisterUser(body)
}
