package ai_tools

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
	"os/exec"
)

// Ensure PathChecker implements the Tool interface
var _ tools.Tool = PathChecker{}

type PathChecker struct {
	CallbacksHandler callbacks.Handler
}

func (p PathChecker) Name() string {
	return "path_checker"
}

func (p PathChecker) Description() string {
	return "Checks if a given executable is in the system's PATH. Takes a single argument: the name of the " +
		"executable. Returns a message indicating whether the executable is found."
}

func (p PathChecker) Call(ctx context.Context, input string) (string, error) {
	if p.CallbacksHandler != nil {
		p.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// Check if the executable is in the PATH
	_, err := exec.LookPath(input)
	if err != nil {
		result := fmt.Sprintf("'%s' is not in the PATH.", input)
		if p.CallbacksHandler != nil {
			p.CallbacksHandler.HandleToolEnd(ctx, result)
		}
		return result, nil
	}

	result := fmt.Sprintf("'%s' is in the PATH.", input)
	if p.CallbacksHandler != nil {
		p.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}
