package cassandra

import "github.com/gocql/gocql"

type ScannerInterface interface {
	Next() bool
	Scan(...interface{}) error
	Err() error
}

type Scanner struct {
	scanner gocql.Scanner
}

func NewScanner(scanner gocql.Scanner) ScannerInterface {
	return &Scanner{
		scanner,
	}
}

func (s *Scanner) Next() bool {
	return s.scanner.Next()
}

func (s *Scanner) Scan(dest ...interface{}) error {
	return s.scanner.Scan(dest...)
}

func (s *Scanner) Err() error {
	return s.scanner.Err()
}





