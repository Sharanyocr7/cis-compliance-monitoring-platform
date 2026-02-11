package collector

import (
	"bufio"
	"bytes"
	"context"
	"strings"

	"cis-agent/internal/model"
	"cis-agent/internal/util"
)

func PackagesUbuntu(ctx context.Context) ([]model.PackageInfo, error) {
	out, err := util.CmdOut(ctx, "dpkg-query", "-W", "-f=${Package}\t${Version}\t${Architecture}\n")
	if err != nil && out == "" {
		return nil, err
	}

	var pkgs []model.PackageInfo
	sc := bufio.NewScanner(bytes.NewBufferString(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		pkgs = append(pkgs, model.PackageInfo{
			Name:    parts[0],
			Version: parts[1],
			Arch:    parts[2],
		})
	}
	return pkgs, nil
}
