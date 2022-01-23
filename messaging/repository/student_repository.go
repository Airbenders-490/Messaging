package repository

import (
	"chat/domain"
	"github.com/gocql/gocql"
)

type StudentRepository struct {
	dbSession *gocql.Session
}

func NewStudentRepository(session *gocql.Session) *StudentRepository {
	return &StudentRepository{
		dbSession: session,
	}
}

const (
	saveStudent = `INSERT INTO chat.student (student_id, first_name, last_name) VALUES (?,?,?);`
	getStudent  = `SELECT * FROM chat.student WHERE student_id=?;`
)

func (r StudentRepository) SaveStudent(student *domain.Student) error {
	return r.dbSession.Query(saveStudent, student.ID, student.FirstName, student.LastName).Consistency(gocql.One).Exec()
}

func (r StudentRepository) GetStudent(userID string) (*domain.Student, error) {
	var student domain.Student
	err := r.dbSession.Query(getStudent, userID).Consistency(gocql.One).Scan(&student.ID, &student.FirstName, &student.LastName)
	if err != nil {
		return nil, err
	}
	return &student, nil
}
