package checks

import (
	"context"
	"bufio"
	"os"
	"strings"
)

type SSHProtocol2 struct{}

func (c SSHProtocol2) ID() string       { return "CIS-SSH-PROTOCOL2" }
func (c SSHProtocol2) Title() string    { return "SSH uses Protocol 2" }
func (c SSHProtocol2) Severity() string { return "MEDIUM" }

func (c SSHProtocol2) Run(ctx context.Context) Result {
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
		if len(fields) >= 2 && strings.EqualFold(fields[0], "Protocol") {
			val = fields[1]
		}
	}
	// On Ubuntu 22+, protocol 2 is default; if unset, we treat as PASS
	if val == "" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "Protocol not set (Ubuntu defaults to 2)"}
	}
	if val == "2" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "Protocol 2"}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "Protocol " + val}
}
