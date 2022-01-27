package cassandra

import (
	"github.com/gocql/gocql"
)

// SessionInterface is a Cassandra connection. It allows mock of gocql.Session
type SessionInterface interface {
	// Close closes the SessionInterface.
	Close()

	// ExecuteBatch executes a BatchInterface for the SessionInterface
	ExecuteBatch(batch BatchInterface) error

	// NewBatch returns a new BatchInterface for the SessionInterface.
	NewBatch(kind BatchKind) BatchInterface

	// Query returns a new QueryInterface for the SessionInterface.
	Query(string, ...interface{}) QueryInterface
}

// NewSession returns a new Session for s.
func NewSession(s *gocql.Session) SessionInterface {
	return &session{s: s}
}

// session is a gocql.Session wrapper
type session struct {
	s *gocql.Session
	b *Batch
}

// Close wraps the session's Close method
func (s *session) Close() {
	s.s.Close()
}

// ExecuteBatch wraps the session's ExecuteBatch method
func (s *session) ExecuteBatch(batch BatchInterface) error {
	return s.s.ExecuteBatch(s.b.B)
}

// NewBatch wraps the session's NewBatch method
func (s *session) NewBatch(kind BatchKind) BatchInterface {
	s.b = &Batch{B: s.s.NewBatch(gocql.BatchType(kind)), s: s.s}
	return s.b
}

// Query wraps the session's Query method
func (s *session) Query(stmt string, values ...interface{}) QueryInterface {
	return NewQuery(s.s.Query(stmt, values...))
}
