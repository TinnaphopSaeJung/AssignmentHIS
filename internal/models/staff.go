package models

type Staff struct {
	ID           int64
	Username     string
	PasswordHash string
	HospitalID   int64
	FirstNameTH  string
	MiddleNameTH string
	LastNameTH   string
	FirstNameEN  string
	MiddleNameEN string
	LastNameEN   string
}
