package validate

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	english "github.com/go-playground/validator/v10/translations/en"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Validator validates struct with translated messages.
type Validator struct {
	validate            *validator.Validate
	universalTranslator *ut.UniversalTranslator
	fallbackLocale      string
}

type validatorOptions struct {
	fallbackTranslator       locales.Translator
	localeTranslators        []locales.Translator
	translationRegistrations []defaultTranslationRegistration
}

type defaultTranslationRegistration struct {
	prefix                         string
	registerDefaultTranslationFunc registerTranslationFunc
}

// ValidatorOptions is the options for NewValidator function.
type ValidatorOptions func(o *validatorOptions)

type registerTranslationFunc func(v *validator.Validate, ut ut.Translator) error

// WithTranslators returns the option for locale translators.
func WithTranslators(fallbackTranslator locales.Translator, localeTranslators ...locales.Translator) ValidatorOptions {
	return func(o *validatorOptions) {
		o.fallbackTranslator = fallbackTranslator
		o.localeTranslators = localeTranslators
	}
}

// WithTranslators returns the option to set translation functions for locales matching with localePrefix.
func WithRegisterDefaultTranslationFunc(localePrefix string, registerDefaultTranslationFunc registerTranslationFunc) ValidatorOptions {
	return func(o *validatorOptions) {
		o.translationRegistrations = append(o.translationRegistrations, defaultTranslationRegistration{
			prefix:                         localePrefix,
			registerDefaultTranslationFunc: registerDefaultTranslationFunc,
		})
	}
}

var defaultOptions = validatorOptions{
	fallbackTranslator:       en.New(),
	localeTranslators:        []locales.Translator{},
	translationRegistrations: []defaultTranslationRegistration{},
}

// NewValidator returns Validator.
func NewValidator(vOpts ...ValidatorOptions) (*Validator, error) {
	opt := defaultOptions
	for _, vo := range vOpts {
		vo(&opt)
	}
	if len(opt.translationRegistrations) == 0 {
		opt.translationRegistrations = []defaultTranslationRegistration{
			{
				prefix:                         en.New().Locale(),
				registerDefaultTranslationFunc: english.RegisterDefaultTranslations,
			},
		}
	}

	allLocaleTranslators := append(opt.localeTranslators, opt.fallbackTranslator)
	universalTranslator := ut.New(opt.fallbackTranslator, allLocaleTranslators...)
	v := &Validator{
		validate:            validator.New(),
		universalTranslator: universalTranslator,
		fallbackLocale:      opt.fallbackTranslator.Locale(),
	}

	if err := v.registerDefaultTranslations(allLocaleTranslators, opt.translationRegistrations); err != nil {
		return nil, err
	}
	return v, nil
}

func (v *Validator) registerDefaultTranslations(locales []locales.Translator, translationFuncs []defaultTranslationRegistration) error {
	for _, locale := range locales {
		localeTranslator, found := v.universalTranslator.GetTranslator(locale.Locale())
		if !found {
			continue
		}

		for _, tmf := range translationFuncs {
			if !strings.HasPrefix(locale.Locale(), tmf.prefix) {
				continue
			}
			if err := tmf.registerDefaultTranslationFunc(v.validate, localeTranslator); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Validator) getTranslator(locale string) ut.Translator {
	if locale == "" {
		return v.universalTranslator.GetFallback()
	}

	t, found := v.universalTranslator.GetTranslator(locale)
	if !found {
		grpclog.Warningf("failed to get translator for locale=%s, use a fallback translator", locale)
	}
	return t
}

// RegisterValidationCtx applies go-playground/validator/v10/Validate.RegisterValidationCtx.
func (v *Validator) RegisterValidationCtx(tag string, fn validator.FuncCtx) error {
	return v.validate.RegisterValidationCtx(tag, fn)
}

// ValidateGRPCRequest validates the request and returns the grpc status if it's error.
func (v *Validator) ValidateGRPCRequest(ctx context.Context, req interface{}) (*status.Status, error) {
	if req == nil {
		return nil, nil
	}

	err := v.validate.StructCtx(ctx, req)
	if err == nil {
		return nil, nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil, err
	}
	if len(errs) == 0 {
		return nil, nil
	}

	locale := localeFromContext(ctx)
	if locale == "" {
		locale = v.fallbackLocale
	}
	t := v.getTranslator(locale)
	br := ConvertValidationErrors(errs.Translate(t))
	if br == nil {
		panic(errors.New("failed to validate request but cannot convert validation errors"))
	}

	// use the first field error as the main error message.
	st := status.New(codes.InvalidArgument, br.FieldViolations[0].Description)
	dstSt, err := st.WithDetails(br)
	if err != nil {
		panic(err)
	}
	return dstSt, nil
}
