package gapi

import (
	"context"

	consts "github.com/MyProject273/Join_Love/pkg/const"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(consts.GrpcGatewayUserAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(consts.UserAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if clientIPs := md.Get(consts.XForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
