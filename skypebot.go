package skypebot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/parnurzeal/gorequest"
)

const (
	_Debug = true

	_tokenURL    = "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token"
	_contentType = "application/x-www-form-urlencoded"
	_grantType   = "client_credentials"
	_scope       = "https://api.botframework.com/.default"

	_Info  = iota //0
	_Error        //1
)

var (
	//Loggers
	logInfo   *log.Logger
	logError  *log.Logger
	logOutput = os.Stderr
)

//Bot main bot object
type Bot struct {
	httpClient   *gorequest.SuperAgent
	clientID     string
	clientSecret string
	bearerToken  AuthToken
	messageCache []ResponseMessage
}

//Called after new, before main
func init() {
	fmt.Println("Init.")
	logInfo = log.New(logOutput, "INFO: ", log.Ldate|log.Ltime)
	logError = log.New(logOutput, "ERROR: ", log.Ldate|log.Ltime)
}

//New creates a new instance of your bot
func New(ID, secret string) *Bot {
	theBot := &Bot{
		httpClient:   gorequest.New(),
		clientID:     ID,
		clientSecret: secret,
	}

	//Set debugging
	if _Debug {
		theBot.httpClient.SetDebug(true)
	}

	//get authentication token
	theBot.GetToken()

	return theBot
}

//GetToken retrieve authentication token
func (obj *Bot) GetToken() *Bot {
	resp, _, errs := obj.httpClient.Post(_tokenURL).
		Type("form").
		Send(`{ "grant_type": "` + _grantType + `", "client_id": "` + obj.clientID +
			`", "client_secret":"` + obj.clientSecret + `", "scope":"` + _scope + `" }`).
		EndStruct(&obj.bearerToken)
	catchHTTPError(resp, errs)

	return obj
}

//MakeMessage sends a message to skype as a response
func (obj *Bot) MakeMessage(message string, pause int) *Bot {
	//Insert new message to cache

	return obj
}

//ShowTyping sends a typing gesture
func (obj *Bot) ShowTyping() *Bot {
	//Insert new message to cache

	return obj
}

//Send sends all messages in the cache/list
func (obj *Bot) Send() *Bot {
	logger(_Info, "Sending: %d messages", len(obj.messageCache))

	//Attemp to send

	//If fail get new token

	return obj
}

func catchHTTPError(resp gorequest.Response, errs []error) {
	for _, err := range errs {
		logger(_Error, err)
	}

	if resp.StatusCode != http.StatusOK {
		logger(_Error, errors.New("http error status "+resp.Status))
	}
}

func logger(typ int, str interface{}, args ...interface{}) {
	if _Debug {
		switch typ {
		case _Info:
			msg := str.(string)
			logInfo.Printf(msg+"\n", args...)
		case _Error:
			logError.Fatal(str)
		}
	}
}
