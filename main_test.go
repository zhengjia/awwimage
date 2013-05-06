package main

import "testing"

func TestEndpoints(t *testing.T) {
  var endpoints = *Endpoints()
  var routes = []string{"/", "count", "random", "bomb"}
  for _, route := range routes {
    if _, ok := endpoints[route] ; !ok {
      t.Errorf("Expected key %s", "instruction")
    }
  }  
}
