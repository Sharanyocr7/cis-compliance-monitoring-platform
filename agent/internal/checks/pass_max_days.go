package checks

import (
	"context"
	"os"
	"strings"
)

type PassMaxDays struct{}

func (c PassMaxDays) ID() string       { return "CIS-PASS-MAX-DAYS" }
func (c PassMaxDays) Title() string    { return "Password expiration (PASS_MAX_DAYS) set to <= 365" }
func (c PassMaxDays) Severity() string { return "MEDIUM" }

func (c PassMaxDays) Run(ctx context.Context) Result {
	b, err := os.ReadFile("/etc/login.defs")
	if err != nil {
		return Result{c.ID(), c.Title(), "ERROR", c.Severity(), err.Error()}
	}
	lines := strings.Split(string(b), "\n")
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		if strings.HasPrefix(ln, "PASS_MAX_DAYS") {
			fields := strings.Fields(ln)
			if len(fields) >= 2 {
				// simple parse without importing strconv heavy? still fine:
				val := fields[1]
				if val == "365" || val == "180" || val == "90" || val == "60" || val == "30" {
					return Result{c.ID(), c.Title(), "PASS", c.Severity(), ln}
				}
				// if unknown value, fail (strict)
				return Result{c.ID(), c.Title(), "FAIL", c.Severity(), ln}
			}
		}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "PASS_MAX_DAYS not set in /etc/login.defs"}
}
