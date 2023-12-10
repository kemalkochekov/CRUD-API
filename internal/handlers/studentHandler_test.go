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

func TestStudentHandler_Get(t *testing.T) {
	t.Parallel()
	var (
		queryParamKey = "id"
	)
	type mockExpected struct {
		result models.StudentRequest
		error  error
	}
	tests := []struct {
		description               string
		mockArguments             int64
		mockExpectedEntities      mockExpected
		result                    models.StudentRequest
		expectedCode              int
		expectedHTTPErrorResponse string
	}{
		{
			description:               "Student not found",
			mockArguments:             4,
			mockExpectedEntities:      mockExpected{error: pkgErrors.ErrNotFound},
			result:                    models.StudentRequest{},
			expectedCode:              http.StatusNotFound,
			expectedHTTPErrorResponse: "Student not found\n",
		},
		{
			description:               "Student exists",
			mockArguments:             1,
			mockExpectedEntities:      mockExpected{result: models.StudentRequest{StudentID: 1, StudentName: "Test", Grade: 90}, error: nil},
			result:                    models.StudentRequest{StudentID: 1, StudentName: "Test", Grade: 90},
			expectedCode:              http.StatusOK,
			expectedHTTPErrorResponse: "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			studentHandler := NewStudentHandler(mockRepo, queryParamKey)

			mockRepo.EXPECT().GetByID(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/student/%d", tc.mockArguments), bytes.NewReader([]byte{}))
			require.NoError(t, err)
			req = mux.SetURLVars(req, map[string]string{queryParamKey: strconv.Itoa(int(tc.mockArguments))})
			rr := httptest.NewRecorder()
			// act
			studentHandler.Get(rr, req)
			// assert
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}
			if rr.Code != http.StatusOK {
				assert.Equal(t, tc.expectedHTTPErrorResponse, rr.Body.String())
				return
			}
			var actual models.StudentRequest
			err = json.Unmarshal(rr.Body.Bytes(), &actual)
			require.NoError(t, err)
			assert.Equal(t, actual, tc.result)
		})
	}
}

func TestStudentHandler_Create(t *testing.T) {
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
		mockArguments             models.StudentRequest
		mockExpectedEntities      mockExpected
		result                    models.StudentRequest
		expectedCode              int
		expectedHTTPErrorResponse string
	}{
		{
			description:               "Successfully Added into Database",
			mockArguments:             models.StudentRequest{StudentID: 1, StudentName: "Test", Grade: 90},
			mockExpectedEntities:      mockExpected{result: 1, error: nil},
			result:                    models.StudentRequest{StudentID: 1, StudentName: "Test", Grade: 90},
			expectedCode:              http.StatusOK,
			expectedHTTPErrorResponse: "",
		},
		{
			description:               "Failed database unable to add",
			mockArguments:             models.StudentRequest{StudentID: 1, StudentName: "test", Grade: 98},
			mockExpectedEntities:      mockExpected{result: 1, error: assert.AnError},
			result:                    models.StudentRequest{},
			expectedCode:              http.StatusInternalServerError,
			expectedHTTPErrorResponse: "Failed to add student: assert.AnError general error for testing\n",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			jsonData, err := json.Marshal(tc.mockArguments)
			require.NoError(t, err)
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			studentHandler := NewStudentHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().Add(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)
			defer ctrl.Finish()

			req, err := http.NewRequest(http.MethodPost, "/student", bytes.NewReader(jsonData))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			// act
			studentHandler.Create(rr, req)
			// assert
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}
			if rr.Code != http.StatusOK {
				assert.Equal(t, tc.expectedHTTPErrorResponse, rr.Body.String())
				return
			}

			var actual models.StudentRequest
			err = json.Unmarshal(rr.Body.Bytes(), &actual)
			require.NoError(t, err)
			assert.Equal(t, actual, tc.result)
		})
	}
}

func TestStudentHandler_Delete(t *testing.T) {
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
			expectedMessage:   "Failed to delete studentByID: assert.AnError general error for testing\n",
			mockArguments:     4,
			mockExpectedError: assert.AnError,
			expectedCode:      http.StatusInternalServerError,
		},
		{
			description:       "Student not found",
			expectedMessage:   "Student with such student_id not found\n",
			mockArguments:     4,
			mockExpectedError: pkgErrors.ErrNotFound,
			expectedCode:      http.StatusNotFound,
		},
		{
			description:       "Successfully Deleted StudentID",
			expectedMessage:   "Successfully Deleted Student Info",
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
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			studentHandler := NewStudentHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().Delete(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedError)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/student/%d", tc.mockArguments), bytes.NewReader([]byte{}))
			require.NoError(t, err)
			req = mux.SetURLVars(req, map[string]string{queryParamKey: strconv.Itoa(int(tc.mockArguments))})
			rr := httptest.NewRecorder()
			// act
			studentHandler.Delete(rr, req)
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

func TestStudentHandler_Update(t *testing.T) {
	t.Parallel()
	var (
		queryParamKey = "id"
	)
	tests := []struct {
		description       string
		expectedMessage   string
		mockArguments     models.StudentRequest
		mockExpectedError error
		expectedCode      int
	}{
		{
			description:       "Succesfully Updated in Database",
			expectedMessage:   "Successfully Updated Student Info",
			mockArguments:     models.StudentRequest{StudentID: 0, StudentName: "Test2", Grade: 92},
			mockExpectedError: nil,
			expectedCode:      http.StatusOK,
		},
		{
			description:       "Not Found",
			expectedMessage:   "Student with such student_id not found\n",
			mockArguments:     models.StudentRequest{StudentID: 0, StudentName: "Test2", Grade: 92},
			mockExpectedError: pkgErrors.ErrNotFound,
			expectedCode:      http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			jsonData, err := json.Marshal(tc.mockArguments)
			require.NoError(t, err)
			studentHandler := NewStudentHandler(mockRepo, queryParamKey)
			mockRepo.EXPECT().Update(gomock.Any(), tc.mockArguments.StudentID, tc.mockArguments).Return(tc.mockExpectedError)
			defer ctrl.Finish()
			req, err := http.NewRequest(http.MethodPut, "/student", bytes.NewReader(jsonData))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			// act
			studentHandler.Update(rr, req)
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
