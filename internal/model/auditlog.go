package model

import "time"

type AuditLog struct {
	ClientID  string                 `json:"client_id"`
	Timestamp time.Time              `json:"timestamp"`
	Entity    string                 `json:"entity"`
	Mutation  string                 `json:"mutation"`
	Before    map[string]interface{} `json:"before"`
	After     map[string]interface{} `json:"after"`
}
