package server

import (
	"backend/judge"

	"github.com/go-fuego/fuego"
)

// Define the request body structure for posting code
type PostCodeBody struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// PostCode handles the POST request to submit code, runs the code using Piston, and checks results using CheckAll
func (s *Server) PostCode(c fuego.ContextWithBody[PostCodeBody]) (judge.Result, error) {

	_, err := s.db.IncrementTries()
	if err != nil {
		return judge.Result{}, err
	}

	// Extract the code and language from the request body
	body, err := c.Body()
	if err != nil {
		return judge.Result{}, err
	}

	result, err := s.Judge.CheckAll(judge.CheckArgs{
		Code:     body.Code,
		Language: body.Language,
	}) // Check if all tests pass
	if err != nil {
		// Log and return an error if there was an issue with CheckAll
		s.logger.Error("Error checking all test cases", "error", err)
		return result, err
	}

	state, err := s.db.GetState()
	if err != nil {
		return judge.Result{}, err
	}

	if state.Winner != "" {
		result.Verdict = "TL"
		return result, err
	}

	// Return the result of the check (AC or WA) along with Piston response
	return result, nil
}

// Register the judge subroute
func (s *Server) registerJudgeSubrouteOn(parentRoute *fuego.Server) *fuego.Server {
	judgeRoute := fuego.Group(parentRoute, "/judge")
	{
		// Define the POST route for submitting code
		fuego.Post(judgeRoute, "/code", s.PostCode)
	}

	return judgeRoute
}
