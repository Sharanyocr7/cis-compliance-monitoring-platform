package checks

import (
	"context"
	"strings"

	"cis-agent/internal/util"
)

type UFWEnabled struct{}

func (c UFWEnabled) ID() string       { return "CIS-UFW-ENABLED" }
func (c UFWEnabled) Title() string    { return "Firewall (UFW) is enabled" }
func (c UFWEnabled) Severity() string { return "HIGH" }

func (c UFWEnabled) Run(ctx context.Context) Result {
	out, err := util.CmdOut(ctx, "ufw", "status")
	s := strings.ToLower(out)
	if err != nil && out == "" {
		return Result{c.ID(), c.Title(), "ERROR", c.Severity(), err.Error()}
	}
	if strings.Contains(s, "status: active") {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), out}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), out}
}
