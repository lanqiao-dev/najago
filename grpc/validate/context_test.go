package validate

import (
	"context"
	"testing"
)

func TestContextWithLocale(t *testing.T) {
	want := "zh_cn"
	ctx := NewContextWithLocale(context.Background(), want)
	if ctx == nil {
		t.Errorf("want not nil, got nil")
	}

	if got := localeFromContext(ctx); want != got {
		t.Errorf("want %s, got %s", want, got)
	}
	if got := localeFromContext(context.Background()); got != "" {
		t.Errorf("want empty string, got %s", got)
	}
}
