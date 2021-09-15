package i18n

import (
	"context"
	"errors"

	"github.com/BurntSushi/toml"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var i18nLocalizerKey struct{}

// NewBundle create the translation bundle with default language tag and translation toml files.
func NewBundle(defaultLanguageTag language.Tag, tomls ...string) (*goi18n.Bundle, error) {
	bundle := goi18n.NewBundle(defaultLanguageTag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, file := range tomls {
		if _, err := bundle.LoadMessageFile(file); err != nil {
			return nil, err
		}
	}
	return bundle, nil
}

// MustLocalize translate message based on localizer in context, it will panic
// when not found localizer in context.
func MustLocalize(c context.Context, lc *goi18n.LocalizeConfig) string {
	return MustLocalizer(c).MustLocalize(lc)
}

// L localize the simple text without template data.
func L(c context.Context, message *goi18n.Message) string {
	return LocalizeText(c, message, nil)
}

// Localize localize the simple text without template data.
func Localize(c context.Context, message *goi18n.Message) string {
	return LocalizeText(c, message, nil)
}

// LT localize the template text with tmplate data.
func LT(c context.Context, message *goi18n.Message, tplData map[string]interface{}) string {
	return localizeText(c, message, tplData)
}

// LocalizeText format the template text with tmplate data.
func LocalizeText(c context.Context, message *goi18n.Message, tplData map[string]interface{}) string {
	return localizeText(c, message, tplData)
}

func localizeText(c context.Context, message *goi18n.Message, tplData map[string]interface{}) string {
	if localizer, ok := Localizer(c); ok {
		return localizer.MustLocalize(&goi18n.LocalizeConfig{
			DefaultMessage: message,
			TemplateData:   tplData,
		})
	}
	return localizeInternalMessage(message, tplData)
}

// MustLocalizer get localizer from context otherwise will panic.
func MustLocalizer(c context.Context) *goi18n.Localizer {
	localizer, ok := Localizer(c)
	if !ok {
		panic(errors.New("context has no localizer"))
	}
	return localizer
}

// Localizer get the localizer
func Localizer(c context.Context) (*goi18n.Localizer, bool) {
	v := c.Value(i18nLocalizerKey)
	if l, b := v.(*goi18n.Localizer); b {
		return l, true
	}
	return nil, false
}

func localizeInternalMessage(message *goi18n.Message, args map[string]interface{}) string {
	if args == nil {
		return message.Other
	}
	tpl := goi18n.NewMessageTemplate(message)
	msg, err := tpl.Execute("other", args, nil)
	if err != nil {
		panic(err)
	}
	return msg
}
