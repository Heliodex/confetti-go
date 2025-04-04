package main

import (
	"errors"
	"fmt"
)

// The Confetti language consists of zero or more directives. A directive consists of one or more arguments and optional subdirectives.

type Argument []rune

type Directive struct {
	Arguments     []Argument
	Subdirectives []Directive
}

func parse(lexed []Token) (p []Directive, err error) {
	var current Directive
	push := func() {
		if len(current.Arguments) == 0 {
			return
		}
		p = append(p, current)
		current = Directive{}
	}

	prevNonWhitespace := func(i int) (prev Token) {
		for i--; i > 0; i-- {
			if prev = lexed[i]; prev.Type != TokWhitespace && prev.Type != TokComment {
				return
			}
		}
		return
	}

	for i := 0; i < len(lexed); i++ {
		switch t := lexed[i]; t.Type {
		case TokArgument:
			current.Arguments = append(current.Arguments, Argument(t.Content))

		case TokSemicolon: // end of directive
			if prev := prevNonWhitespace(i); prev.Type == TokSemicolon || prev.Type == TokNewline {
				return nil, errors.New("unexpected ';'")
			}
			fallthrough

		case TokNewline: // end of directive
			push()

		case TokComment, TokWhitespace: // Ignore whitespace and comments

		case TokOpenBrace:
			if i == len(lexed)-1 || prevNonWhitespace(i).Type == TokSemicolon {
				// fmt.Println(prevNonWhitespace(i).Type == TokSemicolon)
				return nil, fmt.Errorf("unexpected '{'")
			}

			// Get all tokens until next close brace
			var ts []Token

			depth := 1 // also account for nested

			for i++; i < len(lexed); i++ {
				// escapes should be dealt with in lexer
				if t = lexed[i]; t.Type == TokOpenBrace {
					depth++
				} else if t.Type == TokCloseBrace {
					depth--
				}

				if depth == 0 {
					break
				}
				ts = append(ts, t)
			}

			if depth != 0 {
				return nil, fmt.Errorf("expected '}'")
			}

			subdirs, err := parse(ts)
			if err != nil {
				return nil, err
			} else if len(current.Arguments) == 0 {
				// push to the previous directive
				p[len(p)-1].Subdirectives = subdirs
				break
			}

			current.Subdirectives = subdirs
			push()

		case TokCloseBrace:
			return nil, errors.New("found '}' without matching '{'")

		case TokReverseSolidus:
			// escape character
			i++
			if i >= len(lexed) {
				break
			} else if t = lexed[i]; t.Type == TokNewline && len(current.Arguments) == 0 {
				return nil, fmt.Errorf("unexpected line continuation")
			}
		}
	}

	push()
	return
}
