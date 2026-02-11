package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"cis-agent/internal/checks"
	"cis-agent/internal/collector"
	"cis-agent/internal/config"
	"cis-agent/internal/model"
	"cis-agent/internal/transport"
	"cis-agent/internal/util"
)

func main() {
	var cfgPath string
	var dryRun bool
	flag.StringVar(&cfgPath, "config", "../config.yaml", "path to config.yaml")
	flag.BoolVar(&dryRun, "dry-run", false, "print JSON only, do not send to backend")
	flag.Parse()

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout())
	defer cancel()

	hostname := collector.Hostname()
	kernel, _ := util.CmdOut(ctx, "uname", "-r")
	osName, osVer := collector.OSInfo()
	ip := collector.PrimaryIPv4()
	hostID := collector.StableHostID(hostname, osName, osVer, kernel)

	pkgs, _ := collector.PackagesUbuntu(ctx)

	checkList := []checks.Check{
		checks.UFWEnabled{},
		checks.SSHRootLoginDisabled{},
		checks.SSHPasswordAuthDisabled{},
		checks.UnattendedUpgradesEnabled{},
		checks.TimeSyncEnabled{},
		checks.NoEmptyPasswords{},
		checks.PassMaxDays{},
		checks.AuditdEnabled{},
		checks.IPForwardingDisabled{},
		checks.SSHProtocol2{},
	}

	var results []model.CheckResult
	for _, chk := range checkList {
		r := chk.Run(ctx)
		results = append(results, model.CheckResult{
			CheckID:  r.CheckID,
			Title:    r.Title,
			Status:   r.Status,
			Severity: r.Severity,
			Evidence: r.Evidence,
		})
	}

	rep := model.Report{
		Agent: model.AgentMeta{
			Name:          cfg.Agent.Name,
			Version:       cfg.Agent.Version,
			Build:         cfg.Agent.Build,
			SchemaVersion: cfg.Agent.SchemaVersion,
			TsUTC:         time.Now().UTC().Format(time.RFC3339),
		},
		Host: model.HostMeta{
			HostID:    hostID,
			Hostname:  hostname,
			OS:        osName,
			OSVersion: osVer,
			Kernel:    kernel,
			IP:        ip,
		},
		Packages:   pkgs,
		CISResults: results,
	}

	// Always print JSON (good for demo)
	b, _ := json.MarshalIndent(rep, "", "  ")
	fmt.Println(string(b))

	if dryRun {
		return
	}

	client := transport.HTTPClient{
		URL:     cfg.Backend.URL,
		APIKey:  cfg.Backend.APIKey,
		Retries: cfg.Runtime.RetryCount,
		Backoff: time.Duration(cfg.Runtime.RetryBackoffMs) * time.Millisecond,
	}

	if err := client.PostJSON(ctx, rep); err != nil {
		fmt.Printf("ERROR sending report: %v\n", err)
	} else {
		fmt.Println("✅ Report sent to backend")
	}
}
