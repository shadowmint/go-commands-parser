package cparser

// ErrNotHandled is raised when a command string didn't match any factory.
type ErrNoHandler struct{}

// ErrBadSyntax is raised when a command string fails parsing.
type ErrBadSyntax struct{}

// ErrCommandFailed is raised when a command fails to execute.
type ErrCommandFailed struct{}