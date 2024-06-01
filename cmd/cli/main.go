package main

import (
	"errors"
	"fmt"
	"github.com/bgrewell/commander/internal/processors"
	"github.com/bgrewell/go-execute/v2"
	"github.com/bgrewell/usage"
	"github.com/joho/godotenv"
	"log"
	"os"
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

	// Create a new usage to handle command line arguments
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

	p := processors.NewDefaultProcessor(
		processors.WithModel("gpt-4o"),
		processors.WithProvider("openai"),
		processors.WithColor(true),
		processors.WithClipboard(*clip),
	)

	response, err := p.Question(*question, *explain)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// Print out the formatted response
	fmt.Printf(response.Answer)

	// Execute if --exec was passed
	if *exec {
		exe := execute.NewExecutor(
			execute.WithEnvironment(os.Environ()),
			// TODO: this could be problematic as users may end up unknowingly 'stuck' inside another shell. Executing
			//   outside a shell could also have problems if the command is a shell built-in or alias or otherwise has
			//   a dependency on the shell environment. This is a good starting point but should be improved in the
			//   future.
			execute.WithDefaultShell(),
		)
		err := exe.ExecuteTTY(response.Command)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}
