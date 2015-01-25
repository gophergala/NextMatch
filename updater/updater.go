package updater

import (
	"encoding/json"
	"fmt"
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
	shortf  = `20060102`
	evtURI  = `https://erikberg.com/events.json?sport=%s&date=%s` // league[nba,nfl] date
	rsltURI = `https://erikberg.com//sport/results/%s.json`       // team_id
	token   = `0c301171-c857-4d46-b73f-d8f46602cae9`
)

var (
	decoder *json.Decoder
)

// BySport returns events in certain sport [with an optional date]
func BySport(sport string, date ...string) (ev Events, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)

	if len(date) < 1 {
		date = append(date, time.Now().Format(shortf))
	}

	uri := fmt.Sprintf(evtURI, sport, date[0])
	if req, err = http.NewRequest(`GET`, uri, nil); err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	if err = decode(resp, &ev); err != nil {
		return
	}

	return
}

// Result
func Result(team_id int) (rslt Results, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	if req, err = http.NewRequest(`GET`, fmt.Sprintf(rsltURI, team_id), nil); err != nil {
		return
	}

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	if err = decode(resp, rslt); err != nil {
		return
	}

	return
}

func decode(resp *http.Response, d interface{}) error {
	decoder = json.NewDecoder(resp.Body)

	return decoder.Decode(d)
}

func Unmarshal(b string, d interface{}) error {
	data := []byte(b)
	return json.Unmarshal(data, d)
}
