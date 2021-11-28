package domain

import "context"

// Student struct
type Student struct {
	ID        string `json:"id"` //uuid string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// StudentRepository implements the contract for student repository. We only save and get a student here
type StudentRepository interface {
	SaveStudent(ctx context.Context, student *Student) error
	GetStudent(ctx context.Context, studentID string) (*Student, error)
}

// StudentUseCase implements the contract for student functionalities. GetStudent is used by the app, but
// ListenAndSyncStudentRecord is used to sync the data from the profile service
type StudentUseCase interface {
	GetStudent(ctx context.Context, studentID string) (*Student, error)
	ListenAndSyncStudentRecord(ctx context.Context, student *Student) error
}