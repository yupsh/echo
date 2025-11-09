package command

import (
	`context`
	`io`

	gloo "github.com/gloo-foo/framework"
)

// Command implementation using StandardCommand abstraction
type command gloo.Inputs[string, flags]

func Echo(parameters ...any) gloo.Command {
	return command(gloo.Initialize[string, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Join positional arguments with space
		output := ""
		for i, arg := range p.Positional {
			if i > 0 {
				output += " "
			}

			// Handle escape sequences if enabled
			if bool(p.Flags.Escape) {
				arg = processEscapes(arg)
			}
			output += arg
		}

		// Add newline unless NoNewline flag is set
		if !bool(p.Flags.NoNewline) {
			output += "\n"
		}

		_, err := io.WriteString(stdout, output)
		return err
	}
}

// processEscapes handles backslash escape sequences
func processEscapes(s string) string {
	result := ""
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result += "\n"
				i += 2
			case 't':
				result += "\t"
				i += 2
			case 'r':
				result += "\r"
				i += 2
			case '\\':
				result += "\\"
				i += 2
			case 'a':
				result += "\a"
				i += 2
			case 'b':
				result += "\b"
				i += 2
			case 'f':
				result += "\f"
				i += 2
			case 'v':
				result += "\v"
				i += 2
			default:
				result += string(s[i])
				i++
			}
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
