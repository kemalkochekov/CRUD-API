package serviceEntities

type StudentRequest struct {
	StudentID   int64  `json:"student_id"`
	StudentName string `json:"student_name"`
	Grade       int64  `json:"grade"`
}
