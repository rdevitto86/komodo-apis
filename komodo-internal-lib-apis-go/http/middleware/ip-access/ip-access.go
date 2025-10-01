package ipfilter

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"komodo-internal-lib-apis-go/config"
	httpUtils "komodo-internal-lib-apis-go/http/utils"
	logger "komodo-internal-lib-apis-go/logger/runtime"
)

// IPAccessMiddleware enforces allow/deny rules based on client IP.
//
// Configuration:
// - `IP_WHITELIST`: comma-separated list of IPs or CIDR ranges. If set, only these addresses are allowed.
// - `IP_BLACKLIST`: comma-separated list of IPs or CIDR ranges. If set and whitelist is empty, listed addresses are denied.
//
// Recommendation: Add these variables to your environment configuration (e.g. `env/.env.dev` and deployment
// system) so they can be centrally managed. For more advanced setups, move them into a config package or
// use a secrets/config service and inject values at startup.

var (
	ipOnce       sync.Once
	wlNets       []*net.IPNet
	wlIPs        []net.IP
	blNets       []*net.IPNet
	blIPs        []net.IP
)

func IPAccessMiddleware(next http.Handler) http.Handler {
	// lazy-parse env config once
	ipOnce.Do(func() {
		parseAndLoadList(config.GetConfigValue("IP_WHITELIST"), &wlIPs, &wlNets)
		parseAndLoadList(config.GetConfigValue("IP_BLACKLIST"), &blIPs, &blNets)
	})

	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		client := httpUtils.GetClientKey(req)
		if client == "" {
			logger.Error("unable to determine client IP", req)
			http.Error(wtr, "Forbidden", http.StatusForbidden)
			return
		}

		ip := net.ParseIP(client)
		if ip == nil {
			// Try to trim potential port if present
			host, _, err := net.SplitHostPort(client)
			if err == nil {
				ip = net.ParseIP(host)
			}
		}
		if ip == nil {
			logger.Error("invalid client IP: "+client, req)
			http.Error(wtr, "Forbidden", http.StatusForbidden)
			return
		}

		allowed := evaluateIP(ip)
		if !allowed {
			logger.Error("access denied for client ip: "+client, req)
			http.Error(wtr, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(wtr, req)
	})
}

// parseAndLoadList parses comma-separated entries (CIDR or IP) into slices of IPs and IPNets
func parseAndLoadList(raw string, ips *[]net.IP, nets *[]*net.IPNet) {
	if strings.TrimSpace(raw) == "" {
		return
	}
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, "/") {
			if _, network, err := net.ParseCIDR(p); err == nil {
				*nets = append(*nets, network)
				continue
			}
		}
		if ip := net.ParseIP(p); ip != nil {
			*ips = append(*ips, ip)
			continue
		}
		// ignore invalid entries but don't crash; log could be added here
	}
}

// evaluateIP returns true if the given IP is allowed based on whitelist/blacklist rules.
func evaluateIP(ip net.IP) bool {
	// If whitelist present, only allow those
	if len(wlIPs) > 0 || len(wlNets) > 0 {
		return ipInList(ip, wlIPs, wlNets)
	}

	// No whitelist: if blacklisted, deny
	if ipInList(ip, blIPs, blNets) {
		return false
	}

	// default allow
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
