package ipaccess

import (
	"net"
	"net/http"
	"sync"

	"komodo-internal-lib-apis-go/common/errors"
	"komodo-internal-lib-apis-go/config"
	utils "komodo-internal-lib-apis-go/http/utils/http"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	ipsvc "komodo-internal-lib-apis-go/security/ip_access"
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
		client := utils.GetClientKey(req)
		if client == "" {
			logger.Error("unable to determine client IP", req)
			errors.WriteErrorResponse(wtr, req, http.StatusForbidden, errors.ERR_ACCESS_DENIED, "forbidden")
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
			errors.WriteErrorResponse(wtr, req, http.StatusForbidden, errors.ERR_ACCESS_DENIED, "forbidden")
			return
		}

		allowed := ipsvc.Evaluate(ip, &lists)
		if !allowed {
			logger.Error("access denied for client ip: "+client, req)
			errors.WriteErrorResponse(wtr, req, http.StatusForbidden, errors.ERR_ACCESS_DENIED, "forbidden")
			return
		}

		next.ServeHTTP(wtr, req)
	})
}
