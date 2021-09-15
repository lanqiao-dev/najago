package i18n

import (
	"context"
	"testing"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
)

func TestUnaryServerInterceptor(t *testing.T) {

	bundle, _ := NewBundle(language.English, "active.en.toml", "active.zh.toml")

	testCases := []struct {
		name       string
		inputMsg   *goi18n.Message
		langs      []string
		acceptLang string
		response   string
	}{
		{
			name:     "request with same lang and accpet lang",
			inputMsg: &goi18n.Message{ID: "invalid-captcha", Other: "Invalid captcha."},
			langs:    []string{language.Chinese.String(), language.Chinese.String()},
			response: "无效的验证码。",
		},
		{
			name:     "request with diff lang and accept lang",
			inputMsg: &goi18n.Message{ID: "invalid-captcha", Other: "Invalid captcha."},
			langs:    []string{language.Chinese.String(), language.English.String()},
			response: "无效的验证码。",
		},
		{
			name:       "request invalid lang and valid accpet lang",
			inputMsg:   &goi18n.Message{ID: "invalid-captcha", Other: "Invalid captcha."},
			langs:      []string{"invalid-lang", language.Chinese.String()},
			acceptLang: language.Chinese.String(),
			response:   "无效的验证码。",
		},
		{
			name:       "request valid lang and ivalid accpet lang",
			inputMsg:   &goi18n.Message{ID: "invalid-captcha", Other: "Invalid captcha."},
			langs:      []string{language.Chinese.String(), "invalid-accpet-lang"},
			acceptLang: "invalid-accpet-lang",
			response:   "无效的验证码。",
		},
		{
			name:     "request invalid lang and ivalid accpet lang",
			inputMsg: &goi18n.Message{ID: "invalid-captcha", Other: "Invalid captcha."},
			langs:    []string{"invalid-lang", "invalid-accpet-lang"},
			response: "Invalid captcha.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := UnaryServerInterceptor(bundle)
			md := metadata.MD{
				"lang": tc.langs,
			}
			ctx := metadata.NewIncomingContext(context.Background(), md)
			resp, _ := interceptor(ctx, tc.inputMsg, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				return L(ctx, tc.inputMsg), nil
			})

			if resp.(string) != tc.response {
				t.Errorf("want message: %#v, got message: %#v", tc.response, resp.(string))
			}
		})
	}
}
