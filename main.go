package main

import(
  "fmt"
	"net/http"
	"os"
)

type ImagemiServer struct{}

func (h ImagemiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello!")
}

func main() {
    var port = os.Getenv("PORT")
    if port == "" {
      port = "4000"
    }
    var h ImagemiServer
    http.ListenAndServe(":" + port, h)
}
