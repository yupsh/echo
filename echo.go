package echo

import (
	"context"
	"fmt"
	"io"
	"strings"

	localopt "github.com/yupsh/echo/opt"
	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
)

// Flags represents the configuration options for the echo command
type Flags = localopt.Flags

// Command implementation using StandardCommand abstraction
type command struct {
	yup.StandardCommand[Flags]
}

// Echo creates a new echo command with the given parameters
func Echo(parameters ...any) yup.Command {
	args := opt.Args[string, Flags](parameters...)
	return command{
		StandardCommand: yup.StandardCommand[Flags]{
			Positional: args.Positional,
			Flags:      args.Flags,
			Name:       "echo",
		},
	}
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	text := strings.Join(c.Positional, " ")

	// Apply escape processing if enabled
	if bool(c.Flags.Escape) {
		text = c.processEscapes(text)
	}

	// Output text
	if bool(c.Flags.NoNewline) {
		fmt.Fprint(stdout, text)
	} else {
		fmt.Fprintln(stdout, text)
	}

	return nil
}

func (c command) processEscapes(text string) string {
	// Simple escape processing - real implementation would be more complete
	replacements := map[string]string{
		"\\n":  "\n",
		"\\t":  "\t",
		"\\r":  "\r",
		"\\\\": "\\",
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	return text
}
