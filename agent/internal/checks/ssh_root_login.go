package checks

import (
	"bufio"
	"context"
	"os"
	"strings"
)

type SSHRootLoginDisabled struct{}

func (c SSHRootLoginDisabled) ID() string       { return "CIS-SSH-ROOT-LOGIN" }
func (c SSHRootLoginDisabled) Title() string    { return "Root login disabled over SSH" }
func (c SSHRootLoginDisabled) Severity() string { return "HIGH" }

func (c SSHRootLoginDisabled) Run(ctx context.Context) Result {
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
		if len(fields) >= 2 && strings.EqualFold(fields[0], "PermitRootLogin") {
			val = fields[1]
		}
	}
	if val == "" {
		return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "PermitRootLogin not explicitly set"}
	}
	if strings.EqualFold(val, "no") || strings.EqualFold(val, "prohibit-password") {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "PermitRootLogin " + val}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "PermitRootLogin " + val}
}
