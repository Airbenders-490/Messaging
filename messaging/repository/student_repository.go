package repository

import (
	"chat/domain"
	"chat/messaging/repository/cassandra"
	"context"
	"github.com/gocql/gocql"
)

type StudentRepository struct {
	dbSession cassandra.SessionInterface
}

func NewStudentRepository(session cassandra.SessionInterface) *StudentRepository {
	return &StudentRepository{
		dbSession: session,
	}
}

const (
	saveStudent = `INSERT INTO chat.student (student_id, first_name, last_name) VALUES (?,?,?);`
	getStudent  = `SELECT * FROM chat.student WHERE student_id=?;`
)

func (r StudentRepository) SaveStudent(ctx context.Context, student *domain.Student) error {
	return r.dbSession.Query(saveStudent, student.ID, student.FirstName, student.LastName).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r StudentRepository) GetStudent(ctx context.Context, userID string) (*domain.Student, error) {
	var student domain.Student
	err := r.dbSession.Query(getStudent, userID).WithContext(ctx).Consistency(gocql.One).Scan(&student.ID, &student.FirstName, &student.LastName)
	if err != nil {
		return nil, err
	}
	return &student, nil
}
