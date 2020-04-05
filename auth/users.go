package main

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users interface {
	Add(user User) error
	Get(email string) (*User, bool)
}
