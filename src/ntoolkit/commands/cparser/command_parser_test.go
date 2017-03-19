package cparser_test

import (
	"testing"

	"ntoolkit/assert"
	"ntoolkit/commands"
	"ntoolkit/commands/cparser"
	"ntoolkit/errors"
)

func fixture() *cparser.CommandParser {
	rtn := cparser.New()

	// Syntax handlers
	rtn.Register(&GoCommandFactory{})
	rtn.Register(&LookCommandFactory{})
	rtn.Register(&UseCommandFactory{})
	registerUse2Factory(rtn)
	registerPutFactory(rtn)

	// Command handlers
	rtn.Commands.Register(&GoCommandHandler{})
	rtn.Commands.Register(&LookCommandHandler{})
	rtn.Commands.Register(&UseCommandHandler{})
	rtn.Commands.Register(&Use2CommandHandler{})
	rtn.Commands.Register(&PutCommandHandler{})

	return rtn
}

func TestNew(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		T.Assert(cparser.New() != nil)
	})
}

func TestWithNoHandlerCommandsFail(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := cparser.New()
		cmd, err := p.Wait("Hello world", nil)
		T.Assert(cmd == nil)
		T.Assert(err != nil)
		T.Assert(errors.Is(err, cparser.ErrNoHandler{}))
	})
}

func TestWaitForCommand(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := fixture()

		cmd, err := p.Wait("go north", nil)
		T.Assert(err == nil)
		T.Assert(cmd != nil)
		gcmd, ok := cmd.(*GoCommand)
		T.Assert(ok)
		T.Assert(gcmd.Success == true)
		T.Assert(gcmd.Direction == "north")

		cmd, err = p.Wait("go \"north road\"", nil)
		T.Assert(err == nil)
		T.Assert(cmd != nil)
		gcmd, ok = cmd.(*GoCommand)
		T.Assert(ok)
		T.Assert(gcmd.Success == true)
		T.Assert(gcmd.Direction == "north road")

		cmd, err = p.Wait("look north", nil)
		lcmd, ok := cmd.(*LookCommand)
		T.Assert(ok)
		T.Assert(lcmd.Success == true)
		T.Assert(lcmd.Direction == "north")
		T.Assert(cmd != nil)
		T.Assert(err == nil)
	})
}

func TestDeferredCommand(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := fixture()

		p.Execute("go north", nil).Then(func(cmd commands.Command) {
			gcmd, ok := cmd.(*GoCommand)
			T.Assert(ok)
			T.Assert(gcmd.Success == true)
			T.Assert(gcmd.Direction == "north")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("look north", nil).Then(func(cmd commands.Command) {
			lcmd, ok := cmd.(*LookCommand)
			T.Assert(ok)
			T.Assert(lcmd.Success == true)
			T.Assert(lcmd.Direction == "north")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("fadsfdasf north", nil).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			T.Assert(err != nil)
			T.Assert(errors.Is(err, cparser.ErrNoHandler{}))
		})
	})
}

func TestComplexCommand(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := fixture()

		p.Execute("use hammer on door", nil).Then(func(cmd commands.Command) {
			ucmd, ok := cmd.(*UseCommand)
			T.Assert(ok)
			T.Assert(ucmd.Tool == "hammer")
			T.Assert(ucmd.Target == "door")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("use spoon on door", nil).Then(func(cmd commands.Command) {
			ucmd, ok := cmd.(*UseCommand)
			T.Assert(ok)
			T.Assert(ucmd.Tool == "spoon")
			T.Assert(ucmd.Target == "door")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("use hammer on door", nil).Then(func(cmd commands.Command) {
			ucmd, ok := cmd.(*UseCommand)
			T.Assert(ok)
			T.Assert(ucmd.Tool == "hammer")
			T.Assert(ucmd.Target == "door")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("use hammer on spoon", nil).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			inner, _ := errors.Inner(err)
			T.Assert(errors.Is(inner, UseCommandErrBadUse{}))
		})

		p.Execute("use \"magic hammer\" on door", nil).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			inner, _ := errors.Inner(err)
			T.Assert(errors.Is(inner, UseCommandErrBadTool{}))
		})
	})
}

func TestStandardCommandFactory(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := fixture()

		p.Execute("use2 hammer on door", nil).Then(func(cmd commands.Command) {
			ucmd, ok := cmd.(*Use2Command)
			T.Assert(ok)
			T.Assert(ucmd.Tool == "hammer")
			T.Assert(ucmd.Target == "door")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("use2 spoon on door", nil).Then(func(cmd commands.Command) {
			ucmd, ok := cmd.(*Use2Command)
			T.Assert(ok)
			T.Assert(ucmd.Tool == "spoon")
			T.Assert(ucmd.Target == "door")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("use2 hammer on door", nil).Then(func(cmd commands.Command) {
			ucmd, ok := cmd.(*Use2Command)
			T.Assert(ok)
			T.Assert(ucmd.Tool == "hammer")
			T.Assert(ucmd.Target == "door")
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("use2 hammer on spoon", nil).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			inner, _ := errors.Inner(err)
			T.Assert(errors.Is(inner, UseCommandErrBadUse{}))
		})

		p.Execute("use2 \"magic hammer\" on door", nil).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			inner, _ := errors.Inner(err)
			T.Assert(errors.Is(inner, UseCommandErrBadTool{}))
		})
	})
}

func TestMultipleCommandSyntax(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := fixture()

		playerId := 123132

		p.Execute("put foo on bar", playerId).Then(func(cmd commands.Command) {
			pcmd, ok := cmd.(*PutCommand)
			T.Assert(ok)
			T.Assert(pcmd.Item == "foo")
			T.Assert(pcmd.Target == "bar")
			T.Assert(pcmd.PlayerId == playerId)
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("put foo in bar", playerId).Then(func(cmd commands.Command) {
			pcmd, ok := cmd.(*PutCommand)
			T.Assert(ok)
			T.Assert(pcmd.Item == "foo")
			T.Assert(pcmd.Container == "bar")
			T.Assert(pcmd.PlayerId == playerId)
		}, func(err error) {
			T.Unreachable()
		})

		p.Execute("put foo under bar", playerId).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			T.Assert(err != nil)
		})

		p.Execute("put foo adfdasf", playerId).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			T.Assert(err != nil)
		})

		p.Execute("put dragon on bar", playerId).Then(func(cmd commands.Command) {
			T.Unreachable()
		}, func(err error) {
			inner, ok := errors.Inner(err)
			T.Assert(ok)
			T.Assert(errors.Is(inner, ErrInvalidDragon{}))
		})
	})
}