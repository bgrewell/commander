package mutations

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"os/user"
	"runtime"
	"strings"
)

type Injection struct {
}

func (i Injection) Apply(input string) (output string) {
	os := runtime.GOOS
	shell := getShell()
	user := getUser()
	homedir := getHomedir()

	details := fmt.Sprintf("os: %s\nshell: %s\nuser: %s\nhomedir: %s", os, shell, user, homedir)
	return fmt.Sprintf("%s\n\nAdditional Context: %s\n", input, details)
}

func getUser() (username string) {
	currentUser, err := user.Current()
	if err != nil {
		return "unknown"
	}

	return currentUser.Username
}

func getHomedir() (homedir string) {
	currentUser, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return currentUser.HomeDir
}

func getShell() (shell string) {
	switch runtime.GOOS {
	case "linux":
	case "darwin":
		if shell := os.Getenv("SHELL"); shell != "" {
			return shell
		}
	case "windows":
		return detectShell()
	}
	return "unknown"
}

func detectShell() string {
	if _, exists := os.LookupEnv("PSModulePath"); exists {
		return "PowerShell"
	}

	parentProcessName := getParentProcessName()
	if strings.Contains(parentProcessName, "powershell") {
		return "PowerShell"
	}

	return "cmd.exe"
}

func getParentProcessName() string {
	pid := int32(os.Getpid())
	p, err := process.NewProcess(pid)
	if err != nil {
		return ""
	}

	pp, err := p.Parent()
	if err != nil {
		return ""
	}

	name, err := pp.Name()
	if err != nil {
		return ""
	}

	return name
}
