package in

import (
	"bufio"
	"fmt"
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

// Consume reads from the input source until EOF
func (i *In) Consume() (string, error) {
	scanner := bufio.NewScanner(i.In)
	var input string
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return "", err
		}

		input += scanner.Text() + "\n"
	}

	return input, nil
}

// Consume reads from the input source until stopword is encountered
func (i *In) ConsumeUntil(prompt string, stopword string, callback func(string, error)) {
	// show prompt
	fmt.Printf("%s", prompt)

	// read from input
	scanner := bufio.NewScanner(i.In)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			callback("", err)

			// show prompt
			fmt.Printf("%s", prompt)
			continue
		}
		text := scanner.Text()
		if text == stopword {
			break
		}

		callback(text, nil)

		// show prompt
		fmt.Printf("%s", prompt)
	}
}
