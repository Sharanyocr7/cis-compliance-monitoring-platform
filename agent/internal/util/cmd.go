package util

import (
	"context"
	"os/exec"
	"strings"
)

func CmdOut(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
