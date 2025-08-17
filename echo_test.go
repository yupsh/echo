package echo_test

import (
	"context"
	"os"

	"github.com/yupsh/echo"
	"github.com/yupsh/echo/opt"
)

func ExampleEcho() {
	ctx := context.Background()

	cmd := echo.Echo("Hello", "World")
	cmd.Execute(ctx, nil, os.Stdout, os.Stderr)
	// Output: Hello World
}

func ExampleEcho_noNewline() {
	ctx := context.Background()

	cmd := echo.Echo("No newline", opt.NoNewline)
	cmd.Execute(ctx, nil, os.Stdout, os.Stderr)
	// Output: No newline
}

func ExampleEcho_withEscapes() {
	ctx := context.Background()

	cmd := echo.Echo("Line1\\nLine2\\tTab", opt.EnableEscape)
	cmd.Execute(ctx, nil, os.Stdout, os.Stderr)
	// Output: Line1
	// Line2	Tab
}
