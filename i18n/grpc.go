package i18n

import (
	"context"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var messagePrinterKey struct{}

// UnaryServerInterceptor set the message printer for every request.
func UnaryServerInterceptor(matcher language.Matcher) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var values []string
		var lang string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values = md.Get("accept-language")
			lang = values[0]
		} else {
			lang = "'"
		}
		t, _, _ := language.ParseAcceptLanguage(lang)
		tag, _, _ := matcher.Match(t...)

		p := message.NewPrinter(tag)
		ctx = context.WithValue(ctx, messagePrinterKey, p)

		return handler(ctx, req)
	}
}
