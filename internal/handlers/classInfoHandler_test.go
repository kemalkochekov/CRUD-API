package handlers

import (
	"CRUD_Go_Backend/internal/handlers/models"
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	mock_repository "CRUD_Go_Backend/internal/repository/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestClassInfoHandler_AddClass(t *testing.T) {
	t.Parallel()
	var (
		queryParamKey = "id"
	)
	type mockExpected struct {
		result int64
		error  error
	}
	tests := []struct {
		description               string
		mockArguments             models.ClassInfo
		mockExpectedEntities      mockExpected
		result                    models.ClassInfo
		expectedCode              int
		expectedHTTPErrorResponce string
	}{
		{
			description:               "Succesfully Added into Database",
			mockArguments:             models.ClassInfo{StudentID: 1, ClassName: "math"},
			mockExpectedEntities:      mockExpected{result: 0, error: nil},
			result:                    models.ClassInfo{StudentID: 1, ClassName: "math"},
			expectedCode:              http.StatusOK,
			expectedHTTPErrorResponce: "",
		},
		{
			description:               "ForeignKey Error",
			mockArguments:             models.ClassInfo{StudentID: 2, ClassName: "math"},
			mockExpectedEntities:      mockExpected{result: -1, error: pkgErrors.ErrForeignKey},
			result:                    models.ClassInfo{},
			expectedCode:              http.StatusInternalServerError,
			expectedHTTPErrorResponce: "Failed to add class_info: ERROR: insert or update on table \"class_info\" violates foreign key constraint \"fk_student\" (SQLSTATE 23503)\n",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			jsonData, err := json.Marshal(tc.mockArguments)
			require.NoError(t, err)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			classInfoHandler := NewClassInfoHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().Add(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodPost, "/class_info", bytes.NewReader(jsonData))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			// act
			classInfoHandler.AddClass(rr, req)
			// assert
			if rr.Code != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.expectedCode)
			}

			if rr.Code != http.StatusOK {
				assert.Equal(t, tc.expectedHTTPErrorResponce, rr.Body.String())
				return
			}

			var actual models.ClassInfo
			err = json.Unmarshal(rr.Body.Bytes(), &actual)
			require.NoError(t, err)
			assert.Equal(t, actual, tc.result)
		})
	}
}

func TestClassInfoHandler_GetAllClassesByStudent(t *testing.T) {
	t.Parallel()
	var (
		queryParamKey = "id"
	)
	type mockExpected struct {
		result []models.ClassInfo
		error  error
	}
	tests := []struct {
		description               string
		mockArguments             int64
		mockExpectedEntities      mockExpected
		result                    []models.ClassInfo
		expectedCode              int
		expectedHTTPErrorResponce string
	}{
		{
			description:               "Succesfully Get ClassInfo By StudentID",
			mockArguments:             1,
			mockExpectedEntities:      mockExpected{result: []models.ClassInfo{{StudentID: 1, ClassName: "math"}}, error: nil},
			result:                    []models.ClassInfo{{StudentID: 1, ClassName: "math"}},
			expectedCode:              http.StatusOK,
			expectedHTTPErrorResponce: "",
		},
		{
			description:               "In ClassInfo with StudentID does not exist",
			mockArguments:             2,
			mockExpectedEntities:      mockExpected{},
			result:                    []models.ClassInfo{},
			expectedCode:              http.StatusNotFound,
			expectedHTTPErrorResponce: "No existing student in class_info\n",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			classInfoHandler := NewClassInfoHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().GetByStudentID(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/class_info/%d", tc.mockArguments), bytes.NewReader([]byte{}))
			require.NoError(t, err)
			req = mux.SetURLVars(req, map[string]string{queryParamKey: strconv.Itoa(int(tc.mockArguments))})
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			// act
			classInfoHandler.GetAllClassesByStudent(rr, req)
			// assert
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}
			if rr.Code != http.StatusOK {
				assert.Equal(t, tc.expectedHTTPErrorResponce, rr.Body.String())
				return
			}
			var actual []models.ClassInfo
			err = json.Unmarshal(rr.Body.Bytes(), &actual)
			require.NoError(t, err)
			assert.Equal(t, actual, tc.result)
		})
	}
}
func TestClassInfoHandler_DeleteClassByStudent(t *testing.T) {
	t.Parallel()
	var (
		queryParamKey = "id"
	)
	tests := []struct {
		description       string
		expectedMessage   string
		mockArguments     int64
		mockExpectedError error
		expectedCode      int
	}{
		{
			description:       "Unable to delete",
			expectedMessage:   "Failed to delete record from class_info by StudentByID: assert.AnError general error for testing\n",
			mockArguments:     4,
			mockExpectedError: assert.AnError,
			expectedCode:      http.StatusInternalServerError,
		},
		{
			description:       "In ClassInfo not found StudentByID",
			expectedMessage:   "Cannot delete the class_info due to non existing studentID\n",
			mockArguments:     4,
			mockExpectedError: pkgErrors.ErrNotFound,
			expectedCode:      http.StatusNotFound,
		},
		{
			description:       "Successfully Deleted StudentID",
			expectedMessage:   "Successfully Deleted Class Info",
			mockArguments:     1,
			mockExpectedError: nil,
			expectedCode:      http.StatusOK,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			classInfoHandler := NewClassInfoHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().DeleteClassByStudentID(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedError)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/class_info/%d", tc.mockArguments), bytes.NewReader([]byte{}))
			require.NoError(t, err)
			req = mux.SetURLVars(req, map[string]string{queryParamKey: strconv.Itoa(int(tc.mockArguments))})
			rr := httptest.NewRecorder()
			// act
			classInfoHandler.DeleteClassByStudent(rr, req)
			// assert
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			if message := rr.Body.String(); message != tc.expectedMessage {
				t.Errorf("handler returned wrong message: got %v want %v", message, tc.expectedMessage)
			}
		})
	}
}
func TestClassInfoHandler_UpdateClass(t *testing.T) {
	t.Parallel()
	var (
		queryParamKey = "id"
	)
	tests := []struct {
		description       string
		expectedMessage   string
		mockArguments     models.ClassInfo
		mockExpectedError error
		expectedCode      int
	}{
		{
			description:       "Successfully Updated in Database",
			expectedMessage:   "Successfully Updated Class Info",
			mockArguments:     models.ClassInfo{StudentID: 1, ClassName: "math"},
			mockExpectedError: nil,
			expectedCode:      http.StatusOK,
		},
		{
			description:       "Not Found",
			expectedMessage:   "Cannot update the student in class_info due to existing references (foreign key constraint).\n",
			mockArguments:     models.ClassInfo{StudentID: 2, ClassName: "math"},
			mockExpectedError: pkgErrors.ErrNotFound,
			expectedCode:      http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			jsonData, err := json.Marshal(tc.mockArguments)
			require.NoError(t, err)
			classInfoHandler := NewClassInfoHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().Update(gomock.Any(), tc.mockArguments.StudentID, tc.mockArguments).Return(tc.mockExpectedError)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodPut, "/class_info", bytes.NewReader(jsonData))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			// act
			classInfoHandler.UpdateClass(rr, req)
			// assert
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			if message := rr.Body.String(); message != tc.expectedMessage {
				t.Errorf("handler returned wrong message: got %v want %v", message, tc.expectedMessage)
			}
		})
	}
}
