//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import (
	"CRUD_Go_Backend/internal/handlers/serviceEntities"
	"context"
)

type StudentPgRepo interface {
	Add(ctx context.Context, studentReq serviceEntities.StudentRequest) (int64, error)
	GetByID(ctx context.Context, studentID int64) (serviceEntities.StudentRequest, error)
	Delete(ctx context.Context, studentID int64) error
	Update(ctx context.Context, studentId int64, studentReq serviceEntities.StudentRequest) error
}
type ClassInfoPgRepo interface {
	Add(ctx context.Context, classInfoReq serviceEntities.ClassInfo) (int64, error)
	GetByStudentID(ctx context.Context, studentId int64) ([]serviceEntities.ClassInfo, error)
	DeleteClassByStudentID(ctx context.Context, studentID int64) error
	Update(ctx context.Context, studentId int64, classInfoReq serviceEntities.ClassInfo) error
}
