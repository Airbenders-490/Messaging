package cassandra

import (
	"context"
	"github.com/gocql/gocql"
)

// QueryInterface allows mock of gocql.Query
type QueryInterface interface {

	// Consistency returns a QueryInterface for the QueryInterface.
	Consistency(c gocql.Consistency) QueryInterface

	// Exec returns an error for the QueryInterface.
	Exec() error

	// Scan returns an error for the QueryInterface.
	Scan(...interface{}) error

	// WithContext returns a QueryInterface for the QueryInterface.
	WithContext(ctx context.Context) QueryInterface
}

// Query is a wrapper for a query for mockability.
type Query struct {
	query *gocql.Query
}

// NewQuery instantiates a new Query
func NewQuery(query *gocql.Query) QueryInterface {
	return &Query{
		query,
	}
}

// Consistency wraps the query's Consistency method
func (q *Query) Consistency(c gocql.Consistency) QueryInterface {
	q.query.Consistency(c)
	return q
}

// Exec wraps the query's Exec method
func (q *Query) Exec() error {
	return q.query.Exec()
}

// Scan wraps the query's Scan method
func (q *Query) Scan(dest ...interface{}) error {
	return q.query.Scan(dest...)
}

// WithContext wraps the query's WithContext method
func (q *Query) WithContext(ctx context.Context) QueryInterface {
	q.query.WithContext(ctx)
	return q
}