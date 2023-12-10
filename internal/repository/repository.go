//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import (
	"CRUD_Go_Backend/internal/handlers/models"
	"context"
)

type StudentPgRepo interface {
	Add(ctx context.Context, studentReq models.StudentRequest) (int64, error)
	GetByID(ctx context.Context, studentID int64) (models.StudentRequest, error)
	Delete(ctx context.Context, studentID int64) error
	Update(ctx context.Context, studentID int64, studentReq models.StudentRequest) error
}
type ClassInfoPgRepo interface {
	Add(ctx context.Context, classInfoReq models.ClassInfo) (int64, error)
	GetByStudentID(ctx context.Context, studentID int64) ([]models.ClassInfo, error)
	DeleteClassByStudentID(ctx context.Context, studentID int64) error
	Update(ctx context.Context, studentID int64, classInfoReq models.ClassInfo) error
}
