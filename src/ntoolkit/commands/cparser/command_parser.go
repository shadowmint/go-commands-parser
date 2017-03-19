package cparser

import (
	"sync"

	"ntoolkit/commands"
	"ntoolkit/errors"
	"ntoolkit/parser/tools"
)

// CommandParser is a high level interface for dispatching text commands.
type CommandParser struct {
	Commands    *commands.Commands
	blockParser *tools.BlockParser
	factory     []CommandFactory
}

// New returns a new command cparser with the attached commands object.
// If no commands object is supplied a new blank on is created.
func New(cmd ...*commands.Commands) *CommandParser {
	var commander *commands.Commands
	if len(cmd) > 0 {
		commander = cmd[0]
	}
	if commander == nil {
		commander = commands.New()
	}
	return &CommandParser{
		Commands:    commander,
		blockParser: tools.NewBlockParser(),
		factory:     make([]CommandFactory, 0)  }
}

func (p *CommandParser) Execute(command string, context interface{}) (promise *DeferredCommand) {
	defer (func() {
		r := recover()
		if r != nil {
			err := r.(error)
			p.failed(errors.Fail(ErrCommandFailed{}, err, err.Error()))
		}
	})()
	p.blockParser.Parse(command)
	tokens, err := p.blockParser.Finished()
	if err != nil {
		return p.failed(errors.Fail(ErrBadSyntax{}, err, "Invalid command string"))
	}
	for i := range p.factory {
		cmd, err := p.factory[i].Parse(tokens, context)
		if err != nil {
			return p.failed(errors.Fail(ErrCommandFailed{}, err, "Command syntax error"))
		}
		if cmd != nil {
			rtn := &DeferredCommand{}
			p.Commands.Execute(cmd).Then(func() {
				rtn.Resolve(cmd)
			}, func(err error) {
				rtn.Reject(errors.Fail(ErrCommandFailed{}, err, "Command failed to execute"))
			})
			return rtn
		}
	}
	return p.failed(errors.Fail(ErrNoHandler{}, nil, "No handler supported the given command"))
}

// Wait for an executed command to resolve and return nil or the error.
func (p *CommandParser) Wait(command string, context interface{}) (commands.Command, error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	var err error
	var cmd commands.Command
	p.Execute(command, context).Then(func(c commands.Command) {
		cmd = c
		wg.Done()
	}, func(errRtn error) {
		err = errRtn
		wg.Done()
	})
	wg.Wait()
	return cmd, err
}

// Register a new command factory to handle some kind of input.
func (p *CommandParser) Register(factory CommandFactory) {
	p.factory = append(p.factory, factory)
}

// Command returns a new standard command factory; you can use .Word() and .Token()
// on the returned object, or just pass in args; "go" -> Word() and "[name]" -> Token().
func (p *CommandParser) Command(words ...string) *StandardCommandFactory {
	factory := newStandardCommandFactory()
	for i := range words {
		word := words[i]
		if len(word) > 2 && word[0] == '[' && word[len(word)-1] == ']' {
			factory.Token(word[1:len(word)-1])
		} else {
			factory.Word(word)
		}
	}
	return factory
}

func (p *CommandParser) failed(err error) *DeferredCommand {
	rtn := &DeferredCommand{}
	rtn.Reject(err)
	return rtn
}
