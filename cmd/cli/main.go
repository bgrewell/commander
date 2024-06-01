package main

import (
	"errors"
	"fmt"
	"github.com/bgrewell/commander/internal/processors"
	"github.com/bgrewell/go-execute/v2"
	"github.com/bgrewell/usage"
	"github.com/joho/godotenv"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"log"
	"os"
)

var (
	version    string = "dev"
	buildDate  string = "dev"
	commitHash string = "dev"
	branch     string = "dev"
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
	update := sage.AddBooleanOption("u", "update", false, "Check for updates to the application", "", nil)
	// Add the question argument
	question := sage.AddArgument(1, "question", "The question to ask the assistant", "Question")

	// Parse the arguments
	parsed := sage.Parse()

	// Print the usage if the arguments were not parsed
	if !parsed {
		sage.PrintError(errors.New("Failed to parse arguments"))
	}

	// Check if the update flag was passed
	if *update {

		log.Printf("Checking for updates...\n")
		log.Printf("Current version: %s\n", version)

		// TODO: This is temporary while commander is tested out. It will be replaced with a public update URL
		//       once the beta version of commander has been released
		updateURL := os.Getenv("COMMANDER_UPDATE_URL")

		var updater = &selfupdate.Updater{
			CurrentVersion: version,
			ApiURL:         updateURL,
			BinURL:         updateURL,
			DiffURL:        updateURL,
			Dir:            "update/",
			CmdName:        "commander",
			ForceCheck:     true,
		}

		ver, err := updater.UpdateAvailable()
		if err != nil {
			log.Fatalf("Failed to get update information: %v\n", err)
		}
		fmt.Printf("Version %s is available\n", ver)

		// If the update was successful, print the new version and exit
		fmt.Printf("Updated to version %s\n", updater.Info.Version)
		os.Exit(0)
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
