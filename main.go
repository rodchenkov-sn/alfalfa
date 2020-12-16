package main

import (
	"github.com/gorilla/mux"
	"github.com/rodchenkov-sn/alfalfa/auth"
	"github.com/rodchenkov-sn/alfalfa/network"
	"github.com/rodchenkov-sn/alfalfa/service"
	"log"
	"net/http"
	"os"
)

func main() {

	settings := service.ReadSettings("settings.ini")
	repository, err := service.NewRepository(settings.RS)

	if err == nil {

		defer repository.Disconnect()

		router := mux.NewRouter()

		// registration handlers

		router.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
			network.AddUser(repository, w, r)
		}).Methods("POST")

		tokenManager := auth.NewTokenManager(repository, settings.PrivateKey)

		// measurement handlers

		router.HandleFunc("/api/{login}/measurements", func(w http.ResponseWriter, r *http.Request) {
			network.AddMeasurement(tokenManager, repository, w, r)
		}).Methods("POST")
		router.HandleFunc("/api/{login}/measurements", func(w http.ResponseWriter, r *http.Request) {
			network.GetMeasurements(tokenManager, repository, w, r)
		})

		// supervisors handlers

		router.HandleFunc("/api/{login}/supervisors", func(w http.ResponseWriter, r *http.Request) {
			network.AddSupervisor(tokenManager, repository, w, r)
		}).Methods("POST")

		// authentication handlers

		router.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
			network.AuthUser(tokenManager, w, r)
		})

		// file handlers

		router.Handle("/home/{rest}",
			http.StripPrefix("/home/", http.FileServer(http.Dir("./static/"))))

		port := os.Getenv("PORT")
		log.Fatal(http.ListenAndServe(":" + port, router))
	}
}

