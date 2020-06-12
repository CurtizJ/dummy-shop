package main

import (
	"fmt"
)

type SenderLog struct {
}

func (sender *SenderLog) Send(notification *Notification) error {
	message := fmt.Sprintf(`Will send message "%s" to "%s"\n`, notification.Message, notification.Email)
	fmt.Println(message)
	return nil
}
