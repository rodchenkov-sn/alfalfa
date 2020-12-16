package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)
var base = "https://alfalfa-project.herokuapp.com/api"
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
		t.Errorf("register returned status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestPostMeasurments(t *testing.T) {
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
		t.Errorf("register returned status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	token, _ := ioutil.ReadAll(resp.Body)
	var jsonStr2 = []byte(`{   "temperature": 36.6,   "timestamp": "2012-04-23T18:25:43.511Z" }`)
	req2, err2 := http.NewRequest("POST", base + "/measurements", bytes.NewBuffer(jsonStr2))
	req2.Header.Set("Bearer", string(token))
	resp2, err2 := client.Do(req)
	if err2 != nil {
		t.Errorf("Error posting request: %s", err	)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("POST /measutments  returned status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestGetMeasurments(t *testing.T) {
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
		t.Errorf("register returned status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	token, _ := ioutil.ReadAll(resp.Body)
	req2, err2 := http.NewRequest("GET", base + "/measurements", nil)
	req2.Header.Set("Bearer", string(token))
	resp2, err2 := client.Do(req)
	if err2 != nil {
		t.Errorf("Error posting request: %s", err	)
	}
	defer resp2.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /measutments returned status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}