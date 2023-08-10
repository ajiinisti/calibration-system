package model

import "time"

type Employee struct {
	BaseModel
	EmployeeId       string `gorm:"unique"`
	Email            string `gorm:"unique" `
	Name             string
	BirthDate        time.Time
	OrganizationUnit string
	BusinessUnit     string
	Department       string
	Grade            string
}
