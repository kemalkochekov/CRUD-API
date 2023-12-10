package models

type ClassInfo struct {
	ID        int64  `json:"id"`
	StudentID int64  `json:"student_id"`
	ClassName string `json:"class_name"`
}
