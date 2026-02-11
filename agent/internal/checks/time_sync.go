package checks

import (
	"context"

	"cis-agent/internal/util"
)

type TimeSyncEnabled struct{}

func (c TimeSyncEnabled) ID() string       { return "CIS-TIME-SYNC" }
func (c TimeSyncEnabled) Title() string    { return "Time synchronization service enabled" }
func (c TimeSyncEnabled) Severity() string { return "MEDIUM" }

func (c TimeSyncEnabled) Run(ctx context.Context) Result {
	// Accept either systemd-timesyncd or chrony
	out1, _ := util.CmdOut(ctx, "systemctl", "is-active", "systemd-timesyncd")
	out2, _ := util.CmdOut(ctx, "systemctl", "is-active", "chrony")
	if out1 == "active" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "systemd-timesyncd active"}
	}
	if out2 == "active" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "chrony active"}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "timesyncd=" + out1 + ", chrony=" + out2}
}
