package main

import(
  "fmt"
	"net/http"
)

type ImagemiServer struct{}

func (h ImagemiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello!")
}

func main() {
    var h ImagemiServer
    http.ListenAndServe("localhost:4000", h)
}
