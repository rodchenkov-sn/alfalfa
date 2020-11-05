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

type Organization struct {
	Credentials Credentials        `json:"credentials"`
	Policy      OrganizationPolicy `json:"policy"`
}

type OrganizationPolicy struct {
	FreeJoin           bool `json:"free_join"`
	FreeLeave          bool `json:"free_leave"`
	CanGetMeasurements bool `json:"can_get_measurements"`
	CanAddMeasurements bool `json:"can_add_measurements"`
}
