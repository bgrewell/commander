package processors

import (
	"errors"
	"github.com/atotto/clipboard"
	"github.com/bgrewell/commander/internal"
	"github.com/bgrewell/commander/internal/assistants"
	"github.com/bgrewell/commander/internal/mutations"
	"github.com/fatih/color"
	"strings"
	"time"
)

func NewDefaultProcessor(options ...Option) Processor {
	p := &DefaultProcessor{
		yellow:       color.New(color.FgHiYellow),
		cyan:         color.New(color.FgHiCyan),
		gray:         color.New(color.FgWhite),
		white:        color.New(color.FgHiWhite),
		provider:     "openai",
		model:        "gpt-4o",
		useClipboard: false,
		useColor:     true,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

type DefaultProcessor struct {
	assistant    assistants.Assistant
	useColor     bool
	model        string
	provider     string
	useClipboard bool
	yellow       *color.Color
	cyan         *color.Color
	gray         *color.Color
	white        *color.Color
}

func (p *DefaultProcessor) SetColor(useColor bool) {
	p.useColor = useColor
}

func (p *DefaultProcessor) SetModel(model string) {
	p.model = model
}

func (p *DefaultProcessor) SetProvider(provider string) {
	p.provider = provider
}

func (p *DefaultProcessor) SetClipboard(useClipboard bool) {
	p.useClipboard = useClipboard
}

func (p *DefaultProcessor) Question(input string, explain bool) (response *internal.QuestionResponse, err error) {

	// Select the provider
	switch p.provider {
	case "openai":
		p.assistant, err = assistants.NewOpenAIAssistant(p.model)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown provider")
	}

	// Modify the question
	modifiedQ := mutations.Injection{}.Apply(input)

	// Create the response object
	response = &internal.QuestionResponse{
		Question:          input,
		AnnotatedQuestion: modifiedQ,
	}

	// Query the assistant
	var queryResponse []string
	var questionStart = time.Now()
	queryResponse, err = p.assistant.Query(modifiedQ)
	if err != nil {
		return nil, err
	}
	response.QuestionResponseTime = time.Since(questionStart)

	// Set the command and answer
	response.Command = queryResponse[0]
	response.Answer = p.cyan.Sprintf("%s\n", response.Command)

	// Copy to clipboard if requested
	if p.useClipboard {
		err = clipboard.WriteAll(response.Command)
		if err != nil {
			return nil, err
		}
	}

	// Explain if the flag is set
	if explain {
		explainStart := time.Now()
		// Get the explanation
		explanation, err := p.assistant.Explain(response.Command)
		if err != nil {
			return nil, err
		}
		response.ExplanationResponseTime = time.Since(explainStart)
		response.Explanation = explanation[0]
		lines := strings.Split(response.Explanation, "\n")
		response.Answer += p.cyan.Sprintf("\nCommand Explanation:\n")

		for _, line := range lines {
			if strings.Contains(line, "→") {
				parts := strings.Split(line, "→")
				response.Answer += p.yellow.Sprintf("  %s", parts[0])
				response.Answer += p.white.Sprintf("→")
				response.Answer += p.gray.Sprintf("%s\n", strings.Join(parts[1:], "→"))
			} else {
				response.Answer += p.gray.Sprintf("  %s\n", line)
			}
		}
	}

	return response, nil
}
