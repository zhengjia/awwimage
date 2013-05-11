package main

import "testing"

func TestEndpoints(t *testing.T) {
  var endpoints = *Endpoints()
  var routes = []string{"/instruction", "/count/:keyword", "/random/:keyword", "/bomb/:keyword/:number", "/all/:keyword"}
  for _, route := range routes {
    if _, ok := endpoints[route] ; !ok {
      t.Errorf("Expected key %s", "instruction")
    }
  }  
}
