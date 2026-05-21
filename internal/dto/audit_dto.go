package dto

import "time"

type AuditLogEvent struct {
	EventType  string                 `json:"event_type"`
	StaffID    *int64                 `json:"staff_id,omitempty"`
	HospitalID *int64                 `json:"hospital_id,omitempty"`
	Username   string                 `json:"username,omitempty"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}
