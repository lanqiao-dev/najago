package i18n

import (
	"context"

	"golang.org/x/text/message"
)

func MessagePrinter(ctx context.Context) *message.Printer {
	val, ok := ctx.Value(messagePrinterKey).(*message.Printer)
	if !ok {
		return nil
	}
	return val
}

func MustMessagePrinter(ctx context.Context) *message.Printer {
	p, ok := ctx.Value(messagePrinterKey).(*message.Printer)
	if !ok {
		panic("could not find message printer from context")
	}
	return p
}
