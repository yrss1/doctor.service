package meet

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type Credentials struct {
	URL         string
	OauthConfig *oauth2.Config
	OauthToken  *oauth2.Token
}

type Client struct {
	credentials Credentials
}

type Request struct {
	UserEmail   string `json:"user_email"`
	DoctorEmail string `json:"doctor_email"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type Response struct {
	MeetLink  string `json:"meet_link"`
	StartTime string `json:"start_time"`
}

func New(credentials Credentials) (client *Client, err error) {
	client = &Client{
		credentials: credentials,
	}
	return
}

func (c *Client) LoginURL() string {
	return c.credentials.OauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (c *Client) CreateMeetEvent(ctx context.Context, req Request) (res Response, err error) {
	client := c.credentials.OauthConfig.Client(ctx, c.credentials.OauthToken)

	srv, err := calendar.New(client)
	if err != nil {
		return res, fmt.Errorf("calendar client error: %v", err)
	}

	now := time.Now()
	event := &calendar.Event{
		Summary:     "Видеозвонок",
		Description: "Google Meet встреча",
		Start: &calendar.EventDateTime{
			DateTime: now.Add(24 * time.Hour).Format(time.RFC3339),
			TimeZone: "Asia/Almaty",
		},
		End: &calendar.EventDateTime{
			DateTime: now.Add(24*time.Hour + 40*time.Minute).Format(time.RFC3339),
			TimeZone: "Asia/Almaty",
		},
		Attendees: []*calendar.EventAttendee{
			{
				Email:       req.DoctorEmail,
				DisplayName: "Doctor",
				Optional:    false,
			},
			{
				Email:       req.UserEmail,
				DisplayName: "Patient",
				Optional:    false,
			}},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: fmt.Sprintf("go-meet-%d", time.Now().UnixNano()),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
	}

	createdEvent, err := srv.Events.Insert("primary", event).ConferenceDataVersion(1).Do()
	if err != nil {
		return res, fmt.Errorf("error creating calendar event: %v", err)
	}

	for _, entry := range createdEvent.ConferenceData.EntryPoints {
		if entry.EntryPointType == "video" {
			res.MeetLink = entry.Uri
			res.StartTime = createdEvent.Start.DateTime
			return res, nil
		}
	}

	return res, fmt.Errorf("no Google Meet link found")
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.credentials.OauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	c.credentials.OauthToken = token

	return token, nil
}
