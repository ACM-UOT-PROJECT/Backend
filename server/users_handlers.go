package server

import (
	"time"

	d "backend/database"
	"github.com/go-fuego/fuego"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) registerAuthSubrouteOn(parentRoute *fuego.Server) *fuego.Server {
	userRoute := fuego.Group(parentRoute, "/auth")
	{
		fuego.Get(userRoute, "/issueToken", s.IssueToken)
		fuego.Post(userRoute, "", s.PostUser)
	}

	return userRoute
}

type UserToken struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func (s *Server) IssueToken(ctx fuego.ContextNoBody) (string, error) {
	s.logger.Info("starting...")
	// Check for the token in cookies
	cookieToken, err := ctx.Cookie("jwt_token")
	if err == nil {
		if cookieToken != nil {
			_, err := s.server.Security.ValidateToken(cookieToken.Value)
			if err == nil {
				return "already have token", nil
			}
		}
	}
	s.logger.Info("no cookie")

	// Create JWT claims with no expiration (effectively never expires)
	claims := UserToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Backend",                          // Issuer of the token
			Subject:   "anon",                             // Subject (can be a user ID or email)
			ExpiresAt: nil,                                // Never expire
			NotBefore: &jwt.NumericDate{},                 // Token valid immediately
			IssuedAt:  &jwt.NumericDate{Time: time.Now()}, // Token issued now
		},
	}

	// Generate token and store in cookie
	token, err := s.server.Security.GenerateTokenToCookies(claims, ctx.Response())
	if err != nil {
		return "already have token", err
	}

	// Optionally, store token in DB or log it
	_, err = s.db.RegisterUser(d.RegisterUserArgs{
		Username: "",
		Token:    token,
	})
	if err != nil {
		return "", err
	}

	// Return the token
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
