package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgnt = "grpcgateway-user-agent"
	userAgent           = "user-agent"
	xForwardedFor       = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgents := md.Get(grpcGatewayUserAgnt); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		} else if userAgents := md.Get(userAgent); len(userAgents) > 0 {
			// 일반 gRPC UserAgent 추출
			mtdt.UserAgent = userAgents[0]
		}

		// 일반적으로 프록시 환경에서 클라이언트 IP는 x-forwarded-for에 담김
		if clientIPs := md.Get(xForwardedFor); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if mtdt.ClientIP == "" {
		if p, ok := peer.FromContext(ctx); ok {
			mtdt.ClientIP = p.Addr.String()
		}
	}

	return mtdt
}
