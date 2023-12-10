package handlers

import (
	"CRUD_Go_Backend/internal/repository"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(
	studentStorage repository.StudentPgRepo,
	classInfoStorage repository.ClassInfoPgRepo,
	queryParamKey string,
) *mux.Router {
	router := mux.NewRouter()

	studentHandler := NewStudentHandler(studentStorage, queryParamKey)
	classInfoHandler := NewClassInfoHandler(classInfoStorage, queryParamKey)

	// Main Page to check
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("WELCOME CRUD GO BACKEND"))
		if err != nil {
			http.Error(writer, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	// Handler for student
	router.HandleFunc("/student", studentHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/student", studentHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(fmt.Sprintf("/student/{%s:[0-9]+}", queryParamKey), studentHandler.Get).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/student/{%s:[0-9]+}", queryParamKey), studentHandler.Delete).Methods(http.MethodDelete)

	// Handler for class_info
	router.HandleFunc("/class_info", classInfoHandler.AddClass).Methods(http.MethodPost)
	router.HandleFunc("/class_info", classInfoHandler.UpdateClass).Methods(http.MethodPut)
	router.HandleFunc(
		fmt.Sprintf("/class_info/{%s:[0-9]+}", queryParamKey),
		classInfoHandler.DeleteClassByStudent,
	).Methods(http.MethodDelete)
	router.HandleFunc(
		fmt.Sprintf("/class_info/{%s:[0-9]+}", queryParamKey),
		classInfoHandler.GetAllClassesByStudent,
	).Methods(http.MethodGet)

	return router
}
