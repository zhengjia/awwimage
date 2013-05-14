package main

import "testing"

func TestEndpointsContainAllRoutes(t *testing.T) {
  var endpoints = *Endpoints()
  var routes = []string{"/instruction", "/count/:keyword", "/random/:keyword", "/bomb/:keyword/:number", "/all/:keyword"}
  for _, route := range routes {
    if _, ok := endpoints[route] ; !ok {
      t.Errorf("Expected key %s", route)
    }
  }  
}

func TestGetJsonString(t *testing.T) {
	json := GetJsonString(&map[string][]string{"urls": []string{"b", "a"}} )
	if json != "{\"urls\":[\"b\",\"a\"]}" {
	  t.Errorf("Returned %s", json)
	}
}
