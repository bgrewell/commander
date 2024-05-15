package assistants

import (
	"context"
	"github.com/briandowns/spinner"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
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

	systemPrompt := []string{
		"You are an AI assistant designed to generate operating system commands directly corresponding to user requests. Adhere to the following guidelines:\n\n",
		"Command Only: Provide only the command itself. Exclude any explanations, descriptions, or supplementary text.\n",
		"Plain Text Format: Ensure all commands are presented in plain text with no special formatting or embellishments.\n",
		"Use of Placeholders: For any values that the user needs to input, use placeholders enclosed in angle brackets (e.g., <filename>). Be explicit in these placeholders to indicate what type of value the user should provide.\n",
		"Response Format: Your responses should be concise and limited to a single command line per user request. This format should be strictly maintained to ensure clarity and ease of use for the user.\n",
		"Your primary goal is to respond with accurate and applicable commands that the user can directly execute, maintaining a straightforward and utilitarian approach in all responses.\n",
	}

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem,
			systemPrompt...,
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

	systemPrompt := "You are an AI assistant tasked with providing concise explanations of system commands. Follow these guidelines for each command explanation:\n\nConciseness: Your responses must be brief and directly related to the command.\nPlaceholders: For any placeholders, such as words in angle brackets (e.g., <placeholder>), include explicit instructions to replace them with actual values when applicable.\nOutput Template: Use the specific template below for your explanations. For the command ls -a /tmp/folderXXX, the output should look like this:\nls → List directory contents.\n-a → Include hidden files (those starting with '.').\n/tmp/folderXXX → Specifies the directory to list.\nBrevity in Formatting: Ensure each line does not exceed 80 characters.\nExclusion of Summaries: Do not include any sort of summary or additional commentary beyond the template explanation."

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem,
			systemPrompt,
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
