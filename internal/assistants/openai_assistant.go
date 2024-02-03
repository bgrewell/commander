package assistants

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"runtime"
)

func NewOpenAIAssistant() (assistant Assistant, err error) {
	// TODO: Get more detailed information on the users OS like version etc.
	// TODO: Get the shell type this is being ran in

	llm, err := openai.New(openai.WithModel(""))
	if err != nil {
		return nil, err
	}

	return &OpenAIAssistant{
		llm: llm,
	}, nil
}

type OpenAIAssistant struct {
	llm *openai.LLM
}

func (a *OpenAIAssistant) Query(message string) (response []string, err error) {

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem,
			"You are an AI assistant that provides operating system commands that accomplish the users request",
			"Your output is only the command. No explanation, no extra text, just the command",
			"Make sure all commands you output are have no formatting. Plain text only",
			fmt.Sprintf("They are running on %s", runtime.GOOS),
		),
		llms.TextParts(schema.ChatMessageTypeHuman, message),
	}

	completion, err := a.llm.GenerateContent(context.Background(), content,
		llms.WithStreamingFunc(
			func(ctx context.Context, chunk []byte) error {
				//fmt.Printf("chunk: %s [%d bytes]\n", string(chunk), len(chunk))
				return nil
			},
		),
		llms.WithTemperature(0.1),
	)
	if err != nil {
		return nil, err
	}

	options := make([]string, 0)
	for _, choice := range completion.Choices {
		options = append(options, choice.Content)
	}

	return options, nil

}

func (a *OpenAIAssistant) Explain(command string) (response []string, err error) {

	template := `
You follow the following template for your output. This output is for the example command 'ls -lah /tmp/folderXXX'

ls: List directory contents.
-l: Long format, showing detailed information.
-a: Include hidden files (those starting with .).
-h: Human-readable sizes (e.g., 1K, 234M, 2G).
/tmp/folderXXX: Specifies the directory to list.
`

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem,
			"You are an AI assistant that provides concise explanations of system commands",
			"Your responses are brief and concise",
			template,
			"You never include any sort of summary or lead-out, just the template",
			fmt.Sprintf("They are running on %s", runtime.GOOS),
		),
		llms.TextParts(schema.ChatMessageTypeHuman, command),
	}

	completion, err := a.llm.GenerateContent(context.Background(), content,
		llms.WithStreamingFunc(
			func(ctx context.Context, chunk []byte) error {
				//fmt.Printf("chunk: %s [%d bytes]\n", string(chunk), len(chunk))
				return nil
			},
		),
		llms.WithTemperature(0.1),
	)
	if err != nil {
		return nil, err
	}

	options := make([]string, 0)
	for _, choice := range completion.Choices {
		options = append(options, choice.Content)
	}

	return options, nil

}
