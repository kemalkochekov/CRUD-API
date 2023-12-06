//go:build integration
// +build integration

package repository

import (
	"CRUD_Go_Backend/internal/handlers/serviceEntities"
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	"CRUD_Go_Backend/internal/repository/postgres"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestCreateClassInfo(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudentID, err := studentRepo.Add(ctx, testStudentReq)

		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)

		//arrange
		classInfoRepo := NewClassInfoStorage(db.DB)
		testClassInfoReq := serviceEntities.ClassInfo{
			ID:        1,
			StudentID: respStudentID,
			ClassName: "Math",
		}
		//act
		createClassInfoID, err := classInfoRepo.Add(ctx, testClassInfoReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, createClassInfoID)
		//act
		testGetByStudentIDClassInfo, err := classInfoRepo.GetByStudentID(ctx, respStudentID)
		//assert
		require.NoError(t, err)
		require.Equal(t, 1, len(testGetByStudentIDClassInfo))
		assert.Equal(t, createClassInfoID, testGetByStudentIDClassInfo[0].StudentID)
		assert.Equal(t, testClassInfoReq.ClassName, testGetByStudentIDClassInfo[0].ClassName)
	})
	t.Run("Fail for: insert or update on table \"class_info\" violates foreign key constraint", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		classInfoRepo := NewClassInfoStorage(db.DB)
		testClassInfoReq := serviceEntities.ClassInfo{ID: 1, StudentID: 1, ClassName: "Math"}
		//act
		createClassInfoID, err := classInfoRepo.Add(ctx, testClassInfoReq)
		//assert
		require.Error(t, err)
		assert.Negative(t, createClassInfoID)
	})
	t.Run("Not Found", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudentID, err := studentRepo.Add(ctx, testStudentReq)

		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)
		//arrange
		classInfoRepo := NewClassInfoStorage(db.DB)
		testClassInfoReq := serviceEntities.ClassInfo{ID: 1, StudentID: 1, ClassName: "Math"}
		createClassInfoID, err := classInfoRepo.Add(ctx, testClassInfoReq)
		require.NoError(t, err)
		assert.NotZero(t, createClassInfoID)
		nonExistentStudentID := -1
		//act
		getByStudentIDClassInfo, err := classInfoRepo.GetByStudentID(ctx, int64(nonExistentStudentID))
		require.NoError(t, err)
		require.Equal(t, 0, len(getByStudentIDClassInfo))
		assert.Nil(t, err)
	})
}

func TestGetByStudentIDClassInfo(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		respStudentID, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)
		//arrange
		classInfoRepo := NewClassInfoStorage(db.DB)
		testClassInfoReq := serviceEntities.ClassInfo{
			ID:        1,
			StudentID: respStudentID,
			ClassName: "Math",
		}
		//act

		createClassInfoID, err := classInfoRepo.Add(ctx, testClassInfoReq)
		getByStudentIDClassInfo, err := classInfoRepo.GetByStudentID(ctx, respStudentID)
		testClassInfoReq.ID = createClassInfoID
		//assert
		require.NoError(t, err)
		require.Equal(t, 1, len(getByStudentIDClassInfo))
		assert.Equal(t, createClassInfoID, getByStudentIDClassInfo[0].StudentID)
		assert.Equal(t, testClassInfoReq.ClassName, getByStudentIDClassInfo[0].ClassName)
	})
	t.Run("Fail Not Found", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudent, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudent)
		//arrange
		classInfoRepo := NewClassInfoStorage(db.DB)
		nonExistentStudentID := -1
		//act
		getByStudentIDClassInfo, err := classInfoRepo.GetByStudentID(ctx, int64(nonExistentStudentID))
		require.NoError(t, err)
		require.Equal(t, 0, len(getByStudentIDClassInfo))
		assert.Nil(t, err)
	})
}

func TestDeleteClassByStudentIDClassInfo(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudentID, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)

		classInfoRepo := NewClassInfoStorage(db.DB)
		//act
		testClassInfoReq := serviceEntities.ClassInfo{
			ID:        1,
			StudentID: respStudentID,
			ClassName: "Math",
		}
		_, err = classInfoRepo.Add(ctx, testClassInfoReq)
		err = classInfoRepo.DeleteClassByStudentID(ctx, respStudentID)
		//assert
		require.NoError(t, err)
		assert.Nil(t, err)
		//act
		getByStudentIDClassInfo, err := classInfoRepo.GetByStudentID(ctx, respStudentID)
		require.NoError(t, err)
		require.Equal(t, 0, len(getByStudentIDClassInfo))
		assert.Nil(t, err)

	})
	t.Run("Fail Not Found", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		respStudentID, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)

		classInfoRepo := NewClassInfoStorage(db.DB)
		//act
		testClassInfoReq := serviceEntities.ClassInfo{
			ID:        1,
			StudentID: respStudentID,
			ClassName: "Math",
		}
		_, err = classInfoRepo.Add(ctx, testClassInfoReq)
		err = classInfoRepo.DeleteClassByStudentID(ctx, respStudentID-1)
		//assert
		require.Error(t, err)
		assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
	})
}

func TestUpdateClassInfo(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act

		respStudentID, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)

		classInfoRepo := NewClassInfoStorage(db.DB)
		//act
		testClassInfoReq := serviceEntities.ClassInfo{
			ID:        1,
			StudentID: respStudentID,
			ClassName: "Math",
		}
		_, err = classInfoRepo.Add(ctx, testClassInfoReq)
		testClassInfoReq.StudentID = 2
		testClassInfoReq.ClassName = "Computer Science"

		err = classInfoRepo.Update(ctx, respStudentID, testClassInfoReq)
		//assert
		require.NoError(t, err)
		assert.Nil(t, err)
	})
	t.Run("Fail Not Found", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudentID, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudentID)

		classInfoRepo := NewClassInfoStorage(db.DB)
		//act
		testClassInfoReq := serviceEntities.ClassInfo{
			ID:        1,
			StudentID: respStudentID,
			ClassName: "Math",
		}
		_, err = classInfoRepo.Add(ctx, testClassInfoReq)
		err = classInfoRepo.Update(ctx, respStudentID-1, testClassInfoReq)
		//assert
		require.Error(t, err)
		assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
	})
}
