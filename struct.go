package skypebot

//AuthToken token received from bot framework
type AuthToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
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
	Sleep int8
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
