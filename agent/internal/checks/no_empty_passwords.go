package checks

import (
	"context"

	"cis-agent/internal/util"
)

type NoEmptyPasswords struct{}

func (c NoEmptyPasswords) ID() string       { return "CIS-NO-EMPTY-PASSWORDS" }
func (c NoEmptyPasswords) Title() string    { return "No users have empty password fields" }
func (c NoEmptyPasswords) Severity() string { return "HIGH" }

func (c NoEmptyPasswords) Run(ctx context.Context) Result {
	out, err := util.CmdOut(ctx, "awk", "-F:", `($2==""){print $1}`, "/etc/shadow")
	if err != nil {
		return Result{c.ID(), c.Title(), "ERROR", c.Severity(), err.Error()}
	}
	if out == "" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "no empty password hashes"}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "users with empty password field: " + out}
}
