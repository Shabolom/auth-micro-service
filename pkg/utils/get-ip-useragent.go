package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func IpUserAgentFromMetadata(ctx context.Context) (ip string, userAgent string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		agents := md.Get("user-agent")
		if len(agents) > 0 {
			userAgent = agents[0]
		}
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		ip = p.Addr.String()
	}

	return ip, userAgent
}
