// +build !appengine

package main

import(
  "fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/bmizerany/pat"
	"time"
	"math/rand"
	"strconv"
	"regexp"
)

var lower_limit = 300
var valid_kinds = make([]string, 3)
var image_mapping = make(map[string][]string)
var api_key string

type PhotoProperty struct {
  Url string
}

type Photo struct {
  OriginalPhoto PhotoProperty `json:"original_size"`
}

type Blog struct {
  Timestamp int
  Photos []Photo
}

type TaggedApiResponse struct {
  Blogs []Blog `json:"response"`
}

func GetJsonString(v interface{}) string{
  result, err := json.Marshal( v )
  check(err)
  return string(result)
}

func instruction(res http.ResponseWriter, req *http.Request) {
  fmt.Fprint(res, GetJsonString( Endpoints()) )
}

func count(res http.ResponseWriter, req *http.Request) {
  kind := req.URL.Query().Get(":kind")
  fmt.Fprint(res, GetJsonString(&map[string]int{"count": len(image_mapping[kind])} ) )
}

func random(res http.ResponseWriter, req *http.Request) {
  kind := req.URL.Query().Get(":kind")
  index := rand.Intn(len(image_mapping[kind]))
  fmt.Fprint(res, GetJsonString(&map[string]string{"url": image_mapping[kind][index]} ) )
}

func bomb(res http.ResponseWriter, req *http.Request) {
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
  fmt.Fprint(res, GetJsonString(&map[string][]string{"urls": result} ) )
}

func reload(w http.ResponseWriter, req *http.Request) {
  localhostRegex := regexp.MustCompile(`127.0.0.1`)
  if localhostRegex.Match([]byte( req.Host) ){
    populate_mapping()
  } else {
    http.NotFound(w, req)
  }
}

// helper methods
func check(err error) {
  if err != nil {
    panic(err)
  }
}

func Endpoints() *map[string]string {
  return &map[string]string{
    "/instruction": "Get a random image. Supported query keywords: pug, corgi, cat, giraffe",
    "/count/:keyword": "Total images available",
    "/random/:keyword": "Get a random image",
    "/bomb/:keyword/:number": "Get images. Default to 4 is number is not specified",
  }
}

func set_api_key() {
  var err error
  config, err := ioutil.ReadFile("config")
  api_key = string(config)
  check(err)
}

func visit(url string) []byte{
  var err error
  var resp *http.Response
  var body_bytes []byte
  
  resp, err = http.Get(url)
  body_bytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
  check(err)
  return body_bytes
}

func populate(kind string) {
  var timestamp int
  var url string
  var url_template string
  var err error
  var body_bytes []byte
  var tagged_api_response *TaggedApiResponse
  url_template = "http://api.tumblr.com/v2/tagged?api_key=" + api_key + "&tag=" + kind
  for len(image_mapping[kind]) < lower_limit {
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
        image_mapping[kind] = append(image_mapping[kind], Photo.OriginalPhoto.Url )
      }  
    }
  }
}

func populate_mapping() {
  kinds := []string{"pug", "corgi", "shiba", "cat", "giraffe",}
  for _, kind := range kinds {
    image_mapping[kind] = []string{}
    go populate(kind)
  }
}

func initialize(){
  set_api_key()
  populate_mapping()
  rand.Seed( time.Now().UTC().UnixNano())
}

func init() {
  initialize()
  m := pat.New()
  m.Get("/", http.HandlerFunc(instruction))
  m.Get("/count/:kind", http.HandlerFunc(count))
  m.Get("/random/:kind", http.HandlerFunc(random))
  m.Get("/bomb/:kind", http.HandlerFunc(bomb))
  m.Get("/bomb/:kind/:number", http.HandlerFunc(bomb))
  m.Get("/reload", http.HandlerFunc(reload) )
  http.Handle("/", m)
  http.HandleFunc("/instruction", instruction)
}
