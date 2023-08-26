package out

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
)

// MarkdownFormatter is a Formatter that formats a string as markdown. if current session is tty.
func MarkdownFormatter(out *os.File, s string) string {
	// if the output is not a tty, just return the string
	if isTTY := isatty.IsTerminal(out.Fd()); !isTTY {
		return s
	}

	// create a new markdown renderer
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		// if the renderer fails to initialize, just return the string
		return s
	}

	formatted, err := r.Render(s)
	if err != nil {
		// if the renderer fails to render, just return the string
		return s
	}

	return formatted
}
