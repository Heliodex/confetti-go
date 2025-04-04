package main

import (
	"errors"
	"fmt"
)

// The Confetti language consists of zero or more directives. A directive consists of one or more arguments and optional subdirectives.

// The entire AST of the language is ONE struct!!!!
type Directive struct {
	Arguments     []string
	Subdirectives []Directive
}

func parse(ts []Token) (p []Directive, err error) {
	var current Directive
	push := func() {
		if current.Arguments == nil {
			return
		}
		p = append(p, current)
		current = Directive{}
	}

	prevSignificant := func(i int) (prev Token) {
		for i--; i > 0; i-- {
			if prev = ts[i]; prev.Type != TokWhitespace && prev.Type != TokComment {
				return
			}
		}
		return
	}

	for i := 0; i < len(ts); i++ {
		switch t := ts[i]; t.Type {
		case TokArgument:
			current.Arguments = append(current.Arguments, t.Content)

		case TokSemicolon: // end of directive
			if prev := prevSignificant(i); prev.Type == TokSemicolon || prev.Type == TokNewline || prev.Type == TokLineContinuation {
				return nil, errors.New("unexpected ';'")
			}
			fallthrough

		case TokNewline: // end of directive
			push()

		case TokComment, TokWhitespace: // Ignore whitespace and comments

		case TokOpenBrace:
			if i == len(ts)-1 || prevSignificant(i).Type == TokSemicolon {
				// fmt.Println(prevNonWhitespace(i).Type == TokSemicolon)
				return nil, fmt.Errorf("unexpected '{'")
			}

			// Get all tokens until next close brace
			var subts []Token
			var depth int // also account for nested

			for i++; i < len(ts); i++ {
				// escapes should be dealt with in lexer
				if t = ts[i]; t.Type == TokOpenBrace {
					depth++
				} else if t.Type == TokCloseBrace {
					depth--
				}

				if depth < 0 {
					break
				}
				subts = append(subts, t)
			}

			if depth >= 0 {
				return nil, fmt.Errorf("expected '}'")
			}

			subp, err := parse(subts)
			if err != nil {
				return nil, err
			} else if current.Arguments == nil {
				// push to the previous directive
				p[len(p)-1].Subdirectives = subp
				break
			}

			current.Subdirectives = subp
			push()

		case TokCloseBrace:
			return nil, errors.New("found '}' without matching '{'")

		case TokLineContinuation:
			if current.Arguments == nil {
				return nil, fmt.Errorf("unexpected line continuation")
			}
		}
	}

	push()
	return
}
