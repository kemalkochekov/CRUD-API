package handlers

import (
	"CRUD_Go_Backend/internal/repository"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(studentStorage repository.StudentPgRepo, classInfoStorage repository.ClassInfoPgRepo, queryParamKey string) *mux.Router {
	router := mux.NewRouter()

	studentHandler := NewStudentHandler(studentStorage, queryParamKey)
	classInfoHandler := NewClassInfoHandler(classInfoStorage, queryParamKey)

	// Handler for student
	router.HandleFunc("/student", studentHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/student", studentHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(fmt.Sprintf("/student/{%s:[0-9]+}", queryParamKey), studentHandler.Get).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/student/{%s:[0-9]+}", queryParamKey), studentHandler.Delete).Methods(http.MethodDelete)

	// Handler for class_info
	router.HandleFunc("/class_info", classInfoHandler.AddClass).Methods(http.MethodPost)
	router.HandleFunc("/class_info", classInfoHandler.UpdateClass).Methods(http.MethodPut)
	router.HandleFunc(fmt.Sprintf("/class_info/{%s:[0-9]+}", queryParamKey), classInfoHandler.DeleteClassByStudent).Methods(http.MethodDelete)
	router.HandleFunc(fmt.Sprintf("/class_info/{%s:[0-9]+}", queryParamKey), classInfoHandler.GetAllClassesByStudent).Methods(http.MethodGet)

	return router
}
