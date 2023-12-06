package entities

import (
	"CRUD_Go_Backend/internal/handlers/serviceEntities"
	"time"
)

type Student struct {
	StudentID   int64     `db:"student_id"`
	StudentName string    `db:"student_name"`
	Grade       int64     `db:"grade"`
	CreatedAt   time.Time `db:"created_at"`
}

func (s *Student) ToStudentDomain() serviceEntities.StudentRequest {
	return serviceEntities.StudentRequest{
		StudentID:   s.StudentID,
		StudentName: s.StudentName,
		Grade:       s.Grade,
	}
}
