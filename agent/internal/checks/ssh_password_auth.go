package checks

import (
	"bufio"
	"context"
	"os"
	"strings"
)

type SSHPasswordAuthDisabled struct{}

func (c SSHPasswordAuthDisabled) ID() string       { return "CIS-SSH-PASSWORD-AUTH" }
func (c SSHPasswordAuthDisabled) Title() string    { return "SSH password authentication disabled" }
func (c SSHPasswordAuthDisabled) Severity() string { return "HIGH" }

func (c SSHPasswordAuthDisabled) Run(ctx context.Context) Result {
	f, err := os.Open("/etc/ssh/sshd_config")
	if err != nil {
		return Result{c.ID(), c.Title(), "ERROR", c.Severity(), err.Error()}
	}
	defer f.Close()

	val := ""
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.EqualFold(fields[0], "PasswordAuthentication") {
			val = fields[1]
		}
	}
	if val == "" {
		return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "PasswordAuthentication not explicitly set"}
	}
	if strings.EqualFold(val, "no") {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "PasswordAuthentication no"}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "PasswordAuthentication " + val}
}
