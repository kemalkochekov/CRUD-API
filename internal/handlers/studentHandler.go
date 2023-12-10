package handlers

import (
	"CRUD_Go_Backend/internal/handlers/models"
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	"CRUD_Go_Backend/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// StudentHandler handles student-related HTTP requests.
type StudentHandler struct {
	studentStorage repository.StudentPgRepo
	queryParamKey  string
}

// NewStudentHandler creates a new StudentHandler with the given student storage service.
func NewStudentHandler(studentStorage repository.StudentPgRepo, queryParamKey string) *StudentHandler {
	return &StudentHandler{
		studentStorage: studentStorage,
		queryParamKey:  queryParamKey,
	}
}

func (h *StudentHandler) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}

	var studentReq models.StudentRequest

	if err = json.Unmarshal(body, &studentReq); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal JSON: %v", err), http.StatusBadRequest)
		return
	}

	if studentReq.StudentName == "" || studentReq.Grade < 0 {
		http.Error(w, fmt.Sprintf("Failed Student name is empty or Grade is negative: %v", err), http.StatusBadRequest)
		return
	}

	studentReq.StudentID, err = h.studentStorage.Add(req.Context(), studentReq)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrInvalidName) {
			http.Error(w, fmt.Sprintf("Name should not be empty: %v", err), http.StatusBadRequest)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to add student: %v", err), http.StatusInternalServerError)

		return
	}

	userInfoJSON, err := json.Marshal(studentReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(userInfoJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *StudentHandler) Update(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}

	var student models.StudentRequest // 1
	if err := json.Unmarshal(body, &student); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal JSON: %v", err), http.StatusBadRequest)
		return
	}

	err = h.studentStorage.Update(req.Context(), student.StudentID, student)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			http.Error(w, "Student with such student_id not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to update studentByID: %v", err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	message := "Successfully Updated Student Info"
	responseByte := []byte(message)

	_, err = w.Write(responseByte)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *StudentHandler) Get(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[h.queryParamKey]
	if !ok {
		http.Error(w, "Invalid request. Missing query parameter.", http.StatusBadRequest)
		return
	}

	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Failed to convert string to int64.", http.StatusBadRequest)
		return
	}

	userInfo, err := h.studentStorage.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to get record by StudentByID: %v", err), http.StatusInternalServerError)

		return
	}

	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(userInfoJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *StudentHandler) Delete(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[h.queryParamKey]
	if !ok {
		http.Error(w, "Invalid request. Missing query parameter.", http.StatusBadRequest)
		return
	}

	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Failed to convert string to int64.", http.StatusBadRequest)
		return
	}

	err = h.studentStorage.Delete(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			http.Error(w, "Student with such student_id not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to delete studentByID: %v", err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	message := "Successfully Deleted Student Info"
	responseByte := []byte(message)

	_, err = w.Write(responseByte)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
