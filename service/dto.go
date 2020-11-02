package service

import "time"

type AuthInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type MeasurementWithAuth struct {
	Temperature float32   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
	User        *AuthInfo `json:"user"`
}

type Measurement struct {
	Temperature float32   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}
