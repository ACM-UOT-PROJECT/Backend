package judge

import (
	"log/slog"
	"strings"
)

type Judge struct {
	logger  *slog.Logger
	output  string
	inputs  []string
	outputs []string
}

type NewJudgeArgs struct {
	Logger  *slog.Logger
	Inputs  []string
	Outputs []string
}

func NewJudge(args NewJudgeArgs) *Judge {
	if args.Logger == nil {
		args.Logger = slog.Default()
	}

	return &Judge{
		logger:  args.Logger,
		inputs:  args.Inputs,
		outputs: args.Outputs,
	}
}

func CreateJudge(logger slog.Logger) *Judge {
	logger.Info("hello")
	judge := NewJudge(NewJudgeArgs{
		Logger: &logger,
		Inputs: []string{
			"1 2",
			"2 4",
			"6 8",
		},
		Outputs: []string{
			"3",
			"3",
			"3",
		},
	})

	logger.Info(strings.Join(judge.inputs, ", "))

	return judge
}
