package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bmizerany/pat"
)

var image_limit = 300
var last_refreshed_at time.Time
var image_mapping = make(map[string][]string)
var api_key string
var kinds = []string{"pug", "corgi", "shiba", "cat", "giraffe"}

type PhotoProperty struct {
	Url string
}

type Photo struct {
	OriginalPhoto PhotoProperty `json:"original_size"`
}

// TODO Timestamp sometime is string
type Blog struct {
	Timestamp int
	Photos    []Photo
}

type TaggedApiResponse struct {
	Blogs []Blog `json:"response"`
}

type FetcherInterface interface {
	Fetch(url string) ([]byte, error)
}

type tumblrFetcher struct{}

func GetJsonString(v interface{}) string {
	result, err := json.MarshalIndent(v, "", "    ")
	panic_on_error(err)
	return string(result)
}

func instruction(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, GetJsonString(endpoints()))
}

func count(res http.ResponseWriter, req *http.Request) {
	kind := req.URL.Query().Get(":kind")
	fmt.Fprint(res, GetJsonString(&map[string]int{"count": len(image_mapping[kind])}))
}

func random(res http.ResponseWriter, req *http.Request) {
	var err error
	refresh_every_hour()
	kind := req.URL.Query().Get(":kind")
	err = check_kind(kind)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, GetJsonString(&map[string]string{"error": err.Error()}))
		return
	}
	action := req.URL.Query().Get(":action")
	if len(image_mapping[kind]) == 0 {
		done := make(chan bool)
		go check_image_presence(kind, done)
		<-done
	}
	index := rand.Intn(len(image_mapping[kind]))
	url := image_mapping[kind][index]
	if action == "preview" {
		fmt.Fprint(res, "<html><img src='"+url+"' /></html>")
	} else if action == "url" {
		fmt.Fprint(res, url)
	} else {
		fmt.Fprint(res, GetJsonString(&map[string]string{"url": url}))
	}
}

func bomb(res http.ResponseWriter, req *http.Request) {
	refresh_every_hour()
	var result []string
	kind := req.URL.Query().Get(":kind")
	number := req.URL.Query().Get(":number")
	if number == "" {
		number = "4"
	}
	number_str, _ := strconv.Atoi(number)
	permutation := rand.Perm(len(image_mapping[kind]))
	for _, pos := range permutation[:number_str] {
		result = append(result, image_mapping[kind][pos])
	}
	fmt.Fprint(res, GetJsonString(&map[string][]string{"urls": result}))
}

func panic_on_error(err error) {
	if err != nil {
		panic(err)
	}
}

func log_on_error(err error) {
	if err != nil {
		log.Println("Error: ", err)
	}
}

func endpoints() *map[string]map[string]string {
	return &map[string]map[string]string{
		"DEMO": {
			"pug":     "http://awwimage.herokuapp.com/random/pug/preview",
			"corgi":   "http://awwimage.herokuapp.com/random/corgi/preview",
			"cat":     "http://awwimage.herokuapp.com/random/cat/preview",
			"giraffe": "http://awwimage.herokuapp.com/random/giraffe/preview",
		},
		"ENDPOINT": {
			"/instruction":             "Get a random image. Supported keywords: pug, corgi, cat, giraffe",
			"/count/:keyword":          "Number of images available",
			"/random/:keyword/:action": "Get a random image. Optional action: url (get the link directly), preview (preview the image)",
			"/bomb/:keyword/:number":   "Get a number of images. Default to 4",
		},
		"ABOUT": {
			"source": "http://github.com/zhengjia/awwimage",
		},
	}
}

func get_port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	return port
}

func set_api_key() {
	var err error
	config, err := ioutil.ReadFile("config")
	if err == nil {
		api_key = strings.TrimSpace(string(config))
	} else {
		api_key = os.Getenv("TUMBLR_KEY")
		if api_key == "" {
			panic_on_error(errors.New("TUMBLR_KEY isn't set"))
		}
	}
}

func (*tumblrFetcher) Fetch(url string) ([]byte, error) {
	var err error
	var resp *http.Response
	var body_bytes []byte
	resp, err = http.Get(url)
	log_on_error(err)
	body_bytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log_on_error(err)
	return body_bytes, err
}

func populate(kind string, fetcher FetcherInterface) {
	var timestamp int
	var url string
	var url_template string
	var err error
	var body_bytes []byte
	var tagged_api_response *TaggedApiResponse
	var results []string
	url_template = "http://api.tumblr.com/v2/tagged?api_key=" + api_key + "&tag=" + kind
	for len(results) < image_limit {
		if timestamp == 0 {
			url = url_template
		} else {
			url = url_template + "&before=" + strconv.Itoa(timestamp)
		}
		body_bytes, err = fetcher.Fetch(url)
		if err != nil {
			continue
		}
		err = json.Unmarshal(body_bytes, &tagged_api_response)
		log_on_error(err)
		for _, Blog := range tagged_api_response.Blogs {
			timestamp = Blog.Timestamp
			for _, Photo := range Blog.Photos {
				results = append(results, Photo.OriginalPhoto.Url)
			}
		}
	}
	image_mapping[kind] = results
}

func populate_mapping() {
	last_refreshed_at = time.Now()
	fetcher := tumblrFetcher{}
	for _, kind := range kinds {
		go populate(kind, &fetcher)
	}
}

func refresh_every_hour() {
	if time.Now().Sub(last_refreshed_at) > time.Minute*60 {
		populate_mapping()
	}
}

func check_image_presence(kind string, done chan bool) {
	for len(image_mapping[kind]) == 0 {
		time.Sleep(time.Second)
	}
	done <- true
}

func check_kind(kind string) (err error) {
	for _, k := range kinds {
		if k == kind {
			return
		}
	}
	err = errors.New("Image type not supported")
	return
}

func getHttpHandler() http.Handler {
	m := pat.New()
	m.Get("/", http.HandlerFunc(instruction))
	m.Get("/instruction", http.HandlerFunc(instruction))
	m.Get("/count/:kind", http.HandlerFunc(count))
	m.Get("/random/:kind", http.HandlerFunc(random))
	m.Get("/random/:kind/:action", http.HandlerFunc(random))
	m.Get("/bomb/:kind", http.HandlerFunc(bomb))
	m.Get("/bomb/:kind/:number", http.HandlerFunc(bomb))
	return m
}

func initialize() {
	set_api_key()
	populate_mapping()
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	initialize()
	handler := getHttpHandler()
	http.Handle("/", handler)
	http.ListenAndServe(":"+get_port(), nil)
}
