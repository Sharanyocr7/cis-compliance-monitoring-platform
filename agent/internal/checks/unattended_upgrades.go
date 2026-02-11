package checks

import (
	"context"
	"os"

	"cis-agent/internal/util"
)

type UnattendedUpgradesEnabled struct{}

func (c UnattendedUpgradesEnabled) ID() string       { return "CIS-AUTO-UPDATES" }
func (c UnattendedUpgradesEnabled) Title() string    { return "Automatic security updates enabled (unattended-upgrades)" }
func (c UnattendedUpgradesEnabled) Severity() string { return "MEDIUM" }

func (c UnattendedUpgradesEnabled) Run(ctx context.Context) Result {
	// quick: file present + service active is good enough for this project
	if _, err := os.Stat("/etc/apt/apt.conf.d/50unattended-upgrades"); err != nil {
		return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "missing /etc/apt/apt.conf.d/50unattended-upgrades"}
	}
	out, err := util.CmdOut(ctx, "systemctl", "is-enabled", "unattended-upgrades")
	if err == nil && (out == "enabled" || out == "static") {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "unattended-upgrades is " + out}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "unattended-upgrades enabled? " + out}
}
