package command_test

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/echo"
)

// ==============================================================================
// Test Basic Functionality
// ==============================================================================

func TestEcho_NoArguments(t *testing.T) {
	result := run.Command(command.Echo()).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{""})
}

func TestEcho_SingleArgument(t *testing.T) {
	result := run.Command(command.Echo("hello")).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"hello"})
}

func TestEcho_MultipleArguments(t *testing.T) {
	result := run.Command(command.Echo("hello", "world", "test")).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"hello world test"})
}

func TestEcho_EmptyString(t *testing.T) {
	result := run.Command(command.Echo("")).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{""})
}

func TestEcho_MultipleEmptyStrings(t *testing.T) {
	result := run.Command(command.Echo("", "", "")).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"  "})  // Two spaces between three empty strings
}

// ==============================================================================
// Test NoNewline Flag
// ==============================================================================

func TestEcho_NoNewline(t *testing.T) {
	result := run.Command(
		command.Echo("hello", command.NoNewline),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	// Without newline, output should not have a trailing newline
	// But our line capture will still split it as empty last element is removed
	assertion.Count(t, result.Stdout, 1)
	assertion.Equal(t, result.Stdout[0], "hello", "output")
}

func TestEcho_NoNewline_MultipleArgs(t *testing.T) {
	result := run.Command(
		command.Echo("one", "two", "three", command.NoNewline),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
	assertion.Equal(t, result.Stdout[0], "one two three", "output")
}

func TestEcho_WithNewline_Explicit(t *testing.T) {
	result := run.Command(
		command.Echo("hello", command.WithNewline),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"hello"})
}

// ==============================================================================
// Test Escape Sequences
// ==============================================================================

func TestEcho_Escape_Newline(t *testing.T) {
	result := run.Command(
		command.Echo(`line1\nline2`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"line1",
		"line2",
	})
}

func TestEcho_Escape_Tab(t *testing.T) {
	result := run.Command(
		command.Echo(`col1\tcol2\tcol3`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"col1\tcol2\tcol3"})
}

func TestEcho_Escape_CarriageReturn(t *testing.T) {
	result := run.Command(
		command.Echo(`text\rwith\rCR`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "\r")
}

func TestEcho_Escape_Backslash(t *testing.T) {
	result := run.Command(
		command.Echo(`path\\to\\file`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`path\to\file`})
}

func TestEcho_Escape_Alert(t *testing.T) {
	result := run.Command(
		command.Echo(`alert\a`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "\a")
}

func TestEcho_Escape_Backspace(t *testing.T) {
	result := run.Command(
		command.Echo(`text\bback`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "\b")
}

func TestEcho_Escape_FormFeed(t *testing.T) {
	result := run.Command(
		command.Echo(`page\fbreak`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "\f")
}

func TestEcho_Escape_VerticalTab(t *testing.T) {
	result := run.Command(
		command.Echo(`line\vvert`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "\v")
}

func TestEcho_Escape_Mixed(t *testing.T) {
	result := run.Command(
		command.Echo(`line1\nline2\ttab\\backslash`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2)
	assertion.Equal(t, result.Stdout[0], "line1", "first line")
	assertion.Contains(t, []string{result.Stdout[1]}, "line2")
	assertion.Contains(t, []string{result.Stdout[1]}, "\t")
	assertion.Contains(t, []string{result.Stdout[1]}, "\\")
}

func TestEcho_Escape_UnknownSequence(t *testing.T) {
	// Unknown escape sequences should be treated as literal backslash
	result := run.Command(
		command.Echo(`\xunknown`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`\xunknown`})
}

func TestEcho_Escape_Disabled(t *testing.T) {
	// Without EnableEscape, backslashes should be literal
	result := run.Command(
		command.Echo(`line1\nline2`, command.DisableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`line1\nline2`})
}

func TestEcho_Escape_Default_Disabled(t *testing.T) {
	// By default, escape sequences should not be processed
	result := run.Command(
		command.Echo(`line1\nline2`),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`line1\nline2`})
}

// ==============================================================================
// Test Combination of Flags
// ==============================================================================

func TestEcho_NoNewline_WithEscape(t *testing.T) {
	result := run.Command(
		command.Echo(`hello\nworld`, command.NoNewline, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2)
	assertion.Equal(t, result.Stdout[0], "hello", "first line")
	assertion.Equal(t, result.Stdout[1], "world", "second line (no trailing newline)")
}

func TestEcho_MultipleArgs_WithEscape(t *testing.T) {
	result := run.Command(
		command.Echo(`line1\nline2`, `tab\there`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2)
	assertion.Equal(t, result.Stdout[0], "line1", "first line")
	assertion.Contains(t, []string{result.Stdout[1]}, "line2")
	assertion.Contains(t, []string{result.Stdout[1]}, "\t")
}

// ==============================================================================
// Test Special Characters and Edge Cases
// ==============================================================================

func TestEcho_Whitespace(t *testing.T) {
	result := run.Command(command.Echo("   spaces   ")).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"   spaces   "})
}

func TestEcho_Tabs(t *testing.T) {
	result := run.Command(command.Echo("\t\ttabs\t\t")).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"\t\ttabs\t\t"})
}

func TestEcho_Unicode(t *testing.T) {
	result := run.Command(
		command.Echo("日本語", "中文", "한국어"),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"日本語 中文 한국어"})
}

func TestEcho_SpecialCharacters(t *testing.T) {
	result := run.Command(
		command.Echo("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"!@#$%^&*()_+-=[]{}|;':\",./<>?"})
}

func TestEcho_LongString(t *testing.T) {
	long := ""
	for i := 0; i < 1000; i++ {
		long += "a"
	}

	result := run.Command(command.Echo(long)).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{long})
}

func TestEcho_ManyArguments(t *testing.T) {
	args := make([]any, 100)
	for i := 0; i < 100; i++ {
		args[i] = "arg"
	}

	result := run.Command(command.Echo(args...)).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
	// 100 "arg"s with 99 spaces
	expected := ""
	for i := 0; i < 100; i++ {
		if i > 0 {
			expected += " "
		}
		expected += "arg"
	}
	assertion.Equal(t, result.Stdout[0], expected, "output")
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestEcho_OutputError(t *testing.T) {
	result := run.Command(
		command.Echo("test"),
	).WithStdoutError(errors.New("write failed")).WithStdin("").Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestEcho_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected []string
	}{
		{
			name:     "no args",
			args:     []any{},
			expected: []string{""},
		},
		{
			name:     "single word",
			args:     []any{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "multiple words",
			args:     []any{"one", "two", "three"},
			expected: []string{"one two three"},
		},
		{
			name:     "empty string",
			args:     []any{""},
			expected: []string{""},
		},
		{
			name:     "with spaces",
			args:     []any{"hello world", "foo bar"},
			expected: []string{"hello world foo bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Echo(tt.args...)).WithStdin("").Run()
			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

func TestEcho_EscapeSequences_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		contains string  // For non-printable chars
	}{
		{
			name:     "newline",
			input:    `line1\nline2`,
			contains: "\n",
		},
		{
			name:     "tab",
			input:    `a\tb`,
			contains: "\t",
		},
		{
			name:     "carriage return",
			input:    `a\rb`,
			contains: "\r",
		},
		{
			name:     "backslash",
			input:    `a\\b`,
			expected: `a\b`,
		},
		{
			name:     "alert",
			input:    `a\ab`,
			contains: "\a",
		},
		{
			name:     "backspace",
			input:    `a\bb`,
			contains: "\b",
		},
		{
			name:     "form feed",
			input:    `a\fb`,
			contains: "\f",
		},
		{
			name:     "vertical tab",
			input:    `a\vb`,
			contains: "\v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(
				command.Echo(tt.input, command.EnableEscape),
			).WithStdin("").Run()

			assertion.NoError(t, result.Err)

			if tt.expected != "" {
				assertion.Contains(t, result.Stdout, tt.expected)
			}
			if tt.contains != "" {
				fullOutput := result.Stdout[0]
				if len(result.Stdout) > 1 {
					// Join with actual newlines for multi-line output
					for i := 1; i < len(result.Stdout); i++ {
						fullOutput += "\n" + result.Stdout[i]
					}
				}
				assertion.True(t,
					len(fullOutput) > 0 && fullOutput != tt.input,
					"escape sequence should be processed")
			}
		})
	}
}

// ==============================================================================
// Test Edge Cases with processEscapes
// ==============================================================================

func TestEcho_TrailingBackslash(t *testing.T) {
	// Trailing backslash with no following character
	result := run.Command(
		command.Echo(`test\`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`test\`})
}

func TestEcho_MultipleBackslashes(t *testing.T) {
	result := run.Command(
		command.Echo(`\\\\`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`\\`})
}

func TestEcho_BackslashAtEnd(t *testing.T) {
	result := run.Command(
		command.Echo(`path\to\file\`, command.EnableEscape),
	).WithStdin("").Run()

	assertion.NoError(t, result.Err)
	// Backslash at end with no following char stays as backslash
	assertion.Contains(t, result.Stdout, "path")
}

