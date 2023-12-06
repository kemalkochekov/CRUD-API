package handlers

import (
	"net/http"
)

// StudentHandlerInterface defines the methods required for handling student-related requests.
type StudentHandlerInterface interface {
	Create(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

// ClassInfoHandlerInterface defines the methods required for handling class information-related requests.
type ClassInfoHandlerInterface interface {
	AddClass(w http.ResponseWriter, req *http.Request)
	UpdateClass(w http.ResponseWriter, req *http.Request)
	DeleteClassByStudent(w http.ResponseWriter, req *http.Request)
	GetAllClassesByStudent(w http.ResponseWriter, req *http.Request)
}
