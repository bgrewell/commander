package assistants

import (
	"context"
	"github.com/briandowns/spinner"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
	"time"
)

func NewOpenAIAssistant(options ...Option) (assistant Assistant, err error) {

	// Create a new OpenAIAssistant
	a := &OpenAIAssistant{
		model: "gpt-4o",
		tools: make([]tools.Tool, 0),
	}

	// Apply the options
	for _, option := range options {
		option(a)
	}

	// Setup the LLM
	a.llm, err = openai.New(
		openai.WithModel(a.model),
	)
	if err != nil {
		return nil, err
	}

	commandSystemPrompt := "You are an AI assistant designed to generate operating system commands directly corresponding to " +
		"user requests. Adhere to the following guidelines:\n\n" +
		"Command Only: Provide only the command itself. Exclude any explanations, descriptions, or supplementary text.\n" +
		"Plain Text Format: Ensure all commands are presented in plain text with no special formatting or embellishments.\n" +
		"Use of Placeholders: For any values that the user needs to input, use placeholders enclosed in angle brackets " +
		"(e.g., <filename>). Be explicit in these placeholders to indicate what type of value the user should provide.\n" +
		"Response Format: Your responses should be concise and limited to a single command line per user request. " +
		"This format should be strictly maintained to ensure clarity and ease of use for the user.\n" +
		"Your primary goal is to respond with accurate and applicable commands that the user can directly execute, " +
		"maintaining a straightforward and utilitarian approach in all responses.\n"

	a.commandAgent = agents.NewOpenAIFunctionsAgent(
		a.llm,
		a.tools,
		agents.NewOpenAIOption().WithSystemMessage(commandSystemPrompt),
		agents.NewOpenAIOption().WithExtraMessages([]prompts.MessageFormatter{}),
	)
	a.commandExecutor = agents.NewExecutor(
		a.commandAgent,
		a.tools,
		agents.NewOpenAIOption().WithSystemMessage(commandSystemPrompt),
	)

	explainSystemPrompt := "You are an AI assistant tasked with providing concise explanations of system commands. " +
		"Follow these guidelines for each command explanation:\n\n" +
		"Conciseness: Your responses must be brief and directly related to the command.\n" +
		"Exclusion of Summaries: Do not include any sort of summary or additional commentary beyond the template explanation.\n" +
		"Placeholders: For any placeholders, such as words in angle brackets (e.g., <placeholder>), include explicit instructions to replace them with actual values when applicable.\n" +
		"Brevity in Formatting: Ensure each line does not exceed 80 characters.\n" +
		"Output Template: Use the specific template below for your explanations. For the command ls -a /tmp/folderXXX, the output should look like this:\n\n" +
		"ls → List directory contents.\n" +
		"-a → Include hidden files (those starting with '.').\n" +
		"/tmp/folderXXX → Specifies the directory to list.\n"

	a.explainAgent = agents.NewOpenAIFunctionsAgent(
		a.llm,
		a.tools,
		agents.NewOpenAIOption().WithSystemMessage(explainSystemPrompt),
		agents.NewOpenAIOption().WithExtraMessages([]prompts.MessageFormatter{}),
	)
	a.explainExecutor = agents.NewExecutor(
		a.explainAgent,
		a.tools,
		agents.NewOpenAIOption().WithSystemMessage(explainSystemPrompt),
	)

	return a, nil
}

type OpenAIAssistant struct {
	commandAgent    *agents.OpenAIFunctionsAgent
	explainAgent    *agents.OpenAIFunctionsAgent
	commandExecutor agents.Executor
	explainExecutor agents.Executor
	llm             *openai.LLM
	model           string
	tools           []tools.Tool
}

func (a *OpenAIAssistant) SetModel(model string) {
	a.model = model
}

func (a *OpenAIAssistant) SetTools(tools []tools.Tool) {
	a.tools = tools
}

func (a *OpenAIAssistant) Query(message string) (response []string, err error) {

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("fgHiGreen")
	s.Start()

	result, err := chains.Run(context.Background(), a.commandExecutor, message)
	if err != nil {
		return nil, err
	}

	s.Stop()

	return []string{result}, nil
}

func (a *OpenAIAssistant) Explain(command string) (response []string, err error) {

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("fgHiGreen")
	s.Start()

	result, err := chains.Run(context.Background(), a.explainExecutor, command)
	if err != nil {
		return nil, err
	}

	s.Stop()
	return []string{result}, nil

}
