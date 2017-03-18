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

type LookCommandFactory struct {
}

func (factory *LookCommandFactory) Parse(tokenList *parser.Tokens) (commands.Command, error) {
	if tokenList.Front != nil {
		if tokenList.Front.Is(tools.TokenTypeBlock, "look") {
			direction := tokenList.Front.Next
			if direction == nil {
				return nil, errors.Fail(cparser.ErrBadSyntax{}, nil, "Invalid look syntax; try 'lookZ DIRECTION'")
			} else {
				return &LookCommand{Direction: direction.CollectRaw(" ")}, nil
			}
		}
	}
	return nil, nil
}

type LookCommand struct {
	eventHandler *events.EventHandler
	Success      bool
	Direction    string
}

func (cmd *LookCommand) EventHandler() *events.EventHandler {
	if cmd.eventHandler == nil {
		cmd.eventHandler = events.New()
	}
	return cmd.eventHandler
}

type LookCommandHandler struct {
}

// Handles returns the type supported by this command handler
func (handler *LookCommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&LookCommand{})
}

// Execute executes the command given and returns an error on failure
func (handler *LookCommandHandler) Execute(command interface{}) *futures.Deferred {
	lcmd := command.(*LookCommand)
	rtn := &futures.Deferred{}
	lcmd.Success = true
	rtn.Resolve()
	return rtn
}
