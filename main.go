package main

import(
  "fmt"
	"net/http"
	"os"
	"encoding/json"
)

type ImagemiServer struct{}

func Endpoints() *map[string]string{
  return &map[string]string{
    "/": "Instruction",
    "count": "Total count of this picture type",
    "random": "Random",
    "bomb": "Photo bomb up to 10 images",
  }
}

func instruction(res http.ResponseWriter, req *http.Request){
  result, err := json.Marshal( *Endpoints() )
  if err != nil {
    fmt.Println(err)
  }
  fmt.Fprint(res, string(result) )
}

func main() {
    var port = os.Getenv("PORT")
    if port == "" {
      port = "4000"
    }
    // var h ImagemiServer
    http.HandleFunc("/", instruction)
    http.ListenAndServe(":" + port, nil)
}
