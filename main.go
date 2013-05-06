package main

import(
  "fmt"
	"net/http"
	"os"
)

type ImagemiServer struct{}

func root(res http.ResponseWriter, req *http.Request){
  fmt.Fprint(res, "Hello!")
}

func main() {
    var port = os.Getenv("PORT")
    if port == "" {
      port = "4000"
    }
    // var h ImagemiServer
    http.HandleFunc("/", root)
    http.ListenAndServe(":" + port, nil)
}
