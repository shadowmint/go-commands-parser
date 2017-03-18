package cparser_test

import (
	"reflect"

	"ntoolkit/commands"
	"ntoolkit/commands/cparser"
	"ntoolkit/errors"
	"ntoolkit/events"
	"ntoolkit/futures"
)

type ErrInvalidDragon struct{}

func registerPutFactory(parser *cparser.CommandParser) {
	parser.Register(parser.Command("put", "[item]", "on", "[target]").With(func(params map[string]string) (commands.Command, error) {
		return &PutCommand{Item: params["item"], Target: params["target"]}, nil
	}))
	parser.Register(parser.Command("put", "[item]", "in", "[container]").With(func(params map[string]string) (commands.Command, error) {
		return &PutCommand{Item: params["item"], Container: params["container"]}, nil
	}))
	parser.Register(parser.Command().Word("put", true).With(func(params map[string]string) (commands.Command, error) {
		return nil, errors.Fail(cparser.ErrBadSyntax{}, nil, "Invalid put command; try put ITEM on TARGET")
	}))
}

type PutCommand struct {
	eventHandler *events.EventHandler
	Item         string
	Target       string
	Container    string
}

func (cmd *PutCommand) EventHandler() *events.EventHandler {
	if cmd.eventHandler == nil {
		cmd.eventHandler = events.New()
	}
	return cmd.eventHandler
}

type PutCommandHandler struct {
}

// Handles returns the type supported by this command handler
func (handler *PutCommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&PutCommand{})
}

// Execute executes the command given and returns an error on failure
func (handler *PutCommandHandler) Execute(command interface{}) *futures.Deferred {
	rtn := &futures.Deferred{}
	pcmd := command.(*PutCommand)
	if pcmd.Item == "dragon" {
		rtn.Reject(errors.Fail(ErrInvalidDragon{}, nil, "You cannot do that to a dragon"))
	} else {
		rtn.Resolve()
	}
	return rtn
}
