package cparser

import (
	"ntoolkit/commands"
	"ntoolkit/parser"
)

type CommandFactory interface {
	// Parse yields a command if the tokenList is valid, or an error.
	// If the token stream is not valid it should return (nil, nil);
	// errors should only be returned if the token stream is invalid.
	Parse(tokenList *parser.Tokens) (commands.Command, error)
}
