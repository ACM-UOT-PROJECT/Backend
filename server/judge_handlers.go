package server

import (
	"backend/judge"
	"fmt"

	"github.com/go-fuego/fuego"
)

// Define the request body structure for posting code
type PostCodeBody struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// PostCode handles the POST request to submit code, runs the code using Piston, and checks results using CheckAll
func (s *Server) PostCode(c fuego.ContextWithBody[PostCodeBody]) (judge.Result, error) {
	s.logger.Info("Posting Code")

	s.logger.Info("Incrementing Tries")
	_, err := s.db.IncrementTries()
	if err != nil {
		return judge.Result{}, err
	}
	s.logger.Info("Incremented Tries")

	s.logger.Info("Getting Body")
	// Extract the code and language from the request body
	body, err := c.Body()
	if err != nil {
		return judge.Result{}, err
	}
	s.logger.Info("Got Body")

	s.logger.Info("Judging Code")
	result, err := s.Judge.CheckAll(judge.CheckArgs{
		Code:     body.Code,
		Language: body.Language,
	}) // Check if all tests pass
	if err != nil {
		// Log and return an error if there was an issue with CheckAll
		s.logger.Error("Error checking all test cases", "error", err)
		return result, err
	}
	s.logger.Info("Judged Code")
	s.logger.Info(result.Verdict)

	state, err := s.db.GetState()
	if err != nil {
		return judge.Result{}, err
	}

	if state.Winner != "NO_ONE_WON_AT_ALL" && result.Verdict == "AC" {
		s.logger.Info(fmt.Sprintf("Winner: %s", state.Winner))
		result.Verdict = "Too Late"
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
