# go-commands-parser

This is a helper utility for the `go-commands` package.

# Usage

    npm install shadowmint/go-commands-parser --save

Now define a command syntax and assign a handler for it:

    import (
        "reflect"

        "ntoolkit/commands"
        "ntoolkit/commands/cparser"
        "ntoolkit/errors"
        "ntoolkit/events"
        "ntoolkit/futures"
    )

    type ErrInvalidDragon struct{}

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

    func (handler *PutCommandHandler) Handles() reflect.Type {
        return reflect.TypeOf(&PutCommand{})
    }

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

...and assign them all to a commands instance:

    parser := cparser.New()

    parser.Register(parser.Command("put", "[item]", "on", "[target]").With(func(params map[string]string, context interface{}) (commands.Command, error) {
        player, ok := context.(*Player)
        if ok {
          ...
          return &PutCommand{Item: params["item"], Target: params["target"]}, nil
        } else {
          // Return some error...
        }
    }))

    parser.Register(parser.Command().Word("put", true).With(func(params map[string]string, context interface{}) (commands.Command, error) {
        return nil, errors.Fail(cparser.ErrBadSyntax{}, nil, "Invalid put command; try put ITEM on TARGET")
    }))

    // Command handlers
    parser.Commands.Register(&PutCommandHandler{})

Finally, you can execute a command:

    p.Execute("put foo on bar", player).Then(func(cmd commands.Command) {
        pcmd, ok := cmd.(*PutCommand)
        ...
    }, func(err error) {
        ...
    })

Notice that `Execute` and `Wait` take an arbitrary context object that allows the `CommandFactory` to build a specific command
given the execution context. For example, you might want to pass in the requester of the command, the application state, etc.
