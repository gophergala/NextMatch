package xmlstats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	Events struct {
		Event []struct {
			EventID       string `json:"event_id"`
			EventStatus   string `json:"event_status"`
			AwayTeam      Team   `json:"away_team"`
			HomeTeam      Team   `json:"home_team"`
			SeasonType    string `json:"season_type"`
			Site          Site   `json:"site"`
			Sport         string `json:"sport"`
			StartDateTime string `json:"start_date_time"`
		} `json:"event"`
		EventsDate string `json:"events_date"`
	}

	Results struct {
		EventID                 string  `json:"event_id"`
		EventSeasonType         string  `json:"event_season_type"`
		EventStartDateTime      string  `json:"event_start_date_time"`
		EventStatus             string  `json:"event_status"`
		Opponent                Team    `json:"opponent"`
		OpponentEventsLost      float64 `json:"opponent_events_lost"`
		OpponentEventsWon       float64 `json:"opponent_events_won"`
		OpponentPointsScored    float64 `json:"opponent_points_scored"`
		Site                    Site    `json:"site"`
		Team                    Team    `json:"team"`
		TeamEventLocationType   string  `json:"team_event_location_type"`
		TeamEventNumberInSeason float64 `json:"team_event_number_in_season"`
		TeamEventResult         string  `json:"team_event_result"`
		TeamEventsLost          float64 `json:"team_events_lost"`
		TeamEventsWon           float64 `json:"team_events_won"`
		TeamPointsScored        float64 `json:"team_points_scored"`
	}

	Team struct {
		Abbreviation string `json:"abbreviation"`
		Active       bool   `json:"active"`
		City         string `json:"city"`
		Conference   string `json:"conference"`
		Division     string `json:"division"`
		FirstName    string `json:"first_name"`
		FullName     string `json:"full_name"`
		LastName     string `json:"last_name"`
		SiteName     string `json:"site_name"`
		State        string `json:"state"`
		TeamID       string `json:"team_id"`
		Logo        string
	}

	Site struct {
		Capacity float64 `json:"capacity"`
		City     string  `json:"city"`
		Name     string  `json:"name"`
		State    string  `json:"state"`
		Surface  string  `json:"surface"`
	}
)

const (
	shortf    = "20060102"
	eventURI  = "https://erikberg.com/events.json?sport=%s&date=%s" // league[nba,nfl] date
	resultURI = "https://erikberg.com/sport/results/%s.json"        // team_id
	userAgent = "nextmatch/0.1 (https://twitter.com/oscarryz)"
	auth      = "Bearer %s"
)

var (
	Token string
)

var cache = make(map[string]interface{})

func doRequest(uri string, result interface{}) error {

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf(auth, Token))
	req.Header.Add("User-agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	log.Printf("%s for URI %s", resp.Status, uri)

	err = decode(resp, result)
	return err
}

// BySport returns events in certain sport [with an optional date]
func BySport(sport string, date ...string) (ev Events, err error) {

	if len(date) < 1 {
		date = append(date, time.Now().Format(shortf))
	}
	uri := fmt.Sprintf(eventURI, sport, date[0])
	if cache[uri] != nil {
		log.Printf("cache  for URI %s", uri)

		return cache[uri].(Events), nil
	}
	err = doRequest(uri, &ev)
	cache[uri] = ev
	return ev, err
}

// Result
func Result(teamId int) (results Results, err error) {
	uri := fmt.Sprintf(resultURI, teamId)
	if cache[uri] != nil {
		return cache[uri].(Results), nil
	}
	err = doRequest(uri, &results)
	cache[uri] = results
	return results, err
}

func decode(resp *http.Response, d interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(d)
}

func Unmarshal(b string, d interface{}) error {
	data := []byte(b)
	return json.Unmarshal(data, d)
}
