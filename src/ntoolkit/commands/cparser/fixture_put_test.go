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

func putOnHandler(params map[string]string, context interface{}) (commands.Command, error) {
	playerId := context.(int)
	return &PutCommand{Item: params["item"], Target: params["target"], PlayerId: playerId}, nil
}

func putInHandler(params map[string]string, context interface{}) (commands.Command, error) {
	playerId := context.(int)
	return &PutCommand{Item: params["item"], Container: params["container"], PlayerId: playerId}, nil
}

func putDefaultHandler(params map[string]string, context interface{}) (commands.Command, error) {
	return nil, errors.Fail(cparser.ErrBadSyntax{}, nil, "Invalid put command; try put ITEM on TARGET")
}

func registerPutFactory(parser *cparser.CommandParser) {
	parser.Register(parser.Command("put", "[item]", "on", "[target]").With(putOnHandler))
	parser.Register(parser.Command("put", "[item]", "in", "[container]").With(putInHandler))
	parser.Register(parser.Command().Word("put", true).With(putDefaultHandler))
}

type PutCommand struct {
	eventHandler *events.EventHandler
	Item         string
	Target       string
	Container    string
	PlayerId     int
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
