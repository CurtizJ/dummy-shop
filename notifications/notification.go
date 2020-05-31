package main

type Notification struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Subject string `json:"subject"`
}
