package ip

import (
	"encoding/binary"
	"net"
	"net/http"
)

func ClientIP(request *http.Request) string {

	var ip string
	remoteAddr := request.RemoteAddr
	if ip = request.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = request.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

func ToInt(addr string) uint32 {

	ip := net.ParseIP(addr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}
