package ipaccess

import (
	"net"
	"net/http"
	"sync"

	"komodo-forge-apis-go/config"
	httpErr "komodo-forge-apis-go/http/errors"
	httpReq "komodo-forge-apis-go/http/request"
	ipsvc "komodo-forge-apis-go/http/services/ip_access"
	logger "komodo-forge-apis-go/logging/runtime"
)

// IPAccessMiddleware enforces allow/deny rules based on client IP.
//
// Configuration:
// - `IP_WHITELIST`: comma-separated list of IPs or CIDR ranges. If set, only these addresses are allowed.
// - `IP_BLACKLIST`: comma-separated list of IPs or CIDR ranges. If set and whitelist is empty, listed addresses are denied.
//

var (
	ipOnce sync.Once
	lists ipsvc.Lists
)

// Enforces allow/deny rules based on client IP.
func IPAccessMiddleware(next http.Handler) http.Handler {
	// lazy-parse env config once
	ipOnce.Do(func() {
		wlIPs, wlNets := ipsvc.ParseList(config.GetConfigValue("IP_WHITELIST"))
		blIPs, blNets := ipsvc.ParseList(config.GetConfigValue("IP_BLACKLIST"))
		lists = ipsvc.Lists{WhitelistIPs: wlIPs, WhitelistNets: wlNets, BlacklistIPs: blIPs, BlacklistNets: blNets}
	})

	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		client := httpReq.GetClientKey(req)
		if client == "" {
			logger.Error("unable to determine client IP", req)
			httpErr.SendError(wtr, req, httpErr.Global.Forbidden, httpErr.WithDetail("unable to determine client IP"))
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
			logger.Error("invalid client IP: " + client, req)
			httpErr.SendError(wtr, req, httpErr.Global.Forbidden, httpErr.WithDetail("invalid client IP"))
			return
		}

		allowed := ipsvc.Evaluate(ip, &lists)
		if !allowed {
			logger.Error("access denied for client ip: " + client, req)
			httpErr.SendError(wtr, req, httpErr.Global.Forbidden, httpErr.WithDetail("access denied for client IP"))
			return
		}

		next.ServeHTTP(wtr, req)
	})
}
