package judge

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type PistonRequest struct {
	Language string `json:"language"`
	Version  string `json:"version"`
	Files    []struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	} `json:"files"`
	Input string   `json:"stdin"`
	Args  []string `json:"args,omitempty"`
}

type PistonResponse struct {
	Run struct {
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
		ExitCode int    `json:"exit_code"`
	} `json:"run"`
}

func (j *Judge) PistonRunner(language, code, input string) (PistonResponse, error) {
	j.logger.Info("Starting Piston execution", "language", language)

	runtimes, err := GetRunTimes()
	if err != nil {
		j.logger.Error("Error fetching runtimes", "error", err)
		return PistonResponse{}, err
	}

	rt, found := runtimes[language]
	if !found {
		err := errors.New("runtime not found for language: " + language)
		j.logger.Error("Runtime not found", "language", language)
		return PistonResponse{}, err
	}

	reqBody := PistonRequest{
		Language: rt.Language,
		Version:  rt.Version,
		Files: []struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}{
			{Name: "main." + language, Content: code},
		},
		Input: input,
		Args:  []string{},
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		j.logger.Error("Error marshaling request", "error", err)
		return PistonResponse{}, err
	}

	const (
		url            = "https://emkc.org/api/v2/piston/execute"
		maxRetries     = 10
		initialBackoff = 500 * time.Millisecond
	)

	var pistonResp PistonResponse
	backoff := initialBackoff
	var lastStatusCode int

	for i := range maxRetries {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
		if err != nil {
			j.logger.Error("HTTP request failed", "attempt", i+1, "error", err)
			return PistonResponse{}, err
		}
		lastStatusCode = resp.StatusCode

		if resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(&pistonResp); err != nil {
				j.logger.Error("Failed to decode Piston response", "error", err)
				return PistonResponse{}, err
			}
			break
		}

		resp.Body.Close()
		time.Sleep(backoff)
		backoff *= 2
	}

	// Handle unsuccessful final status code
	if lastStatusCode != http.StatusOK {
		err := errors.New("Piston API returned non-200 status after retries")
		j.logger.Error("Final request failed", "status", lastStatusCode)
		return PistonResponse{}, err
	}

	// Log runtime stderr
	if pistonResp.Run.Stderr != "" {
		j.logger.Error("Runtime error", "stderr", pistonResp.Run.Stderr)
		return pistonResp, nil
	}

	j.logger.Info("Execution succeeded",
		"stdout", pistonResp.Run.Stdout,
		"exit_code", pistonResp.Run.ExitCode,
	)

	return pistonResp, nil
}
