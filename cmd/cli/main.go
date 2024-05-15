package main

import (
	"errors"
	"fmt"
	execute "github.com/BGrewell/go-execute/v2"
	"github.com/atotto/clipboard"
	"github.com/bgrewell/commander/internal/assistants"
	"github.com/bgrewell/commander/internal/mutations"
	"github.com/bgrewell/usage"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

var (
	version    string = "0.0.1"
	buildDate  string = "debug"
	commitHash string = "debug"
	branch     string = "debug"
)

func main() {

	// Load a .env if it's present. If it's not, that's okay we will ignore that error
	_ = godotenv.Load()

	// Create a new sage to handle command line arguments
	sage := usage.NewUsage(
		usage.WithApplicationName("commander"),
		usage.WithApplicationVersion(version),
		usage.WithApplicationBuildDate(buildDate),
		usage.WithApplicationCommitHash(commitHash),
		usage.WithApplicationBranch(branch),
		usage.WithApplicationDescription("Commander is a command line tool that uses large language models like OpenAI's GPT-4 to generate commands based on a question. It can also explain the command and execute it. Use command execution with caution as you may execute a command you do not wish to run"),
	)

	// Add standard options
	explain := sage.AddBooleanOption("e", "explain", false, "Provide an explanation of the output", "", nil)
	exec := sage.AddBooleanOption("x", "exec", false, "Execute the returned command", "", nil)
	clip := sage.AddBooleanOption("c", "clip", false, "Place the command in the clipboard", "", nil)
	// Add the question argument
	question := sage.AddArgument(1, "question", "The question to ask the assistant", "Question")

	// Parse the arguments
	parsed := sage.Parse()

	// Print the usage if the arguments were not parsed
	if !parsed {
		sage.PrintError(errors.New("Failed to parse arguments"))
	}

	// Check if the question was provided
	if *question == "" {
		sage.PrintError(errors.New("You need to ask a question"))
	}

	// TODO: Move this elsewhere
	yellow := color.New(color.FgHiYellow)
	cyan := color.New(color.FgHiCyan)
	white := color.New(color.FgWhite)
	c := color.New(color.FgCyan)

	// Create a new assistant using the GPT-4 Turbo model
	assistant, err := assistants.NewOpenAIAssistant("gpt-4o")
	if err != nil {
		panic(err)
	}

	// TODO: TEMP
	modifiedQ := mutations.Injection{}.Apply(*question)

	// Query the assistant
	response, err := assistant.Query(modifiedQ)
	if err != nil {
		panic(err)
	}

	//cyan.Print("Command:")
	command := response[0]

	if *clip {
		err = clipboard.WriteAll(command)
		if err != nil {
			fmt.Println("Error copying to clipboard:", err)
			return
		}
	}

	// Print out the command
	c.Printf("%s\n", command)

	// Explain if the flag is set
	if *explain {
		// Get the explanation
		explanation, err := assistant.Explain(response[0])
		if err != nil {
			panic(err)
		}
		lines := strings.Split(explanation[0], "\n")
		cyan.Print("\nCommand Explanation:\n")
		for _, line := range lines {
			if strings.Contains(line, "→") {
				parts := strings.Split(line, "→")
				yellow.Printf("  %s→", parts[0])
				white.Printf("%s\n", strings.Join(parts[1:], "→"))
			} else {
				white.Printf("  %s\n", line)
			}
		}
	}

	// Execute if --exec was passed
	if *exec {
		exe := execute.NewExecutor()
		err := exe.ExecuteTTY(command)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}
