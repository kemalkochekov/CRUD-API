package entities

import "CRUD_Go_Backend/internal/handlers/models"

type ClassInfo struct {
	ID        int64  `db:"id"`
	StudentID int64  `db:"student_id"`
	ClassName string `db:"class_name"`
}

func (c *ClassInfo) ToClassInfoDomain() models.ClassInfo {
	return models.ClassInfo{
		ID:        c.ID,
		StudentID: c.StudentID,
		ClassName: c.ClassName,
	}
}
