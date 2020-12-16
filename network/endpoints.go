package network

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rodchenkov-sn/alfalfa/auth"
	"github.com/rodchenkov-sn/alfalfa/common"
	"github.com/rodchenkov-sn/alfalfa/service"
	"log"
	"net/http"
)

func AddUser(repository *service.Repository, writer http.ResponseWriter, request *http.Request) {
	var authInfo common.UserInfo
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

func AddMeasurement(tokenManager *auth.TokenManager, repository *service.Repository,
	writer http.ResponseWriter, request *http.Request) {
	var measurement common.Measurement
	if json.NewDecoder(request.Body).Decode(&measurement) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token := request.Header.Get("Bearer")
	issuer, err := tokenManager.ValidateToken(token)
	if err != nil {
		writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	subject := mux.Vars(request)["login"]
	rights := repository.GetRights(issuer, subject)
	if !rights.Write {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := repository.AddMeasurement(subject, measurement); err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		log.Println(err)
		return
	}
}

func GetMeasurements(tokenManager *auth.TokenManager, repository *service.Repository,
	writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Bearer")
	issuer, err := tokenManager.ValidateToken(token)
	if err != nil {
		writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	subject := mux.Vars(request)["login"]
	rights := repository.GetRights(issuer, subject)
	if !rights.Read {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	measurements, err := repository.GetMeasurements(subject)
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

func AuthUser(tokenManager *auth.TokenManager, writer http.ResponseWriter, request *http.Request) {
	var authInfo common.Credentials
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

func AddSupervisor(tokenManager *auth.TokenManager, repository *service.Repository,
	writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Bearer")
	issuer, err := tokenManager.ValidateToken(token)
	if err != nil {
		writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	var supervisors []common.Supervisor
	if json.NewDecoder(request.Body).Decode(&supervisors) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if repository.AddSupervisor(issuer, supervisors) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}
