package opt

// Boolean flag types with constants
type NoNewlineFlag bool

const (
	NoNewline   NoNewlineFlag = true
	WithNewline NoNewlineFlag = false
)

type EscapeFlag bool

const (
	EnableEscape  EscapeFlag = true
	DisableEscape EscapeFlag = false
)

// Flags represents the configuration options for the echo command
type Flags struct {
	NoNewline NoNewlineFlag // Don't output trailing newline
	Escape    EscapeFlag    // Enable interpretation of backslash escapes
}

// Flag configuration methods
func (f NoNewlineFlag) Configure(flags *Flags) { flags.NoNewline = f }
func (f EscapeFlag) Configure(flags *Flags)    { flags.Escape = f }
