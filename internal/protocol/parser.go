package protocol

// CommandParser is an interface that is used to parse and retrive a command from input
type CommandParser interface {
	// Parse and return a command or an error
	Parse() (*Command, error)
}
