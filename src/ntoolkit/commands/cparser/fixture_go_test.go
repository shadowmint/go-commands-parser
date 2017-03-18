package cparser_test

import (
	"reflect"

	"ntoolkit/commands"
	"ntoolkit/commands/cparser"
	"ntoolkit/errors"
	"ntoolkit/events"
	"ntoolkit/futures"
	"ntoolkit/parser"
	"ntoolkit/parser/tools"
)

type GoCommandFactory struct {
}

func (factory *GoCommandFactory) Parse(tokenList *parser.Tokens) (commands.Command, error) {
	if tokenList.Front != nil {
		if tokenList.Front.Is(tools.TokenTypeBlock, "go") {
			direction := tokenList.Front.Next
			if direction == nil {
				return nil, errors.Fail(cparser.ErrBadSyntax{}, nil, "Invalid go syntax; try 'go DIRECTION'")
			} else {
				return &GoCommand{Direction: direction.CollectRaw( " ")}, nil
			}
		}
	}
	return nil, nil
}

type GoCommand struct {
	eventHandler *events.EventHandler
	Success      bool
	Direction    string
}

func (cmd *GoCommand) EventHandler() *events.EventHandler {
	if cmd.eventHandler == nil {
		cmd.eventHandler = events.New()
	}
	return cmd.eventHandler
}

type GoCommandHandler struct {
}

// Handles returns the type supported by this command handler
func (handler *GoCommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&GoCommand{})
}

// Execute executes the command given and returns an error on failure
func (handler *GoCommandHandler) Execute(command interface{}) *futures.Deferred {
	gcmd := command.(*GoCommand)
	rtn := &futures.Deferred{}
	gcmd.Success = true
	rtn.Resolve()
	return rtn
}
