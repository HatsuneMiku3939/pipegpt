package stdin

import (
	"bufio"
	"os"
)

// ConsumeStdin reads from stdin and returns the contents as a string.
func ConsumeStdin() string {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}
	return input
}
