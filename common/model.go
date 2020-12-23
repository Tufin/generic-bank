package common

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type RequestAccount struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type Balance struct {
	Label  string `json:"label" form:"label" binding:"required"`
	Amount int    `json:"amount" form:"amount" binding:"required"`
}

type SSNAccount struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	SSN      string `json:"ssn"`
}
