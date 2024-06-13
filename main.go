package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	var authCode string
	waitCh := make(chan struct{})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode = r.URL.Query().Get("code")
		close(waitCh)
		fmt.Fprint(w, "Authentication successful! You can close this window.")
	})

	go func() {
		http.ListenAndServe("localhost:3000", nil)
	}()

	go openBrowser(authURL)

	<-waitCh

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func createEventsFromCSV(csvFile string) ([]*calendar.Event, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open CSV file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read CSV file: %w", err)
	}

	var events []*calendar.Event
	for _, record := range records[1:] {
		summary := record[0]
		startDateTime := record[1]
		endDateTime := record[2]
		timeZone := record[3]
		attendees := strings.Split(record[4], ";")

		var eventAttendees []*calendar.EventAttendee
		for _, attendee := range attendees {
			eventAttendees = append(eventAttendees, &calendar.EventAttendee{Email: attendee})
		}

		event := &calendar.Event{
			Summary: summary,
			Start: &calendar.EventDateTime{
				DateTime: startDateTime,
				TimeZone: timeZone,
			},
			End: &calendar.EventDateTime{
				DateTime: endDateTime,
				TimeZone: timeZone,
			},
			Attendees: eventAttendees,
		}

		events = append(events, event)
	}

	return events, nil
}

func insertEventsToCalendar(srv *calendar.Service, events []*calendar.Event) {
	for _, event := range events {
		calendarId := "primary"
		createdEvent, err := srv.Events.Insert(calendarId, event).Do()
		if err != nil {
			log.Printf("Unable to create event for %s: %v\n", event.Summary, err)
		} else {
			fmt.Printf("Event created: %s\n", createdEvent.HtmlLink)
		}
	}
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	events, err := createEventsFromCSV("events.csv")
	if err != nil {
		log.Fatalf("Unable to create events from CSV: %v", err)
	}

	insertEventsToCalendar(srv, events)
}
