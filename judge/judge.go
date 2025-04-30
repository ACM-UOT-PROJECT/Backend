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
			"3 6",
			"7 5",
			"93 35",
			"86 92",
			"49 21",
			"62 27",
			"690 59",
			"763 926",
			"540 426",
			"172 736",
		},
		Outputs: []string{
			"9",
			"4",
			"312",
			"36",
			"90",
			"289",
			"4725",
			"673",
			"368",
			"2610",
		},
	})

	logger.Info(strings.Join(judge.inputs, ", "))

	return judge
}
