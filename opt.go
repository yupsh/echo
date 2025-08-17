package command

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

type flags struct {
	NoNewline NoNewlineFlag
	Escape    EscapeFlag
}

func (f NoNewlineFlag) Configure(flags *flags) { flags.NoNewline = f }
func (f EscapeFlag) Configure(flags *flags)    { flags.Escape = f }
