package processors

import "github.com/bgrewell/commander/internal"

type Option func(processor Processor)

func WithModel(model string) Option {
	return func(p Processor) {
		p.SetModel(model)
	}
}

func WithProvider(provider string) Option {
	return func(p Processor) {
		p.SetProvider(provider)
	}
}

func WithColor(useColor bool) Option {
	return func(p Processor) {
		p.SetColor(useColor)
	}
}

func WithClipboard(useClipboard bool) Option {
	return func(p Processor) {
		p.SetClipboard(useClipboard)
	}
}

type Processor interface {
	Question(input string, explain bool) (response *internal.QuestionResponse, err error)
	SetColor(useColor bool)
	SetModel(model string)
	SetProvider(provider string)
	SetClipboard(useClipboard bool)
}
