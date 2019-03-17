package protocol

type CommandParser interface {
	Parse() (*Command, error)
}
