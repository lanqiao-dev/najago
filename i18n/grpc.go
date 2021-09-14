package i18n

import (
	"context"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	grpcLangMetaKey = "lang"
)

// UnaryServerInterceptor set the localizerfor every request.
func UnaryServerInterceptor(bundle *goi18n.Bundle) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var langs []string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			langs = md.Get(grpcLangMetaKey)

		}
		localizer := goi18n.NewLocalizer(bundle, langs...)
		ctx = context.WithValue(ctx, i18nLocalizerKey, localizer)
		return handler(ctx, req)
	}
}
