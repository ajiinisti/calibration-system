package request

import "time"

type CreateUser struct {
	Email            string
	Name             string
	Nik              string
	DateOfBirth      time.Time
	SupervisorNik    string
	BusinessUnitId   string
	OrganizationUnit string
	Division         string
	Department       string
	JoinDate         time.Time
	Grade            string
	Roles            []string
	HRBP             string
	Position         string
	ScoringMethod    string
	PhoneNumber      string
}

type UpdateUser struct {
	ID               string
	Email            string
	Name             string
	Nik              string
	DateOfBirth      time.Time
	SupervisorNik    string
	BusinessUnitId   string
	OrganizationUnit string
	Division         string
	Department       string
	JoinDate         time.Time
	Grade            string
	HRBP             string
	Position         string
	Roles            []string
	ScoringMethod    string
	PhoneNumber      string
}
