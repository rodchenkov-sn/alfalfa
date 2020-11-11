package common

import "time"

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Measurement struct {
	Temperature float32   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

type UserInfo struct {
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	Supervisors []string `json:"supervisors"`
}

type Supervisor struct {
	Login string `json:"login"`
}
