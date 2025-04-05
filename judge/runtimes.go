package judge

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type ApiRunTime struct {
	Language string   `json:"language"`
	Version  string   `json:"version"`
	Aliases  []string `json:"aliases"`
	Runtime  string   `json:"runtime"`
}

type ApiRunTimes []ApiRunTime

var apiRunTimes ApiRunTimes

type RunTimes map[string]ApiRunTime

var runTimes RunTimes

func GetRunTimes() (RunTimes, error) {
	if runTimes != nil {
		return runTimes, nil
	}

	const (
		url            = "https://emkc.org/api/v2/piston/runtimes"
		maxRetries     = 3
		initialBackoff = 500 * time.Millisecond
	)

	var (
		resp *http.Response
		err  error
	)

	backoff := initialBackoff
	for i := range maxRetries {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		if i < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
	}

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch runtimes from API: non-200 response")
	}

	var fetched ApiRunTimes
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		return nil, err
	}

	apiRunTimes = fetched
	runTimes = make(RunTimes)
	for _, rt := range apiRunTimes {
		runTimes[rt.Language] = rt
		for _, alias := range rt.Aliases {
			runTimes[alias] = rt
		}
	}

	return runTimes, nil
}
