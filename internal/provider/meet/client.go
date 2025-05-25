package meet

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type Credentials struct {
	URL          string
	OauthConfig  *oauth2.Config
	OauthToken   *oauth2.Token
	TokenStorage TokenStorage
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

	// Try to load token from storage if not provided
	if client.credentials.OauthToken == nil && client.credentials.TokenStorage != nil {
		client.credentials.OauthToken, err = client.credentials.TokenStorage.LoadToken()
		if err != nil {
			return nil, fmt.Errorf("failed to load token: %w", err)
		}
	}

	return
}

func (c *Client) LoginURL() string {
	return c.credentials.OauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (c *Client) CreateMeetEvent(ctx context.Context, req Request) (res Response, err error) {
	// Check if we have a token
	if c.credentials.OauthToken == nil {
		return res, fmt.Errorf("no OAuth token available")
	}

	// Create OAuth2 client with automatic token refresh
	tokenSource := c.credentials.OauthConfig.TokenSource(ctx, c.credentials.OauthToken)
	client := oauth2.NewClient(ctx, tokenSource)

	// Create Calendar service
	srv, err := calendar.New(client)
	if err != nil {
		return res, fmt.Errorf("calendar client error: %v", err)
	}

	event := &calendar.Event{
		Summary:     "Видеозвонок",
		Description: "Google Meet встреча",
		Start: &calendar.EventDateTime{
			DateTime: req.StartTime,
			TimeZone: "Asia/Almaty",
		},
		End: &calendar.EventDateTime{
			DateTime: req.EndTime,
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

	// After successful API call, save the current token
	if c.credentials.TokenStorage != nil {
		if token, err := tokenSource.Token(); err == nil {
			if err := c.credentials.TokenStorage.SaveToken(token); err != nil {
				// Log the error but don't fail the operation
				fmt.Printf("Failed to save token: %v\n", err)
			}
		}
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

	// Save the token if storage is available
	if c.credentials.TokenStorage != nil {
		if err := c.credentials.TokenStorage.SaveToken(token); err != nil {
			return nil, fmt.Errorf("failed to save token: %w", err)
		}
	}

	return token, nil
}
