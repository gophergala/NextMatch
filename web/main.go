package main

import (
	"flag"
	"github.com/gophergala/NextMatch/updater/xmlstats"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var port = flag.String("p", "80", "the port on wich we're serving")

func init() {
	addTfunc("parse", time.Parse)
}

func handler(w http.ResponseWriter, r *http.Request) {

	e, err := xmlstats.BySport("nba", "20150123")
	if err != nil {
		log.Printf("Didn't get data :( err = %v", err)
	}

	for i, _ := range e.Event {
		e.Event[i].AwayTeam.Logo = teamLogos[e.Event[i].AwayTeam.TeamID]
		e.Event[i].HomeTeam.Logo = teamLogos[e.Event[i].HomeTeam.TeamID]
		//ev.AwayTeam.Logo = teamLogos[ev.AwayTeam.TeamID]
		//ev.HomeTeam.Logo = teamLogos[ev.HomeTeam.TeamID]
	}

	renderArgs := args{
		"events": e,
		"title":  "Home",
	}

	execT(w, "home", renderArgs)
}

func reload(rw http.ResponseWriter, req *http.Request) {
	loadTmpl()
	http.Redirect(rw, req, "/", 302)
}

func main() {

	xmlstats.Token = os.Getenv("XMLSTATS_TOKEN")
	if len(xmlstats.Token) == 0 {
		log.Fatal("Specify XMLSTATS_TOKEN environment variable")
	}

	flag.Parse()
	loadTmpl()
	http.HandleFunc("/refresh", reload)
	http.Handle("/static/", static(http.FileServer(http.Dir("."))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func static(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/static", "/resources", 1)
		h.ServeHTTP(w, r)
	}
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
