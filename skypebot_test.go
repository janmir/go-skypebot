package skypebot

import (
	"testing"
)

const (
	clientID     = "******"
	clientSecret = "******"
)

var (
	bot = New(clientID, clientSecret)
)

func TestMain(t *testing.T) {
	bot.Set("").
		ShowTyping(0).
		MakeMessage("", 0).
		Send()
}
