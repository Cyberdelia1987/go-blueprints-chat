package main

import (
	"time"
)

type message struct {
	UserId    string
	Name      string
	Message   string
	AvatarURL string
	When      time.Time
}
