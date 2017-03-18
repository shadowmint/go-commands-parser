package cparser_test

import (
	"fmt"
	"reflect"

	"ntoolkit/commands"
	"ntoolkit/commands/cparser"
	"ntoolkit/errors"
	"ntoolkit/events"
	"ntoolkit/futures"
)

func registerUse2Factory(parser *cparser.CommandParser) {
	parser.Register(parser.Command("use2", "[tool]", "on", "[target]").With(func(params map[string]string) (commands.Command, error) {
		return &Use2Command{Target: params["target"], Tool: params["tool"]}, nil
	}))
}

type Use2Command struct {
	eventHandler *events.EventHandler
	Target       string
	Tool         string
}

func (cmd *Use2Command) EventHandler() *events.EventHandler {
	if cmd.eventHandler == nil {
		cmd.eventHandler = events.New()
	}
	return cmd.eventHandler
}

type Use2CommandHandler struct {
}

// Handles returns the type supported by this command handler
func (handler *Use2CommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&Use2Command{})
}

// Execute executes the command given and returns an error on failure
func (handler *Use2CommandHandler) Execute(command interface{}) *futures.Deferred {
	rtn := &futures.Deferred{}
	ucmd := command.(*Use2Command)
	if !handler.in(validToolTypes, ucmd.Tool) {
		rtn.Reject(errors.Fail(UseCommandErrBadTool{}, nil, fmt.Sprintf("%s is not a valid tool", ucmd.Tool)))
	} else if !handler.in(validObjectTypes, ucmd.Target) {
		rtn.Reject(errors.Fail(UseCommandErrBadUse{}, nil, fmt.Sprintf("%s is not a valid target for tool %s", ucmd.Target, ucmd.Tool)))
	} else {
		rtn.Resolve()
	}
	return rtn
}

func (handler *Use2CommandHandler) in(targets []string, target string) bool {
	for i := 0; i < len(targets); i++ {
		if targets[i] == target {
			return true
		}
	}
	return false
}
