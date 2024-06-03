package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bgrewell/commander/internal/processors"
	"github.com/bgrewell/go-execute/v2"
	"github.com/bgrewell/usage"
	"github.com/joho/godotenv"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	version    string = "dev"
	buildDate  string = "dev"
	commitHash string = "dev"
	branch     string = "dev"
)

func installBinary() {
	installDir := "/opt/commander/bin/"
	linkLocation := "/usr/local/bin/commander"

	if runtime.GOOS == "windows" {
		installDir = "C:\\Program Files\\commander\\bin\\"
		linkLocation = "C:\\Windows\\commander"
	} else if runtime.GOOS == "darwin" {
		installDir = "/opt/commander/bin/"
		linkLocation = "/usr/local/bin/commander"
	}

	// 1. Ensure the directories exists
	err := ensureDir(installDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ensureDir(path.Dir(linkLocation))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Move the binary to the directory
	err = moveSelf(installDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. Create the symlink
	targetPath := filepath.Join(installDir, filepath.Base(os.Args[0]))
	err = createLink(targetPath, linkLocation)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. Display a success message
	fmt.Println("Installation successful.")
}

// checkPermissions checks if the program has permissions to write to the directory.
func checkPermissions(dir string) error {
	// Try to create a temporary file to check permissions.
	tempFile := filepath.Join(dir, "temp.txt")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("permission denied: must be run with elevated privileges")
	}
	file.Close()
	os.Remove(tempFile)
	return nil
}

// ensureDir ensures that the directory and all its parents are created.
func ensureDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

// moveSelf moves the running binary to the specified directory.
func moveSelf(installDir string) error {
	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	targetPath := filepath.Join(installDir, filepath.Base(executablePath))
	err = copyFile(executablePath, targetPath)
	if err != nil {
		return err
	}

	// Make the target file executable
	err = makeExecutable(targetPath)
	if err != nil {
		return err
	}

	// Remove the original executable after copying
	err = os.Remove(executablePath)
	if err != nil {
		return fmt.Errorf("failed to remove original executable: %v", err)
	}

	return nil
}

// makeExecutable ensures that the file is executable.
func makeExecutable(filePath string) error {
	err := os.Chmod(filePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to make file executable: %v", err)
	}
	return nil
}

// copyFile copies a file from src to dst. If src and dst files exist, and are the same, then return success. Otherise, attempt to copy the file contents.
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	if err = out.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %v", err)
	}

	return nil
}

// createLink creates a symbolic link at the specified location pointing to the executable.
func createLink(target, linkLocation string) error {
	// Check if the link already exists
	if _, err := os.Lstat(linkLocation); err == nil {
		// If it exists, remove it
		err := os.Remove(linkLocation)
		if err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	err := os.Symlink(target, linkLocation)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}
	return nil
}

// aliasExists checks if the alias already exists in the shell configuration file.
func aliasExists(shellConfigFile, aliasName string) (bool, error) {
	file, err := os.Open(shellConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // File does not exist, so alias can't exist either
		}
		return false, fmt.Errorf("failed to open shell configuration file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	aliasPrefix := fmt.Sprintf("alias %s=", aliasName)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, aliasPrefix) {
			return true, nil // Alias already exists
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to read shell configuration file: %v", err)
	}

	return false, nil // Alias does not exist
}

type Alias struct {
	Name    string
	Command string
}

// createAlias creates an alias in the user's shell configuration file if it doesn't already exist.
func createAlias(alias ...Alias) (shellcfg string, err error) {
	var shellConfigFile string

	for _, a := range alias {
		// Get the current user
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user: %v", err)
		}

		// Detect the user's shell
		shell := filepath.Base(os.Getenv("SHELL"))
		switch shell {
		case "bash":
			shellConfigFile = filepath.Join(usr.HomeDir, ".bashrc")
		case "zsh":
			shellConfigFile = filepath.Join(usr.HomeDir, ".zshrc")
		default:
			return "", fmt.Errorf("unsupported shell: %s", shell)
		}

		// Check if the alias already exists
		exists, err := aliasExists(shellConfigFile, a.Name)
		if err != nil {
			return "", err
		}
		if exists {
			fmt.Printf("Alias '%s' already exists in %s\n", a.Name, shellConfigFile)
			return shellConfigFile, nil
		}

		// Create the alias string
		aliasStr := fmt.Sprintf("alias %s='%s'\n", a.Name, a.Command)

		// Open the shell configuration file
		file, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return "", fmt.Errorf("failed to open shell configuration file: %v", err)
		}
		defer file.Close()

		// Write the alias to the file
		_, err = file.WriteString(aliasStr)
		if err != nil {
			return "", fmt.Errorf("failed to write alias to shell configuration file: %v", err)
		}

		fmt.Printf("Alias '%s' added to %s\n", a.Name, shellConfigFile)
	}

	return shellConfigFile, nil
}

func getExecutableLocation() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func main() {

	// Load a .env if it's present. If it's not, that's okay we will ignore that error
	locations := []string{".env"}
	binDir, err := getExecutableLocation()
	if err == nil {
		locations = append(locations, path.Join(binDir, ".env"))
	}
	_ = godotenv.Load(locations...)

	// Create a new usage to handle command line arguments
	sage := usage.NewUsage(
		usage.WithApplicationName("commander"),
		usage.WithApplicationVersion(version),
		usage.WithApplicationBuildDate(buildDate),
		usage.WithApplicationCommitHash(commitHash),
		usage.WithApplicationBranch(branch),
		usage.WithApplicationDescription("Commander is a command line tool that uses large language models like OpenAI's GPT-4 to generate commands based on a question. It can also explain the command and execute it. Use command execution with caution as you may execute a command you do not wish to run"),
	)

	// Add a group for query options
	query := sage.AddGroup(1, "Query", "Query options")

	// Add default options
	update := sage.AddBooleanOption("u", "update", false, "Check for updates to the application", "", nil)
	install := sage.AddBooleanOption("i", "install", false, "Install the application", "", nil)
	alias := sage.AddBooleanOption("a", "alias", false, "Create alias's 'c' (commander), 'ce' (commander -explain) 'cx' (commander -exec) for the application", "", nil)

	// Add query options
	clip := sage.AddBooleanOption("c", "clip", false, "Place the command in the clipboard", "", query)
	explain := sage.AddBooleanOption("e", "explain", false, "Provide an explanation of the output", "", query)
	exec := sage.AddBooleanOption("x", "exec", false, "Execute the returned command", "", query)

	// Add the question argument
	question := sage.AddArgument(1, "question", "The question to ask the assistant", "Question")

	// Parse the arguments
	parsed := sage.Parse()

	// Print the usage if the arguments were not parsed
	if !parsed {
		sage.PrintError(errors.New("Failed to parse arguments"))
	}

	if *install {
		installBinary()
		return
	}

	if *alias {
		cfgFile, err := createAlias(
			Alias{
				Name:    "c",
				Command: "commander",
			},
			Alias{
				Name:    "ce",
				Command: "commander -explain",
			},
			Alias{
				Name:    "cx",
				Command: "commander -exec",
			},
		)
		if err != nil {
			fmt.Printf("Failed to create alias: %v\n", err)
			return
		}

		fmt.Printf("Aliases created. Please run 'source %s' to apply the changes.\n", cfgFile)
		return
	}

	// Check if the update flag was passed
	if *update {

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

		if ver != "" {
			// Check to see if the version available is newer than the current version
			if versionIsNewer(ver, version) {
				log.Printf("Checking for updates...\n")
				log.Printf("Current version: %s\n", version)
				log.Printf("Version %s is available\n", ver)
				if err := updater.Update(); err != nil {
					log.Fatalf("Failed to update to version %s: %v\n", ver, err)
				}
				log.Printf("Updated to version %s\n", updater.Info.Version)
				os.Exit(0)
			}
		}
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

	// Print out the installation instructions if they are present
	if response.InstallInstructions != "" {
		fmt.Printf(response.InstallInstructions)
	}

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

func versionIsNewer(availableVersion string, currentVersion string) bool {

	if currentVersion == "dev" || strings.Contains(currentVersion, "dirty") {
		return false
	}

	// Split the available version
	available := splitVersion(availableVersion)

	// Split the current version
	current := splitVersion(currentVersion)

	// Compare the versions
	if available[0] > current[0] {
		return true
	} else if available[0] == current[0] {
		if available[1] > current[1] {
			return true
		} else if available[1] == current[1] {
			if available[2] > current[2] {
				return true
			}
		}
	}

	return false
}

func splitVersion(version string) []int {
	var major, minor, patch int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		log.Fatalf("Failed to split version: %v\n", err)
	}
	return []int{major, minor, patch}
}
