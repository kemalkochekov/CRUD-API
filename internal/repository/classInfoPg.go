package repository

import (
	"CRUD_Go_Backend/internal/pkg/connection"
	"context"

	"CRUD_Go_Backend/internal/handlers/models"
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	"CRUD_Go_Backend/internal/pkg/utils"
	"CRUD_Go_Backend/internal/repository/entities"
)

type ClassInfoStorage struct {
	db connection.DBops
}

func NewClassInfoStorage(database connection.DBops) ClassInfoStorage {
	return ClassInfoStorage{db: database}
}
func ToClassInfoStorage(c models.ClassInfo) entities.ClassInfo {
	return entities.ClassInfo{
		StudentID: c.StudentID,
		ClassName: c.ClassName,
	}
}

func (r *ClassInfoStorage) Add(ctx context.Context, classInfoReq models.ClassInfo) (int64, error) {
	classInfoPg := ToClassInfoStorage(classInfoReq)
	var id int64

	err := r.db.ExecQueryRow(ctx, `INSERT INTO class_info(student_id, class_name) VALUES($1, $2) RETURNING id;`,
		classInfoPg.StudentID,
		classInfoPg.ClassName,
	).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *ClassInfoStorage) GetByStudentID(ctx context.Context, studentID int64) ([]models.ClassInfo, error) {
	var classInfo []entities.ClassInfo

	rows, err := r.db.ExecQuery(ctx, `SELECT id, student_id, class_name FROM class_info WHERE student_id=$1;`, studentID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tempClassInfo entities.ClassInfo

		err := rows.Scan(&tempClassInfo.ID, &tempClassInfo.StudentID, &tempClassInfo.ClassName)
		if err != nil {
			return nil, err
		}

		classInfo = append(classInfo, tempClassInfo)
	}

	classesInfo := utils.Map(
		classInfo,
		func(p entities.ClassInfo) models.ClassInfo {
			return p.ToClassInfoDomain()
		},
	)

	return classesInfo, nil
}

func (r *ClassInfoStorage) DeleteClassByStudentID(ctx context.Context, studentID int64) error {
	command, err := r.db.Exec(ctx, "DELETE FROM class_info WHERE student_id = $1", studentID)
	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return pkgErrors.ErrNotFound
	}

	return nil
}

func (r *ClassInfoStorage) Update(ctx context.Context, studentID int64, classInfoReq models.ClassInfo) error {
	classInfo := ToClassInfoStorage(classInfoReq)

	command, err := r.db.Exec(ctx, `
		UPDATE class_info
		SET class_name = $2
		WHERE student_id = $1
	`, studentID, classInfo.ClassName)

	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return pkgErrors.ErrNotFound
	}

	return nil
}
