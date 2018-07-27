package snapshot

import (
	"io"
)

// DBSnapshot Represents a database in the system
type DBSnapshot interface {
	Load(reader io.Reader) (map[string]interface{}, error)
	Save(writer io.Writer) error
}
