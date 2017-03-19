package cparser

import (
	"fmt"
	"strings"

	"ntoolkit/commands"
	"ntoolkit/errors"
	"ntoolkit/parser"
)

const (
	standardCommandTypeWord  = iota
	standardCommandTypeToken = iota
)

// standardCommandWord
type standardCommandWord struct {
	// Type of this word
	Type int

	// The name of this word
	Name string

	// If unique, matching this token and not all others generates a syntax error.
	// For example, if you want 'go home now' to be an error, make word 'go' unique.
	Unique bool
}

// StandardCommandFactory is a CommandFactory for a command in the form
// WORD [token | WORD]... that generates a map of keyword -> token if
// matching, and returns a syntax error otherwise. If this context a WORD
// is a constant value, like 'go' in 'go north'.
type StandardCommandFactory struct {
	// List of items that work
	items []standardCommandWord

	// Invoked after successful parse check to generate a command.
	handler func(params map[string]string, context interface{}) (commands.Command, error)
}

// newStandardCommandFactory creates an returns a command factory
func newStandardCommandFactory() *StandardCommandFactory {
	return &StandardCommandFactory{items: make([]standardCommandWord, 0)}
}

// Word adds a word to the command syntax, and returns the instance.
func (factory *StandardCommandFactory) Word(word string, isUnique ...bool) *StandardCommandFactory {
	isWordUnique := false
	if len(isUnique) > 0 {
		isWordUnique = isUnique[0]
	}
	factory.items = append(factory.items, standardCommandWord{
		Type:   standardCommandTypeWord,
		Name:   word,
		Unique: isWordUnique})
	return factory
}

// Token adds a token to the command syntax, and returns the instance.
func (factory *StandardCommandFactory) Token(tokenName string) *StandardCommandFactory {
	factory.items = append(factory.items, standardCommandWord{
		Type:   standardCommandTypeToken,
		Name:   tokenName,
		Unique: false})
	return factory
}

// With sets the handler to generate a command on the factory
func (factory *StandardCommandFactory) With(factoryFunc func(params map[string]string, context interface{}) (commands.Command, error)) *StandardCommandFactory {
	factory.handler = factoryFunc
	return factory
}

// String renders the factory as a string list
func (factory *StandardCommandFactory) String() string {
	buffer := make([]string, len(factory.items))
	for i := range factory.items {
		item := factory.items[i]
		if item.Type == standardCommandTypeWord {
			buffer[i] = item.Name
		} else if item.Type == standardCommandTypeToken {
			buffer[i] = fmt.Sprintf("[%s]", item.Name)
		}
	}
	return strings.Join(buffer, " ")
}

// Parse checks the token list against the defined syntax and raises and error if it doesn't work.
// Notice that
func (factory *StandardCommandFactory) Parse(tokenList *parser.Tokens, context interface{}) (cmd commands.Command, err error) {
	defer (func() {
		r := recover()
		if r != nil {
			err = r.(error)
			cmd = nil
		}
	})()

	// setup
	params := make(map[string]string)
	foundUnique := false
	foundMatch := 0
	marker := tokenList.Front

	// collect values
	for offset := 0; offset < len(factory.items) && marker != nil; offset++ {
		raw := marker.CollectRaw()
		item := factory.items[offset]
		if item.Type == standardCommandTypeWord {
			// TODO: Capitialization check?
			if raw == item.Name {
				foundMatch += 1
				if item.Unique {
					foundUnique = true
				}
			}
		} else if item.Type == standardCommandTypeToken {
			params[item.Name] = raw
			foundMatch += 1
		}
		marker = marker.Next
	}

	// validate; error if not right length but we found any unique tokens
	// If we found no match, this handler isn't the right one.
	if foundMatch != len(factory.items) {
		if foundUnique {
			return nil, errors.Fail(ErrBadSyntax{}, nil, fmt.Sprintf("Invalid syntax for command, did not match: %s", factory))
		} else {
			return nil, nil
		}
	}

	// ! Someone forget to call With()
	if factory.handler == nil {
		return nil, errors.Fail(ErrBadSyntax{}, nil, "No handler attached to standard command factory")
	}

	// Try to get a command back
	return factory.handler(params, context)
}
