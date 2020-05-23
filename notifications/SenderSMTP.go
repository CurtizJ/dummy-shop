package main

type SenderSMTP struct {
	server   string
	login    string
	password string
}

func (sender *SenderSMTP) Send(notification *Notification) error {
	return nil
}
