package main

type Config struct {
	SMTP_SERVER   string `env:"SMTP_SERVER" default:""`
	SMTP_LOGIN    string `env:"SMTP_LOGIN" default:""`
	SMTP_PASSWORD string `env:"SMTP_PASSWORD" default:""`
	SenderType    string `env:"SENDER" default:"Log"`
}
