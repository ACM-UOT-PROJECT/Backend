package server

import (
	"net/http"
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
	cookieToken, err := ctx.Cookie("acm_jwt_token")
	if err == nil {
		if cookieToken != nil {
			_, err := s.server.Security.ValidateToken(cookieToken.Value)
			if err == nil {
				s.logger.Info("User already has a token")
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
	token, err := s.server.Security.GenerateToken(claims)
	if err != nil {
		return "could not generate token", err
	}

	ctx.SetCookie(http.Cookie{
		Name:     "acm_jwt_token",
		Value:    token,
		Path:     "/",
		Secure:   false,                // Set to true if using HTTPS (recommended in production)
		HttpOnly: true,                 // Prevents JavaScript access (for security)
		SameSite: http.SameSiteLaxMode, // Required for cross-origin requests
		MaxAge:   86400,                // 1 day expiration, adjust as needed
	})

	s.db.IncrementPeople()

	_, err = s.db.RegisterUser(d.RegisterUserArgs{
		UserName: "",
		Token:    token,
	})
	if err != nil {
		return "", err
	}

	// Return the token
	return token, nil
}

type PostUserBody struct {
	UserName string
}

func (s *Server) PostUser(c fuego.ContextWithBody[PostUserBody]) (d.User, error) {
	body, err := c.Body()
	if err != nil {
		s.logger.Error(err.Error())
		return d.User{}, err
	}
	cookie, err := c.Cookie("jwt_token")
	if err != nil {
		return d.User{}, err
	}

	args := d.RegisterUserArgs{
		UserName: body.UserName,
		Token:    cookie.Value,
	}

	_, err = s.db.SetWinner(d.SetWinnerArgs{
		Winner: body.UserName,
	})
	if err != nil {
		return d.User{}, err
	}

	return s.db.RegisterUser(args)
}
