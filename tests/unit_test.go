package tests

import (
	"github.com/rodchenkov-sn/alfalfa/auth"
	"github.com/rodchenkov-sn/alfalfa/common"
	"github.com/rodchenkov-sn/alfalfa/service"
	"testing"
	"time"
)

func TestHashing(t *testing.T) {
	var testVar = "abcdef"
	var _, error1 = service.HashPassword(testVar)
	var secondHash, error2 = service.HashPassword(testVar)
	if error1 != nil || error2 != nil {
		t.Error("HashPassword returned an error")
	}
	if !service.ComparePasswords(secondHash, testVar) {
		t.Errorf("ComparePasswords failed on password %s with hash %s", testVar, secondHash)
	}
}

func TestRepositoryConnection(t *testing.T) {
	var settings = service.ReadSettings("../settings.ini")
	_, err := service.NewRepository(settings.RS)
	if err != nil {
		t.Errorf("NewRepository: Error on repository creation: %s", err)
	}
}

func TestToken(t *testing.T) {
	var settings = service.ReadSettings("../settings.ini")
	repository, _ := service.NewRepository(settings.RS)
	var user = common.UserInfo{"admin", "admin", nil}
	repository.AddUser(user)
	var tokenManager = auth.NewTokenManager(repository, settings.PrivateKey)
	var authInfo = common.Credentials{"admin", "admin"}
	var token, err = tokenManager.GenerateToken(authInfo)
	if err != nil {
		t.Errorf("GenerateToken: Error on token generation: %s", err)
	}
	var _, errorValidation = tokenManager.ValidateToken(token)
	if errorValidation != nil {
		t.Errorf("ValidateToken: Error on token validation: %s", err)
	}
}

func TestRepository(t *testing.T) {
	var settings = service.ReadSettings("../settings.ini")
	repository, _ := service.NewRepository(settings.RS)
	var user = common.UserInfo{"admin", "admin", nil}
	repository.AddUser(user)
	var supervisors []common.Supervisor
	var supervisor = common.Supervisor{"admin"}
	_ = append(supervisors, supervisor)
	var err = repository.AddSupervisor(user.Login,supervisors)
	if err != nil {
		t.Errorf("AddSupervisor: Error on supervisor creation: %s", err)
	}
	var measurment = common.Measurement{
		Temperature: -100,
		Timestamp:   time.Time{},
	}
	err = repository.AddMeasurement("admin", measurment)
	if err != nil {
		t.Errorf("AddMeasurement:Error on adding measurment: %s", err)
	}
	var measurments, errMesur = repository.GetMeasurements("admin")
	if errMesur != nil {
		t.Errorf("GetMeasurements: Error getting measurments: %s", err)
	}
	if measurments[0] != measurment {
		t.Error("GetMeasurements: Measurment not added or lost")
	}
	repository.Disconnect()
}
