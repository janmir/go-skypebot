package skypebot

//AuthToken token received from bot framework
type AuthToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

//RequestMessage is for request
type RequestMessage struct {
	Type           string        `json:"type,omitempty"`
	Action         string        `json:"action,omitempty"`
	ID             string        `json:"id,omitempty"`
	ChannelID      string        `json:"channelId,omitempty"`
	Text           string        `json:"text,omitempty"`
	Locale         string        `json:"locale,omitempty"`
	TextFormat     string        `json:"textFormat,omitempty"`
	Timestamp      string        `json:"timestamp,omitempty"`
	LocalTimestamp string        `json:"localTimestamp,omitempty"`
	From           _From         `json:"from,omitempty"`
	Recipient      _Recipient    `json:"recipient,omitempty"`
	Conversation   _Conversation `json:"conversation,omitempty"`
	MembersAdded   []struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"membersAdded,omitempty"`
	MembersRemoved []struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"membersRemoved,omitempty"`
	ChannelData struct {
		ClientActivityID string `json:"clientActivityId,omitempty"`
		Text             string `json:"text,omitempty"`
	} `json:"channelData,omitempty"`
	Entities []struct {
		Type              string `json:"type,omitempty"`
		RequiresBotState  bool   `json:"requiresBotState,omitempty"`
		SupportsTts       bool   `json:"supportsTts,omitempty"`
		SupportsListening bool   `json:"supportsListening,omitempty"`
		Text              string `json:"text,omitempty,omitempty"`
		Locale            string `json:"locale,omitempty,omitempty"`
		Country           string `json:"country,omitempty,omitempty"`
		Platform          string `json:"platform,omitempty,omitempty"`
		Mentioned         struct {
			ID string `json:"id,omitempty"`
		} `json:"mentioned,omitempty"`
	} `json:"entities,omitempty"`
	ServiceURL string `json:"serviceUrl,omitempty"`
}

//ResponseMessage is for response
type ResponseMessage struct {
	Type         string         `json:"type"`
	From         _From          `json:"from"`
	Conversation _Conversation  `json:"conversation"`
	Recipient    _Recipient     `json:"recipient"`
	Text         string         `json:"text"`
	ReplyToID    string         `json:"replyToId"`
	TextFormat   string         `json:"textFormat"`
	Locale       string         `json:"locale"`
	InputHint    string         `json:"inputHint"`
	Attachments  []_Attachments `json:"attachments"`
	//custom
	Sleep int
}

type _From struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type _Conversation struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsGroup bool   `json:"isGroup"`
}

type _Recipient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type _Attachments struct {
	ContentType string   `json:"contentType"`
	ContentURL  string   `json:"contentUrl"`
	Content     _Content `json:"content"`
}

type _Content struct {
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Text     string    `json:"text"`
	Type     string    `json:"type"`
	Version  string    `json:"version"`
	To       []string  `json:"to"` // In singIn: list of users that this will be visible
	Items    []_Items  `json:"items"`
	Tax      string    `json:"tax"`
	Total    string    `json:"total"`
	Images   []_Images `json:"images"`
	Tap      _Tap      `json:"tap"`
	Facts    []struct {
		Value string `json:"value"`
		Key   string `json:"key"`
	} `json:"facts"`
	Body []struct {
		Type       string `json:"type"`
		Text       string `json:"text"`
		Size       string `json:"size,omitempty"`
		Separation string `json:"separation,omitempty"`
	} `json:"body"`
	Actions []struct {
		Type  string `json:"type"` //Type of action to perform.
		URL   string `json:"url"`
		Title string `json:"title"` //In signIn card: Text of the button. Only applicable for a button's action.

		Image string `json:"image"` //Image to display
		Text  string `json:"text"`  //Text for the action
		Value string `json:"value"`
	} `json:"actions"`
	Buttons []struct {
		Type  string `json:"type"`
		Title string `json:"title"`
		Value string `json:"value"`
		Image string `json:"image,omitempty"`
	} `json:"buttons"`
}

type _Items struct {
	Price string `json:"price"`
	Title string `json:"title"`
}

type _Tap struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Title string `json:"title"`
}
type _Images struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}
