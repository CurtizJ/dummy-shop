package main

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	RoleUnknown int32 = iota
	RoleUser
	RoleAdmin
)
