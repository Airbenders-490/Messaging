package cassandra

import (
	"context"
	"github.com/gocql/gocql"
)

// BatchInterface is an ordered collection of CQL queries. It allows mock of gocql.Batch
type BatchInterface interface {
	// AddBatchEntry adds a gocql.BatchEntry for the BatchInterface.
	AddBatchEntry(entry *gocql.BatchEntry)

	// WithContext returns a BatchInterface for the BatchInterface.
	WithContext(ctx context.Context) BatchInterface
}

// BatchKind is the kind of Batch. The choice of kind mostly affects performance.
type BatchKind byte

// Kinds of batches.
const (
	// BatchLogged queries are atomic. Queries are only isolated within a single partition.
	BatchLogged BatchKind = 0

	// BatchUnlogged queries are not atomic. Atomic queries spanning multiple partitions cost performance.
	BatchUnlogged BatchKind = 1

	// BatchCounter queries update counters and are not idempotent.
	BatchCounter BatchKind = 2
)

// Batch is a gocql specific batch
type Batch struct {
	B *gocql.Batch
	s *gocql.Session
}

// AddBatchEntry encapsulates the adding of batch entries for better mockability
func (b *Batch) AddBatchEntry(entry *gocql.BatchEntry) {
	b.B.Entries = append(b.B.Entries, *entry)
}

// WithContext wraps the batch's WithContext method
func (b *Batch) WithContext(ctx context.Context) BatchInterface {
	b.B.WithContext(ctx)
	return b
}