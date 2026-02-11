package collector

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"strings"
)

func Hostname() string {
	h, _ := os.Hostname()
	return h
}

func OSInfo() (string, string) {
	b, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", ""
	}
	lines := strings.Split(string(b), "\n")
	var id, ver string
	for _, ln := range lines {
		if strings.HasPrefix(ln, "ID=") {
			id = strings.Trim(strings.TrimPrefix(ln, "ID="), "\"")
		}
		if strings.HasPrefix(ln, "VERSION_ID=") {
			ver = strings.Trim(strings.TrimPrefix(ln, "VERSION_ID="), "\"")
		}
	}
	return id, ver
}

func PrimaryIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			ipNet, ok := a.(*net.IPNet)
			if !ok || ipNet.IP == nil {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return ""
}

func StableHostID(hostname, osName, osVer, kernel string) string {
	h := sha256.Sum256([]byte(hostname + "|" + osName + "|" + osVer + "|" + kernel))
	return hex.EncodeToString(h[:16])
}
