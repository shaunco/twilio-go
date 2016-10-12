package twilio

import (
	"net/url"

	types "github.com/kevinburke/go-types"
)

const pathPart = "Messages"

type MessageService struct {
	client *Client
}

type Direction string

const DirectionOutboundReply = Direction("outbound-reply")

type Status string

const StatusSent = Status("sent")
const StatusReceived = Status("received")
const StatusDelivered = Status("delivered")

type Message struct {
	Sid                 string           `json:"sid"`
	Body                string           `json:"body"`
	From                string           `json:"from"`
	To                  string           `json:"to"`
	Price               string           `json:"price"`
	Status              Status           `json:"status"`
	AccountSid          string           `json:"account_sid"`
	MessagingServiceSid types.NullString `json:"messaging_service_sid"`
	DateCreated         TwilioTime       `json:"date_created"`
	DateUpdated         TwilioTime       `json:"date_updated"`
	DateSent            TwilioTime       `json:"date_sent"`
	NumSegments         Segments         `json:"num_segments"`
	// TODO fix type here... UintStr or something ?
	NumMedia        Segments          `json:"num_media"`
	PriceUnit       string            `json:"price_unit"`
	Direction       Direction         `json:"direction"`
	SubresourceURIs map[string]string `json:"subresource_uris"`
	URI             string            `json:"uri"`
	APIVersion      string            `json:"api_version"`
}

type MessagePage struct {
	Page
	Messages []*Message `json:"messages"`
}

type MessageIterator struct {
	pos         int
	messages    []*Message
	nextPageURI string
	data        url.Values
	service     *MessageService
}

// Create a message with the given values.
func (m *MessageService) Create(data url.Values) (*Message, error) {
	msg := new(Message)
	err := m.client.CreateResource(pathPart, data, msg)
	return msg, err
}

// SendMessage is a convenience wrapper around Create.
func (m *MessageService) SendMessage(from string, to string, body string, mediaURLs []url.URL) (*Message, error) {
	v := url.Values{
		"Body": []string{body},
		"From": []string{from},
		"To":   []string{to},
	}
	if mediaURLs != nil {
		for _, mediaURL := range mediaURLs {
			v.Add("MediaUrl", mediaURL.String())
		}
	}
	return m.Create(v)
}

type MessagePageIterator struct {
	client      *Client
	nextPageURI types.NullString
	data        url.Values
	count       uint
}

// GetPage returns a single page of resources. To retrieve multiple pages, use
// GetPageIterator.
func (m *MessageService) GetPage(data url.Values) (*MessagePage, error) {
	mp := new(MessagePage)
	err := m.client.ListResource(pathPart, data, mp)
	return mp, err
}

// GetPageIterator returns an iterator which can be used to retrieve pages.
func (m *MessageService) GetPageIterator(data url.Values) *MessagePageIterator {
	return &MessagePageIterator{
		client:      m.client,
		nextPageURI: types.NullString{},
		data:        data,
		count:       0,
	}
}

func (m *MessagePageIterator) Next() (*MessagePage, error) {
	mp := new(MessagePage)
	var err error
	if m.count == 0 {
		err = m.client.ListResource(pathPart, m.data, mp)
	} else {
		err = m.client.GetNextPage(m.nextPageURI.String, mp)
	}
	if err != nil {
		return nil, err
	}
	if mp.NextPageURI.Valid == false {
		return nil, NoMoreResults
	}
	m.count++
	m.nextPageURI = mp.NextPageURI
	return mp, nil
}