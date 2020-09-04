package router

import (
	"challenge-junior/middleware"

	"github.com/gorilla/mux"
)

// Router vai ser exportado e usado la no main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/project/{id}", middleware.GetProject).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/project", middleware.GetAllProject).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newproject", middleware.CreateProject).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/project/{id}", middleware.UpdateProject).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteproject/{id}", middleware.DeleteProject).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/lot/{lotid}", middleware.GetLot).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/lot/{projectID}", middleware.GetAllLot).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newlot", middleware.CreateLot).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletelot/{lotid}", middleware.DeleteLot).Methods("DELETE", "OPTIONS")

	return router
}
