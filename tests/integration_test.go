package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)
var base = "https://alfalfa-project.herokuapp.com/api"
func authenticate() string {
	var jsonStr = []byte(`{
	"login": "abcdef",
	"password": "abcdef"
}`)
	req, _ := http.NewRequest("POST", base + "/auth", bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	token, _ := ioutil.ReadAll(resp.Body)
	return string(token)
}
func TestLogin(t *testing.T) {
	var jsonStr = []byte(`{
	"login": "abcdef",
	"password": "abcdef"
}`)
	req, err := http.NewRequest("POST", base + "/auth", bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error posting request: %s", err	)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("auth returned status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestPostMeasurements(t *testing.T) {
	client := &http.Client{}
	token := authenticate()
	var jsonStr2 = []byte(`{   "temperature": 36.6,   "timestamp": "2012-04-23T18:25:43.511Z" }`)
	req2, err2 := http.NewRequest("POST", base + "/abcdef/measurements", bytes.NewBuffer(jsonStr2))
	req2.Header.Set("Bearer", string(token))
	resp2, err2 := client.Do(req2)
	if err2 != nil {
		t.Errorf("Error posting request: %s", err2	)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("POST /measurements  returned status %d, expected %d", resp2.StatusCode, http.StatusOK)
	}
}

func TestGetMeasurements(t *testing.T) {
	client := &http.Client{}
	token := authenticate()
	req2, err2 := http.NewRequest("GET", base + "/abcdef/measurements", nil)
	req2.Header.Set("Bearer", string(token))
	resp2, err2 := client.Do(req2)
	if err2 != nil {
		t.Errorf("Error posting request: %s", err2	)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("GET /measurements returned status %d, expected %d", resp2.StatusCode, http.StatusOK)
	}
}