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
	"fmt"
)

// UseCommandErrBadTool is for when trying to use something that isn't a tool
type UseCommandErrBadTool struct{}

// UseCommandErrBadUse is for when trying to use something on something that just isn't right
type UseCommandErrBadUse struct{}

type UseCommandFactory struct {
}

func (factory *UseCommandFactory) Parse(tokenList *parser.Tokens) (commands.Command, error) {
	if tokenList.Front != nil {
		if tokenList.Front.Is(tools.TokenTypeBlock, "use") {
			tool := tokenList.Front.Next
			if tool != nil && tool.Next != nil && tool.Next.Is(tools.TokenTypeBlock, "on") {
				target := tool.Next.Next
				if target != nil {
					return &UseCommand{Tool: tool.CollectRaw(" "), Target: target.CollectRaw(" ")}, nil
				}
			}
			return nil, errors.Fail(cparser.ErrBadSyntax{}, nil, "Invalid Use syntax; try 'use TOOL on OBJECT'")
		}
	}
	return nil, nil
}

type UseCommand struct {
	eventHandler *events.EventHandler
	Target       string
	Tool         string
}

func (cmd *UseCommand) EventHandler() *events.EventHandler {
	if cmd.eventHandler == nil {
		cmd.eventHandler = events.New()
	}
	return cmd.eventHandler
}

type UseCommandHandler struct {
}

// Handles returns the type supported by this command handler
func (handler *UseCommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&UseCommand{})
}

var validToolTypes = []string{"hammer", "spoon"}
var validObjectTypes = []string{"door", "chair"}

// Execute executes the command given and returns an error on failure
func (handler *UseCommandHandler) Execute(command interface{}) *futures.Deferred {
	rtn := &futures.Deferred{}
	ucmd := command.(*UseCommand)
	if !handler.in(validToolTypes, ucmd.Tool) {
		rtn.Reject(errors.Fail(UseCommandErrBadTool{}, nil, fmt.Sprintf("%s is not a valid tool", ucmd.Tool)))
	} else if !handler.in(validObjectTypes, ucmd.Target) {
		rtn.Reject(errors.Fail(UseCommandErrBadUse{}, nil, fmt.Sprintf("%s is not a valid target for tool %s", ucmd.Target, ucmd.Tool)))
	} else {
		rtn.Resolve()
	}
	return rtn
}

func (handler *UseCommandHandler) in(targets []string, target string) bool {
	for i := 0; i < len(targets); i++ {
		if targets[i] == target {
			return true
		}
	}
	return false
}