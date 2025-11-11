package ip_access

import (
	"net"
	"strings"
)

// Lists contains parsed whitelist/blacklist IPs and CIDR networks.
type Lists struct {
	WhitelistIPs  []net.IP
	WhitelistNets []*net.IPNet
	BlacklistIPs  []net.IP
	BlacklistNets []*net.IPNet
}

// Evaluate returns true if the ip is allowed according to the provided lists.
// Behavior:
// - If lists == nil the function returns true (no restrictions).
// - If any whitelist entries exist, the address must be present in the
//   whitelist (IPs or CIDRs) to be allowed.
// - Otherwise, if the address appears in the blacklist it is denied.
// - Default: allow.
func Evaluate(ip net.IP, lists *Lists) bool {
	if lists == nil {
		return true
	}

	// If whitelist present, only allow those entries
	if len(lists.WhitelistIPs) > 0 || len(lists.WhitelistNets) > 0 {
		return ipInList(ip, lists.WhitelistIPs, lists.WhitelistNets)
	}

	// If blacklisted, deny
	if ipInList(ip, lists.BlacklistIPs, lists.BlacklistNets) {
		return false
	}

	return true
}

func ipInList(ip net.IP, ips []net.IP, nets []*net.IPNet) bool {
	for _, a := range ips {
		if a.Equal(ip) {
			return true
		}
	}
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// ParseList parses a comma-separated list of IPs or CIDR ranges and returns
// the parsed IPs and networks. Invalid entries are ignored.
func ParseList(raw string) (ips []net.IP, nets []*net.IPNet) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, "/") {
			if _, network, err := net.ParseCIDR(p); err == nil {
				nets = append(nets, network)
				continue
			}
		}
		if ip := net.ParseIP(p); ip != nil {
			ips = append(ips, ip)
			continue
		}
		// ignore invalid entries
	}
	return ips, nets
}

// ParseLists is a convenience helper that parses both whitelist and blacklist
// raw strings and returns a fully populated Lists struct.
func ParseLists(whitelistRaw, blacklistRaw string) *Lists {
	wlIPs, wlNets := ParseList(whitelistRaw)
	blIPs, blNets := ParseList(blacklistRaw)
	return &Lists{WhitelistIPs: wlIPs, WhitelistNets: wlNets, BlacklistIPs: blIPs, BlacklistNets: blNets}
}
