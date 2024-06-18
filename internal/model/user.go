package model

type User struct {
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Accounts []Account `json:"accounts,omitempty"`
}
