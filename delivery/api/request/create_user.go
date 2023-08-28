package request

import "time"

type CreateUser struct {
	Email            string
	Name             string
	Nik              string
	DateOfBirth      time.Time
	SupervisorName   string
	BusinessUnitId   string
	OrganizationUnit string
	Division         string
	Department       string
	HireDate         time.Time
	Grade            string
	Role             []string
}

type UpdateUser struct {
	ID               string
	Email            string
	Name             string
	Nik              string
	DateOfBirth      time.Time
	SupervisorName   string
	BusinessUnitId   string
	OrganizationUnit string
	Division         string
	Department       string
	HireDate         time.Time
	Grade            string
	Role             []string
}
