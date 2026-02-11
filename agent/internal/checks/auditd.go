package checks

import (
	"context"

	"cis-agent/internal/util"
)

type AuditdEnabled struct{}

func (c AuditdEnabled) ID() string       { return "CIS-AUDITD" }
func (c AuditdEnabled) Title() string    { return "Audit daemon (auditd) installed and running" }
func (c AuditdEnabled) Severity() string { return "HIGH" }

func (c AuditdEnabled) Run(ctx context.Context) Result {
	_, err := util.CmdOut(ctx, "dpkg", "-s", "auditd")
	if err != nil {
		return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "auditd not installed"}
	}
	out, _ := util.CmdOut(ctx, "systemctl", "is-active", "auditd")
	if out == "active" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "auditd active"}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "auditd status: " + out}
}
