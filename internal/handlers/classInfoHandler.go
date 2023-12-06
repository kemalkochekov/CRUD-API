package handlers

import (
	"CRUD_Go_Backend/internal/handlers/serviceEntities"
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	"CRUD_Go_Backend/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// ClassInfoHandler handles class information-related HTTP requests
type ClassInfoHandler struct {
	classInfoStorage repository.ClassInfoPgRepo
	queryParamKey    string
}

// NewClassInfoHandler creates a new ClassInfoHandler with the given class information storage service
func NewClassInfoHandler(classInfoStorage repository.ClassInfoPgRepo, queryParamKey string) *ClassInfoHandler {
	return &ClassInfoHandler{
		classInfoStorage: classInfoStorage,
		queryParamKey:    queryParamKey,
	}
}

func (h *ClassInfoHandler) AddClass(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}
	var classInfo serviceEntities.ClassInfo
	if err := json.Unmarshal(body, &classInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal JSON: %v", err), http.StatusBadRequest)
		return
	}
	classInfo.ID, err = h.classInfoStorage.Add(req.Context(), classInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add class_info: %v", err), http.StatusInternalServerError)
		return
	}

	classInfoJson, err := json.Marshal(classInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(classInfoJson)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *ClassInfoHandler) UpdateClass(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}
	var classInfo serviceEntities.ClassInfo // 1
	if err := json.Unmarshal(body, &classInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal JSON: %v", err), http.StatusBadRequest)
		return
	}

	err = h.classInfoStorage.Update(req.Context(), classInfo.StudentID, classInfo)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			http.Error(w, "Cannot update the student in class_info due to existing references (foreign key constraint).", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update class_info: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	message := "Successfully Updated Class Info"
	responseByte := []byte(message)

	_, err = w.Write(responseByte)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *ClassInfoHandler) DeleteClassByStudent(w http.ResponseWriter, req *http.Request) {
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

	err = h.classInfoStorage.DeleteClassByStudentID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			http.Error(w, "Cannot delete the class_info due to non existing studentID", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete record from class_info by StudentByID: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	message := "Successfully Deleted Class Info"
	responseByte := []byte(message)

	_, err = w.Write(responseByte)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *ClassInfoHandler) GetAllClassesByStudent(w http.ResponseWriter, req *http.Request) {
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

	classesInfo, err := h.classInfoStorage.GetByStudentID(req.Context(), keyInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get all records who is attenting that class from class_info: %v", err), http.StatusInternalServerError)
		return
	}

	if len(classesInfo) == 0 {
		// classesInfo is empty
		http.Error(w, "No existing student in class_info", http.StatusNotFound)
		return
	}
	userInfoJson, err := json.Marshal(classesInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(userInfoJson)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
