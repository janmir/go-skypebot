package skypebot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/parnurzeal/gorequest"
)

const (
	_Debug = true

	_tokenURL    = "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token"
	_contentType = "application/x-www-form-urlencoded"
	_grantType   = "client_credentials"
	_scope       = "https://api.botframework.com/.default"
	_templateURL = "%s/v3/conversations/%s/activities/%s"

	_Info  = iota //0
	_Error        //1

	_keyValueDB = "cache.db"
	_dataBucket = "data"
)

var (
	f = fmt.Sprintf

	//Loggers
	logInfo   *log.Logger
	logError  *log.Logger
	logOutput = os.Stderr

	//Key values store database
	database *bolt.DB

	/**Persistent data**/

	//bearerToken Global auth token
	bearerToken AuthToken
	//service url
	serviceURL string
)

//BotManager manages bot instances
type BotManager struct {
	Bots map[string]Bot
}

//Bot main bot object
type Bot struct {
	httpClient   *gorequest.SuperAgent
	clientID     string
	clientSecret string
	messageCache []ResponseMessage
	request      RequestMessage
	replyURL     string
}

/*
 * ┌────────────────────────────────────────┐
 * │        BotManager Main Methods         │
 * └────────────────────────────────────────┘
 */

//Get the bot base on id/key
func (man *BotManager) Get(key string) *Bot {
	b := man.Bots[key]
	return &b
}

/*
 * ┌─────────────────────────────────────┐
 * │        Initializations/New          │
 * └─────────────────────────────────────┘
 */
//Called after new, before main
func init() {
	fmt.Println("Init.")
	logInfo = log.New(logOutput, "INFO: ", log.Ldate|log.Ltime)
	logError = log.New(logOutput, "ERROR: ", log.Ldate|log.Ltime)

	//Init using previous serviceURL
	database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_dataBucket))

		if value := bucket.Get([]byte("serviceURL")); value != nil {
			serviceURL = string(value)
		}

		return nil
	})
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
	if bearerToken.AccessToken == "" {
		initDatabase()

		//get first from key-data-store
		database.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(_dataBucket))

			if value := bucket.Get([]byte("bearerToken")); value != nil {
				bearerToken.AccessToken = string(value)
			} else {
				//request
				theBot.GetToken()

				//save
				return bucket.Put([]byte("bearerToken"),
					[]byte(bearerToken.AccessToken))
			}
			return nil
		})
	}

	return theBot
}

/*
 * ┌─────────────────────────────────┐
 * │        Bot Main Methods         │
 * └─────────────────────────────────┘
 */

//GetToken retrieve authentication token
func (obj *Bot) GetToken() *Bot {
	resp, _, errs := obj.httpClient.Post(_tokenURL).
		Type("form").
		Send(`{ "grant_type": "` + _grantType + `", "client_id": "` + obj.clientID +
			`", "client_secret":"` + obj.clientSecret + `", "scope":"` + _scope + `" }`).
		EndStruct(&bearerToken)
	catchHTTPError(resp, errs, nil)

	return obj
}

//SetDefaultServiceURL  sets the default service url
func (obj *Bot) SetDefaultServiceURL(url string) *Bot {
	//Save in key-value store
	serviceURL = url

	//Save to cache file
	database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_dataBucket))
		return bucket.Put([]byte("serviceURL"), []byte(serviceURL))
	})

	return obj
}

//Set sets the request from client data
func (obj *Bot) Set(req interface{}) *Bot {
	switch req.(type) {
	//For responses
	case RequestMessage:
		obj.request = req.(RequestMessage)
		obj.replyURL = f(_templateURL, obj.request.ServiceURL, url.QueryEscape(obj.request.Conversation.ID),
			url.QueryEscape(obj.request.ID))
	//For pro-active messages
	case string:
		conversationID := req.(string)
		obj.request = RequestMessage{
			ServiceURL: "", //required access from db
			From: _From{
				ID:   "", //bot
				Name: "", //bot id
			},
			Recipient: _Recipient{
				ID:   "", //can be blank
				Name: "", //can be blank, i think
			},
			Conversation: _Conversation{
				ID:   conversationID,
				Name: "",
			},
		}

		obj.replyURL = f(_templateURL, obj.request.ServiceURL, url.QueryEscape(conversationID), "")
	}
	return obj
}

//MakeMessage sends a message to skype as a response
func (obj *Bot) MakeMessage(message string, pause int) *Bot {
	request := obj.request

	//Insert new message to cache
	response := ResponseMessage{
		Type: "message",
		From: _From{
			ID:   request.Recipient.ID,
			Name: request.Recipient.Name,
		},
		Conversation: _Conversation{
			ID:   request.Conversation.ID,
			Name: request.Conversation.Name,
		},
		Recipient: _Recipient{
			ID:   request.From.ID,
			Name: request.From.Name,
		},
		Locale:     "en-US",
		Text:       message,
		TextFormat: "markdown", //plain", //"xml",
		ReplyToID:  request.ID,
		InputHint:  "ignoringInput",
	}

	//Append
	obj.messageCache = append(obj.messageCache, response)

	return obj
}

//ShowTyping sends a typing gesture
func (obj *Bot) ShowTyping(pause int) *Bot {
	request := obj.request

	//Insert new message to cache
	response := ResponseMessage{
		Text: "...",
		Type: "typing",
		From: _From{
			ID:   request.Recipient.ID,
			Name: request.Recipient.Name,
		},
		Conversation: _Conversation{
			ID:   request.Conversation.ID,
			Name: request.Conversation.Name,
		},
		Recipient: _Recipient{
			ID:   request.From.ID,
			Name: request.From.Name,
		},
		ReplyToID: request.ID,
		Sleep:     pause,
	}

	//Append
	obj.messageCache = append(obj.messageCache, response)

	return obj
}

//Send sends all messages in the cache/list
func (obj *Bot) Send() *Bot {
	logger(_Info, "Sending: %d messages", len(obj.messageCache))

	retries := 2
	counter := 0

	//Loop thru all
	for _, message := range obj.messageCache {

		//Attemp to send
		err := retry(&counter, func() bool {
			resp, _, errs := obj.httpClient.Post(obj.replyURL).
				Set("Authorization", bearerToken.AccessToken).
				Send(message).
				End()

			return catchHTTPError(resp, errs, func(status int) {
				if resp.StatusCode == http.StatusUnauthorized {
					logger(_Info, "Auth Token Expired...")
					logger(_Info, "Requesting new...")

					//If fail get new token
					database.Update(func(tx *bolt.Tx) error {
						bucket := tx.Bucket([]byte(_dataBucket))

						//request
						obj.GetToken()

						//save
						return bucket.Put([]byte("bearerToken"),
							[]byte(bearerToken.AccessToken))
					})
				}
			})
		}, retries)
		catch(err)

		//Sleep gamay
		time.Sleep(time.Millisecond * time.Duration(message.Sleep))
	}

	return obj
}

/*
 * ┌─────────────────────────────┐
 * │      Utility Functions      │
 * └─────────────────────────────┘
 */
func retry(counter *int, fn func() bool, try int) error {
	for *counter < try {
		if fn() {
			return nil
		}
		*counter++
	}

	return errors.New("retry failed")
}

func catchHTTPError(resp gorequest.Response, errs []error, callback func(int)) bool {
	for _, err := range errs {
		catch(err)
	}

	if resp.StatusCode != http.StatusOK {
		//maybe add log here, response body

		//check if auth err
		if callback != nil {
			callback(resp.StatusCode)
			return false
		} else {
			logger(_Error, errors.New("http error status "+resp.Status))
		}
	}

	return true
}

func catch(err error) {
	if err != nil {
		logger(_Error, err)
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

func initDatabase() {
	//Key value store
	db, err := bolt.Open(_keyValueDB, 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	catch(err)
	database = db

	database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(_dataBucket))
		return err
	})
}
