package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func IpUserAgentFromMetadata(ctx context.Context) (ip string, userAgent string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ""
	}

	ips := md.Get("client-ip")
	if len(ips) > 0 {
		ip = ips[0]
	}

	agents := md.Get("user-agent")
	if len(agents) > 0 {
		userAgent = agents[0]
	}

	return ip, userAgent
}
