package bar

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type tokens []token

type token interface {
	debug(*Bar) string
	print(*Bar) string
}

type tokenFormat struct {
	stream *bufio.Reader
}

type spaceToken struct{}
type barToken struct{}
type percentToken struct{}
type rateToken struct{}
type customVerbToken struct {
	verb string
}
type literalToken struct {
	content string
}

func tokenize(f string, customVerbs []string) tokens {
	var t tokens

	sr := strings.NewReader(f)
	r := &tokenFormat{bufio.NewReader(sr)}

	for {
		tkn, err := r.nextToken(customVerbs)
		if err != nil {
			if err == io.EOF {
				return t
			}

			panic(fmt.Sprintf("tokenize: %v", err))
		}

		t = append(t, tkn)
	}
}

func (f *tokenFormat) nextToken(customVerbs []string) (token, error) {
	for {
		r, _, err := f.stream.ReadRune()
		if err != nil {
			return nil, err
		}

		switch r {
		case ' ':
			return spaceToken{}, nil
		case ':':
			return f.readAction(customVerbs)
		default:
			return f.readLiteral(r)
		}
	}
}

func (f *tokenFormat) readAction(customVerbs []string) (token, error) {
	var verb bytes.Buffer

	for {
		r, _, err := f.stream.ReadRune()

		if err != nil {
			if err == io.EOF {
				t, _ := tokenFromString(verb.String(), customVerbs)
				return t, nil
			}

			return nil, err
		}

		verb.Write([]byte(string([]rune{r})))

		if t, literal := tokenFromString(verb.String(), customVerbs); !literal {
			return t, nil
		}

		if f.readSeparator() {
			t, _ := tokenFromString(verb.String(), customVerbs)
			return t, nil
		}
	}
}

func (f *tokenFormat) readLiteral(prefix rune) (token, error) {
	var value bytes.Buffer

	value.Write([]byte(string([]rune{prefix})))

	for {
		r, _, err := f.stream.ReadRune()

		if err != nil {
			if err == io.EOF {
				t, _ := tokenFromString(value.String(), nil)
				return t, nil
			}

			return nil, err
		}

		value.Write([]byte(string([]rune{r})))

		if f.readSeparator() {
			t, _ := tokenFromString(value.String(), nil)
			return t, nil
		}
	}
}

func (f *tokenFormat) readSeparator() bool {
	p, err := f.stream.Peek(1)
	if err != nil || p[0] == ' ' || p[0] == ':' {
		return true
	}
	return false
}

// tokenFromString will return the token parsed from s; literal is true as
// long as the token is not a literal
func tokenFromString(s string, customVerbs []string) (token, bool) {
	// check for standard verbs
	switch s {
	case "bar":
		return barToken{}, false
	case "percent":
		return percentToken{}, false
	case "rate":
		return rateToken{}, false
	}

	// check for custom verbs
	for _, verb := range customVerbs {
		if s == verb {
			return customVerbToken{verb}, false
		}
	}

	return literalToken{s}, true
}

func (t spaceToken) print(_ *Bar) string {
	return " "
}

func (t barToken) print(b *Bar) string {
	p := int(b.prog() * float64(b.width))
	return fmt.Sprintf(
		"%s%s%s%s%s",
		b.start,
		strings.Repeat(b.complete, int(math.Max(0, float64(p-1)))),
		b.head,
		strings.Repeat(b.incomplete, b.width-p),
		b.end,
	)
}

func (t percentToken) print(b *Bar) string {
	return fmt.Sprintf("%.1f%%", b.prog()*100)
}

func (t rateToken) print(b *Bar) string {
	return fmt.Sprintf("%.1f", b.rate)
}

func (t customVerbToken) print(b *Bar) string {
	for _, def := range b.context {
		if def.verb == t.verb {
			return def.value.String()
		}
	}

	fmt.Fprintf(os.Stderr, "tokenize: only use `:` for custom verbs")
	return t.verb
}

func (t literalToken) print(_ *Bar) string {
	return t.content
}

func (t spaceToken) debug(b *Bar) string {
	return " "
}

func (t barToken) debug(b *Bar) string {
	return fmt.Sprintf("<barToken p={%d} t={%d}>", b.progress, b.total)
}

func (t percentToken) debug(b *Bar) string {
	return fmt.Sprintf("<percentToken \"%s\">", t.print(b))
}

func (t rateToken) debug(b *Bar) string {
	return fmt.Sprintf("<rateToken \"%s\">", t.print(b))
}

func (t customVerbToken) debug(b *Bar) string {
	return fmt.Sprintf("<customVerbToken verb=\"%s\" value=\"%s\">", t.verb, t.print(b))
}

func (t literalToken) debug(b *Bar) string {
	return fmt.Sprintf("<literalToken \"%s\">", t.content)
}
