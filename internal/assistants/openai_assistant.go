package assistants

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"os"
	"runtime"
	"time"
)

func NewOpenAIAssistant(model string) (assistant Assistant, err error) {
	// TODO: Get more detailed information on the users OS like version etc.
	// TODO: Get the shell type this is being ran in

	llm, err := openai.New(openai.WithModel(model))
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

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	instructions := []string{
		"You are an AI assistant that provides operating system commands that accomplish the users request",
		"Your output is only the command. No explanation, no extra text, just the command",
		"Make sure all commands you output are have no formatting. Plain text only",
		"Use placeholders in angle brackets for unknown values the user should provide, example <subnet>",
		fmt.Sprintf("The user is running %s", runtime.GOOS),
	}

	if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && os.Getenv("SHELL") != "" {
		instructions = append(instructions, fmt.Sprintf("The user is using the %s shell", os.Getenv("SHELL")))
	}

	if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && os.Getenv("USER") != "" {
		instructions = append(instructions, fmt.Sprintf("The users name is %s. Use that for commands instead of ~ or $HOME", os.Getenv("USER")))
	}

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem,
			instructions...,
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

	s.Stop()
	return options, nil

}

func (a *OpenAIAssistant) Explain(command string) (response []string, err error) {

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	template := `
You follow the following template for your output. This output is for the example command 'ls -a /tmp/folderXXX'

ls → List directory contents.
-a → Include hidden files (those starting with .).
/tmp/folderXXX → Specifies the directory to list.
`

	instructions := []string{
		"You are an AI assistant that provides concise explanations of system commands",
		"Your responses are brief and concise",
		"When you see placeholders, words in angle brackets, always include instructions to replace with the real value",
		template,
		"You never include any sort of summary or lead-out, just the template",
		"Limit line length to 80 characters",
		fmt.Sprintf("They are running on %s", runtime.GOOS),
	}

	if runtime.GOOS == "linux" && os.Getenv("SHELL") != "" {
		instructions = append(instructions, fmt.Sprintf("The user is using the %s shell", os.Getenv("SHELL")))
	}

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem,
			instructions...,
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

	s.Stop()
	return options, nil

}
