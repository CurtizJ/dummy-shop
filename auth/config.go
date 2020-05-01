package main

type Config struct {
	AccessExpiration  uint   `env:"ACCESS_TOKEN_EXPIRATION" default:"86400"`
	RefreshExpiration uint   `env:"REFRESH_TOKEN_EXPIRATION" default:"604800"`
	AuthSecret        string `env:"AUTH_SECRET"`
}
