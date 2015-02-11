package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bmizerany/pat"
)

var image_limit = 300
var server_started_at = time.Now()
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

func GetJsonString(v interface{}) string {
	result, err := json.Marshal(v)
	check(err)
	return string(result)
}

func instruction(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, GetJsonString(Endpoints()))
}

func count(res http.ResponseWriter, req *http.Request) {
	kind := req.URL.Query().Get(":kind")
	fmt.Fprint(res, GetJsonString(&map[string]int{"count": len(image_mapping[kind])}))
}

func random(res http.ResponseWriter, req *http.Request) {
	populate_uptime()
	kind := req.URL.Query().Get(":kind")
	action := req.URL.Query().Get(":action")
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
	populate_uptime()
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

// helper methods
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Endpoints() *map[string]string {
	return &map[string]string{
		"/instruction":             "Get a random image. Supported keywords: pug, corgi, cat, giraffe",
		"/count/:keyword":          "Number of images available",
		"/random/:keyword/:action": "Get a random image. Optional action: url (get the link directly), preview (preview the image)",
		"/bomb/:keyword/:number":   "Get a number of images. Default to 4",
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
			check(errors.New("TUMBLR_KEY isn't set"))
		}
	}
}

func visit(url string) []byte {
	var err error
	var resp *http.Response
	var body_bytes []byte
	resp, err = http.Get(url)
	check(err)
	body_bytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	check(err)
	return body_bytes
}

func populate(kind string, count int) {
	var timestamp int
	var url string
	var url_template string
	var err error
	var body_bytes []byte
	var tagged_api_response *TaggedApiResponse
	var results []string
	url_template = "http://api.tumblr.com/v2/tagged?api_key=" + api_key + "&tag=" + kind
	for len(results) < count {
		if timestamp == 0 {
			url = url_template
		} else {
			url = url_template + "&before=" + strconv.Itoa(timestamp)
		}
		body_bytes = visit(url)
		err = json.Unmarshal(body_bytes, &tagged_api_response)
		check(err)
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
	for _, kind := range kinds {
		go populate(kind, image_limit)
	}
}

func populate_uptime() {
	if time.Now().Sub(server_started_at) > time.Minute*60 {
		populate_mapping()
	}
}

func initialize() {
	set_api_key()
	for _, kind := range kinds {
		populate(kind, 5)
	}
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	initialize()
	m := pat.New()
	m.Get("/", http.HandlerFunc(instruction))
	m.Get("/count/:kind", http.HandlerFunc(count))
	m.Get("/random/:kind", http.HandlerFunc(random))
	m.Get("/random/:kind/:action", http.HandlerFunc(random))
	m.Get("/bomb/:kind", http.HandlerFunc(bomb))
	m.Get("/bomb/:kind/:number", http.HandlerFunc(bomb))
	http.Handle("/", m)
	http.HandleFunc("/instruction", instruction)
	populate_mapping()
	http.ListenAndServe(":"+get_port(), nil)
}
