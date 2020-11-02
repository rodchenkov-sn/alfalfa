package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rodchenkov-sn/alfalfa/service"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
)

func addUser(repository *service.Repository, writer http.ResponseWriter, request *http.Request) {
	var authInfo service.AuthInfo
	if json.NewDecoder(request.Body).Decode(&authInfo) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := repository.AddUser(authInfo); err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
}

func addMeasurement(repository *service.Repository, writer http.ResponseWriter, request *http.Request) {
	var measurementWithAuth service.MeasurementWithAuth
	if json.NewDecoder(request.Body).Decode(&measurementWithAuth) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := repository.AddMeasurement(measurementWithAuth); err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
}

func getMeasurements(repository *service.Repository, writer http.ResponseWriter, request *http.Request) {
	var authInfo service.AuthInfo
	if json.NewDecoder(request.Body).Decode(&authInfo) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	measurements, err := repository.GetMeasurements(authInfo)
	if err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
	if err := json.NewEncoder(writer).Encode(measurements); err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
}

func readSettings(file string) service.RepositorySettings {
	settings, err := ini.Load(file)
	if err != nil {
		panic(err)
	}
	return service.RepositorySettings{
		Uri: settings.Section("").Key("db_link").String(),
		UsersSettings: service.CollectionSettings{
			Database:   settings.Section("users").Key("db_name").String(),
			Collection: settings.Section("users").Key("collection_name").String(),
		},
		MeasurementsSettings: service.CollectionSettings{
			Database:   settings.Section("measurements").Key("db_name").String(),
			Collection: settings.Section("measurements").Key("collection_name").String(),
		},
	}
}

func main() {

	repository, err := service.NewRepository(readSettings("settings.ini"))

	if err == nil {

		defer repository.Disconnect()

		router := mux.NewRouter()
		router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
			addUser(repository, w, r)
		}).Methods("POST")
		router.HandleFunc("/api/measurements", func(w http.ResponseWriter, r *http.Request) {
			addMeasurement(repository, w, r)
		}).Methods("POST")
		router.HandleFunc("/api/measurements", func(w http.ResponseWriter, r *http.Request) {
			getMeasurements(repository, w, r)
		})
		router.Handle("/home/{rest}",
			http.StripPrefix("/home/", http.FileServer(http.Dir("./static/"))))

		port := os.Getenv("PORT")
		log.Fatal(http.ListenAndServe(":" + port, router))
	}
}

