package helper

import (
	"context"
	"fmt"
	"strings"

	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
	Lang      string
}

func ExtractMetadata(ctx context.Context) *Metadata {
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

		if lang := md.Get(consts.AcceptLanguage); len(lang) > 0 {
			mtdt.Lang = lang[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}
	fmt.Println(mtdt)

	return mtdt
}

func CustomMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "accept-language":
		return key, true
	case "authorization":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
