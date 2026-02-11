package checks

import (
	"context"
	"os"
	"strings"
)

type IPForwardingDisabled struct{}

func (c IPForwardingDisabled) ID() string       { return "CIS-IP-FORWARDING" }
func (c IPForwardingDisabled) Title() string    { return "IPv4 forwarding disabled (net.ipv4.ip_forward=0)" }
func (c IPForwardingDisabled) Severity() string { return "MEDIUM" }

func (c IPForwardingDisabled) Run(ctx context.Context) Result {
	b, err := os.ReadFile("/proc/sys/net/ipv4/ip_forward")
	if err != nil {
		return Result{c.ID(), c.Title(), "ERROR", c.Severity(), err.Error()}
	}
	v := strings.TrimSpace(string(b))
	if v == "0" {
		return Result{c.ID(), c.Title(), "PASS", c.Severity(), "net.ipv4.ip_forward=0"}
	}
	return Result{c.ID(), c.Title(), "FAIL", c.Severity(), "net.ipv4.ip_forward=" + v}
}
