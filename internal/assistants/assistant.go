package assistants

import "github.com/tmc/langchaingo/tools"

type Option func(assistant Assistant)

func WithModel(model string) Option {
	return func(a Assistant) {
		a.SetModel(model)
	}
}

func WithTools(tools []tools.Tool) Option {
	return func(a Assistant) {
		a.SetTools(tools)
	}
}

type Assistant interface {
	Query(message string) (response []string, err error)
	Explain(command string) (response []string, err error)
	SetModel(model string)
	SetTools(tools []tools.Tool)
}
