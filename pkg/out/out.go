package out

import (
	"os"
)

// Formatter is an interface that formats a string.
type Formatter func(*os.File, string) string

// Out is a struct that contains an output destination.
type Out struct {
	Out       *os.File
	Formatter Formatter
}

// NewOutput returns a new Output struct.
func New(out *os.File, formatter Formatter) *Out {
	return &Out{
		Out:       out,
		Formatter: formatter,
	}
}

// Emit writes to the output destination.
func (o *Out) Emit(s string) {
	// if no formatter is provided, just write the string
	if o.Formatter == nil {
		if _, err := o.Out.Write([]byte(s)); err != nil {
			panic(err)
		}
		return
	}

	// otherwise, format the string and write it
	formatted := o.Formatter(o.Out, s)
	if _, err := o.Out.Write([]byte(formatted)); err != nil {
		panic(err)
	}
}
