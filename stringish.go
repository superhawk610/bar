package bar

// since Go strings don't satisfy the fmt.Stringer interface, we
// need to wrap them in a struct that does

import (
	"fmt"
)

type stringish struct {
	content fmt.Stringer
}

type stringStringer struct {
	s string
}

func newStringish(content interface{}) *stringish {
	switch content := content.(type) {
	case string:
		return &stringish{
			content: &stringStringer{s: content},
		}
	case fmt.Stringer:
		return &stringish{content}
	default:
		panic(fmt.Sprintf("stringish: %v", content))
	}
}

func (s *stringStringer) String() string {
	return s.s
}

func (s *stringish) String() string {
	return s.content.String()
}
