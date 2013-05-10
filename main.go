package main

import(
  "fmt"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
)

var valid_kinds = make([]string, 3)
var image_mapping = make(map[string][]string)
var api_key string

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
    "/": "help page",
    "count": "Total image count for this picture kind",
    "random": "Random",
    "bomb": "Photo bomb up to 10 images",
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

func visit(url string) {
  
}

func populate_image_mapping(kind string) {
  url_template := "http://api.tumblr.com/v2/tagged?api_key=" + api_key + "&tag=" + kind
  url := url_template
  resp, err := http.Get(url)
  body_bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
  check(err)
  body_string := string(body_bytes)
  fmt.Println(body_string)
  // image_mapping[kind] = api_body
  // err := json.Unmarshal(api_body, &u)
  // check(err)
}

func initialize(){
  set_api_key()
  kinds := []string{"pug", "corgi", "cat"}
  for _, kind := range kinds {
    image_mapping[kind] = []string{}
    populate_image_mapping(kind)
  }
}

func main() {
  initialize()
  http.HandleFunc("/", instruction)
  http.HandleFunc("/count", count)
  http.HandleFunc("/random", random)
  http.HandleFunc("/bomb", bomb)
  http.ListenAndServe(":" + get_port(), nil)
}
