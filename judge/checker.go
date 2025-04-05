package judge

import (
	"errors"
	"fmt"
	"strings"
)

type Result struct {
	Verdict        string
	PistonResponse PistonResponse
}

type CheckArgs struct {
	Code     string
	Language string
}

func (j *Judge) CheckAll(args CheckArgs) (Result, error) {
	j.logger.Info(args.Code)
	// Ensure inputs and outputs have the same length
	if len(j.inputs) != len(j.outputs) {
		return Result{}, errors.New("mismatch inputs and outputs")
	}

	// Iterate through inputs and outputs
	j.logger.Info(strings.Join(j.inputs, ", "))
	for i := range j.inputs {
		// Prepare the CheckArgs for the current pair of input and output
		args := CheckOneArgs{
			Code:          args.Code,
			Language:      args.Language,
			Input:         j.inputs[i],
			CorrectOutput: j.outputs[i],
		}

		// Check the current input/output pair
		result, err := j.Check(args)
		if err != nil {
			return Result{}, errors.New("error checking input/output pair")
		}

		// If the result is not "AC", return "WA"
		if result.Verdict != "AC" {
			return result, nil
		}
	}
	j.logger.Info("done")

	// If all checks pass, return "AC" (Accepted)
	return Result{
		Verdict: "AC",
	}, nil
}

type CheckOneArgs struct {
	Code          string
	Language      string
	Input         string
	CorrectOutput string
}

func (j *Judge) Check(args CheckOneArgs) (Result, error) {
	j.logger.Info(fmt.Sprintf("code: %s", args.Code))
	// Run the code using Piston
	resp, err := j.PistonRunner(args.Language, args.Code, args.Input)
	if err != nil {
		j.logger.Error("Failed to execute code with Piston", "error", err)
		return Result{Verdict: "RE", PistonResponse: resp}, err
	}
	j.logger.Info(fmt.Sprintf("out: %s", resp.Run.Stdout))

	actualOutput := resp.Run.Stdout

	// Trim input and correctOutput of all irrelevant characters like trailing whitespaces
	trimmedActual := strings.TrimSpace(actualOutput)
	trimmedCorrect := strings.TrimSpace(args.CorrectOutput)

	// Split by lines
	actualLines := strings.Split(trimmedActual, "\n")
	expectedLines := strings.Split(trimmedCorrect, "\n")

	// Compare number of lines
	if len(actualLines) != len(expectedLines) {
		return Result{
			Verdict:        "WA",
			PistonResponse: resp,
		}, nil
	}

	// Compare each line
	for i := range actualLines {
		if strings.TrimSpace(actualLines[i]) != strings.TrimSpace(expectedLines[i]) {
			return Result{
				Verdict:        "WA",
				PistonResponse: resp,
			}, nil
		}
	}

	// If everything matches
	return Result{
		Verdict:        "AC",
		PistonResponse: resp,
	}, nil
}
