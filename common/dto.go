package common

import "time"

type AuthInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Measurement struct {
	Temperature float32   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}
