// package confetti-go implements the Confetti configuration language.
package main

import (
	"fmt"
	"strings"
)

type Extensions map[string]string

func (e Extensions) Has(name string) bool {
	_, ok := e[name]
	return ok
}

var extensions = Extensions{
	"c_style_comments": "",
	"expression_arguments":  "",
	// "punctuator_arguments": "",
}

const text = `foo
{
    bar
}
;
baz`

func printDirective(d Directive, depth int) {
	prefix := strings.Repeat("  ", depth)

	fmt.Println(prefix + "Directive:")
	for _, arg := range d.Arguments {
		fmt.Printf(prefix+"  %q\n", arg)
	}
	for _, sub := range d.Subdirectives {
		printDirective(sub, depth+1)
	}
}

func main() {
	fmt.Println("Lexing")

	ts, err := lex(text, extensions)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsing")

	// for _, t := range ts {
	// 	if t.Type == TokWhitespace {
	// 		continue
	// 	}
	// 	fmt.Printf("%14s  %s\n", t.Type, t.Content)
	// }

	p, err := parse(ts, extensions)
	if err != nil {
		panic(err)
	}

	for _, d := range p {
		printDirective(d, 0)
	}
}
