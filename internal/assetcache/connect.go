package assetcache

import (
	"context"
	"net"
	"net/netip"
	"strconv"
)

// tcpConnect checks the first reported address that is still assigned locally.
// An unavailable endpoint or expired status budget leaves the check unattempted.
func tcpConnect(ctx context.Context, status Status) (success, checked bool) {
	if !status.Active || status.Port == nil || *status.Port < 1 || *status.Port > 65535 {
		return false, false
	}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return false, false
	}

	for _, reported := range status.PrivateAddresses {
		ip, err := netip.ParseAddr(reported)
		if err != nil || (ip.IsLinkLocalUnicast() && ip.Zone() == "") {
			continue
		}
		for _, address := range addresses {
			prefix, err := netip.ParsePrefix(address.String())
			if err != nil || prefix.Addr().Unmap() != ip.WithZone("").Unmap() {
				continue
			}
			if ctx.Err() != nil {
				return false, false
			}
			var dialer net.Dialer
			conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ip.String(), strconv.Itoa(*status.Port)))
			if err != nil {
				return false, true
			}
			_ = conn.Close()
			return true, true
		}
	}
	return false, false
}
