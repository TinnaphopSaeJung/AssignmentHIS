package models

type Staff struct {
	ID           int64
	Username     string
	PasswordHash string
	HospitalID   int64
}
