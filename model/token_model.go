package model

type TokenModel struct {
	Username string
	Email    string
	Role     []string
	ID       string
	Name     string
}

type ModifiedTokenModel struct {
	Username     string
	Email        string
	Role         []string
	ID           string
	Name         string
	Nik          string
	Division     string
	BusinessUnit BusinessUnit
}
