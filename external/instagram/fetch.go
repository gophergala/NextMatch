package instagram

import (
	"encoding/json"
	"net/http"
	"os"
)

var (
	at  = `2275718.1fb234f.580dabd7008f4529995d071001a5abe5` // access_token
	api = `https://api.instagram.com/v1/`                    //api_url
)

func init() {
	// Allow override of access_token and api urk
	if envAt := os.Getenv(`INSTAGRAM_ACCESS_TOKEN`); len(at) > 0 {
		at = envAt
	}

	if envAPI := os.Getenv(`INSTAGRAM_API`); len(envAPI) > 0 {
		api = envAPI
	}
}

func ByTag(tag string) (data Obj, err error) {
	var resp *http.Response

	if resp, err = http.Get(api + `/tags/` + tag + `/media/recent?access_token=` + at); err != nil {
		return
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&data); err != nil {
		return
	}

	return
}

func BuildTag(away, home string) string {
	return away + home
}
