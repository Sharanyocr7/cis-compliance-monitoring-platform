package checks

import "context"

type Result struct {
	CheckID  string
	Title    string
	Status   string // PASS/FAIL/ERROR
	Severity string
	Evidence string
}

type Check interface {
	ID() string
	Title() string
	Severity() string
	Run(ctx context.Context) Result
}
