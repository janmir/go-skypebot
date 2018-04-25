package skypebot

import (
	"testing"
)

const (
	clientID     = "*****"
	clientSecret = "*****"
)

var (
	bot = New(clientID, clientSecret)
)

func TestMain(t *testing.T) {
	bot.MakeMessage("", 0).
		Send()
}