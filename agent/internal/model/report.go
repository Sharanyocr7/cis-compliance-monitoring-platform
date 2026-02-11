package model

type AgentMeta struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	Build         string `json:"build"`
	SchemaVersion string `json:"schema_version"`
	TsUTC         string `json:"ts_utc"`
}

type HostMeta struct {
	HostID    string `json:"host_id"`
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	OSVersion string `json:"os_version"`
	Kernel    string `json:"kernel"`
	IP        string `json:"ip"`
}

type PackageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
}

type CheckResult struct {
	CheckID  string `json:"check_id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Evidence string `json:"evidence"`
}

type Report struct {
	Agent      AgentMeta     `json:"agent"`
	Host       HostMeta      `json:"host"`
	Packages   []PackageInfo `json:"packages"`
	CISResults []CheckResult `json:"cis_results"`
}
