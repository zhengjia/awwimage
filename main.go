package main

import(
  "fmt"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/bmizerany/pat"
)

var limit = 100
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

func instruction(res http.ResponseWriter, req *http.Request) {
  result, err := json.Marshal( *Endpoints() )
  check(err)
  fmt.Fprint(res, string(result) )
}

func count(res http.ResponseWriter, req *http.Request) {
  var kind = req.FormValue("kind")
  fmt.Fprint(res, map[string]int{"count": len(image_mapping[kind])} )
}

func random(res http.ResponseWriter, req *http.Request) {
}

func bomb(res http.ResponseWriter, req *http.Request) {
}

// helper methods
func check(err error) {
  if err != nil {
    panic(err)
  }
}

func Endpoints() *map[string]string {
  return &map[string]string{
    "/instruction": "Get a random image. Supported query keywords: pug, corgi, cat",
    "/count/:keyword": "Total images available",
    "/random/:keyword": "Get a random image",
    "/bomb/:keyword/:number": "Get up to 10 images",
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

func populate_image_mapping(kind string) {
  var timestamp int32
  var url string
  var url_template string
  var err error
  var body_bytes []byte
  var tagged_api_response *TaggedApiResponse
  url_template = "http://api.tumblr.com/v2/tagged?api_key=" + api_key + "&tag=" + kind
  for len(image_mapping[kind] < limit) {
    if timestamp == 0 {
      url = url_template
    } else {
      url = url_template + "&timestamp=" + string(timestamp)
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

func initialize(){
  set_api_key()
  kinds := []string{"pug", "corgi", "cat"}
  for _, kind := range kinds {
    image_mapping[kind] = []string{}
    populate_image_mapping(kind)
  }
  fmt.Println(image_mapping)
}

func main() {
  initialize()
  m := pat.New()
  m.Get("/", http.HandlerFunc(instruction))
  m.Get("/count/:kind", http.HandlerFunc(count))
  m.Get("/random/:kind", http.HandlerFunc(random))
  m.Get("/bomb/:kind/:number", http.HandlerFunc(bomb))
  http.Handle("/", m)
  
  http.HandleFunc("/instruction", instruction)
  http.ListenAndServe(":" + get_port(), nil)
}
