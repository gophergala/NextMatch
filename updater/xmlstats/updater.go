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
		Logo         string
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

	for i, _ := range ev.Event {
		ev.Event[i].AwayTeam.Logo = teamLogos[ev.Event[i].AwayTeam.TeamID]
		ev.Event[i].HomeTeam.Logo = teamLogos[ev.Event[i].HomeTeam.TeamID]
		//ev.AwayTeam.Logo = teamLogos[ev.AwayTeam.TeamID]
		//ev.HomeTeam.Logo = teamLogos[ev.HomeTeam.TeamID]
	}

	return ev, err
}

// Result
func Result(teamId string) (results Results, err error) {
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

var teamLogos = map[string]string{
	"atlanta-hawks":          "http://content.sportslogos.net/logos/6/220/thumbs/5mdhgjh3aa92kih09pgi.gif",
	"boston-celtics":         "http://content.sportslogos.net/logos/6/213/thumbs/slhg02hbef3j1ov4lsnwyol5o.gif",
	"brooklyn-nets":          "http://content.sportslogos.net/logos/6/3786/thumbs/hsuff5m3dgiv20kovde422r1f.gif",
	"charlotte-hornets":      "http://content.sportslogos.net/logos/6/5120/thumbs/512019262015.gif",
	"chicago-bulls":          "http://content.sportslogos.net/logos/6/221/thumbs/hj3gmh82w9hffmeh3fjm5h874.gif",
	"cleveland-cavaliers":    "http://content.sportslogos.net/logos/6/222/thumbs/e4701g88mmn7ehz2baynbs6e0.gif",
	"dallas-mavericks":       "http://content.sportslogos.net/logos/6/228/thumbs/ifk08eam05rwxr3yhol3whdcm.gif",
	"denver-nuggets":         "http://content.sportslogos.net/logos/6/229/thumbs/xeti0fjbyzmcffue57vz5o1gl.gif",
	"detroit-pistons":        "http://content.sportslogos.net/logos/6/223/thumbs/3079.gif",
	"golden-state-warriors":  "http://content.sportslogos.net/logos/6/235/thumbs/qhhir6fj8zp30f33s7sfb4yw0.gif",
	"houston-rockets":        "http://content.sportslogos.net/logos/6/230/thumbs/8xe4813lzybfhfl14axgzzqeq.gif",
	"indiana-pacers":         "http://content.sportslogos.net/logos/6/224/thumbs/3083.gif",
	"los-angeles-clippers":   "http://content.sportslogos.net/logos/6/236/thumbs/bvv028jd1hhr8ee8ii7a0fg4i.gif",
	"los-angeles-lakers":     "http://content.sportslogos.net/logos/6/237/thumbs/uig7aiht8jnpl1szbi57zzlsh.gif",
	"memphis-grizzlies":      "http://content.sportslogos.net/logos/6/231/thumbs/793.gif",
	"miami-heat":             "http://content.sportslogos.net/logos/6/214/thumbs/burm5gh2wvjti3xhei5h16k8e.gif",
	"milwaukee-bucks":        "http://content.sportslogos.net/logos/6/225/thumbs/0295onf2c4xsbfsxye6i.gif",
	"minnesota-timberwolves": "http://content.sportslogos.net/logos/6/232/thumbs/zq8qkfni1g087f4245egc32po.gif",
	"new-orleans-pelicans":   "http://content.sportslogos.net/logos/6/4962/thumbs/496226812014.gif",
	"new-york-knicks":        "http://content.sportslogos.net/logos/6/216/thumbs/2nn48xofg0hms8k326cqdmuis.gif",
	"oklahoma-city-thunder":  "http://content.sportslogos.net/logos/6/2687/thumbs/khmovcnezy06c3nm05ccn0oj2.gif",
	"orlando-magic":          "http://content.sportslogos.net/logos/6/217/thumbs/wd9ic7qafgfb0yxs7tem7n5g4.gif",
	"philadelphia-76ers":     "http://content.sportslogos.net/logos/6/218/thumbs/qlpk0etqwelv8artgc7tvqefu.gif",
	"phoenix-suns":           "http://content.sportslogos.net/logos/6/238/thumbs/23843702014.gif",
	"portland-trail-blazers": "http://content.sportslogos.net/logos/6/239/thumbs/bahmh46cyy6eod2jez4g21buk.gif",
	"sacramento-kings":       "http://content.sportslogos.net/logos/6/240/thumbs/832.gif",
	"san-antonio-spurs":      "http://content.sportslogos.net/logos/6/233/thumbs/827.gif",
	"toronto-raptors":        "http://content.sportslogos.net/logos/6/227/thumbs/yfypcwqog6qx8658sn5w65huh.gif",
	"utah-jazz":              "http://content.sportslogos.net/logos/6/234/thumbs/m2leygieeoy40t46n1qqv0550.gif",
	"washington-wizards":     "http://content.sportslogos.net/logos/6/219/thumbs/b3619brnphtx65s2th4p9eggf.gif",
}
