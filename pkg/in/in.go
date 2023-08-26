package in

import (
	"bufio"
	"os"
)

// In is a struct that contains an input source.
type In struct {
	In *os.File
}

// NewInput returns a new Input struct.
func New(in *os.File) *In {
	return &In{In: in}
}

// Consume reads from the input source until break character is encountered.
func (i *In) Consume(breakChar byte) string {
	scanner := bufio.NewScanner(i.In)
	var input string
	for scanner.Scan() {
		input += scanner.Text() + "\n"

		// if break character is encountered, break
		bytes := scanner.Bytes()
		if len(bytes) > 1 && bytes[len(bytes)-1] == breakChar {
			break
		}
	}

	return input
}
