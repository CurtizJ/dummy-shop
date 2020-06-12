package main

type Sender interface {
	Send(notification *Notification) error
}
