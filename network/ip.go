package network

import (
	"net"
	"net/http"
	"strings"
)

func GetRemoteIP(r *http.Request, enableForwarding bool) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !enableForwarding {
		return host
	}

	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			host = ip.String()
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			host = ip.String()
		}
	}

	return host
}
