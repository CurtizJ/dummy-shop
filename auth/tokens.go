package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func NewToken(email string) (*Token, error) {
	var token Token
	var err error

	token.Access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Second * 30).Unix(),
	}).SignedString([]byte("kek"))

	if err != nil {
		return nil, err
	}

	bytes := make([]byte, 16)
	rand.Read(bytes)
	token.Refresh = fmt.Sprintf("%x", bytes)

	return &token, nil
}
