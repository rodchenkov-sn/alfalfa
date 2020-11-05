package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rodchenkov-sn/alfalfa/auth"
	"github.com/rodchenkov-sn/alfalfa/common"
	"github.com/rodchenkov-sn/alfalfa/service"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
)

func addUser(repository *service.Repository, writer http.ResponseWriter, request *http.Request) {
	var authInfo common.AuthInfo
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

func addMeasurement(tokenManager *auth.TokenManager, repository *service.Repository,
					writer http.ResponseWriter, request *http.Request) {
	var measurement common.Measurement
	if json.NewDecoder(request.Body).Decode(&measurement) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token := request.Header.Get("Bearer")
	login, err := tokenManager.ValidateToken(token)
	if err != nil {
		writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	if err := repository.AddMeasurement(login, measurement); err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
}

func getMeasurements(tokenManager *auth.TokenManager, repository *service.Repository,
	                 writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Bearer")
	login, err := tokenManager.ValidateToken(token)
	if err != nil {
		writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	measurements, err := repository.GetMeasurements(login)
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

func authUser(tokenManager *auth.TokenManager, writer http.ResponseWriter, request *http.Request) {
	var authInfo common.AuthInfo
	if json.NewDecoder(request.Body).Decode(&authInfo) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := tokenManager.GenerateToken(authInfo)
	if err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
	if _, err := writer.Write([]byte(token)); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
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

		router.HandleFunc("/api/users/register", func(w http.ResponseWriter, r *http.Request) {
			addUser(repository, w, r)
		}).Methods("POST")

		tokenManager := auth.NewTokenManager(repository, "secret_key")

		router.HandleFunc("/api/measurements", func(w http.ResponseWriter, r *http.Request) {
			addMeasurement(tokenManager, repository, w, r)
		}).Methods("POST")
		router.HandleFunc("/api/measurements", func(w http.ResponseWriter, r *http.Request) {
			getMeasurements(tokenManager, repository, w, r)
		})
		router.HandleFunc("/api/users/auth", func(w http.ResponseWriter, r *http.Request) {
			authUser(tokenManager, w, r)
		})

		router.Handle("/home/{rest}",
			http.StripPrefix("/home/", http.FileServer(http.Dir("./static/"))))

		// port := os.Getenv("PORT")
		log.Fatal(http.ListenAndServe(":8080", router))
	}
}

