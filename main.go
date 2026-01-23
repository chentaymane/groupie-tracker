package main

import (
	"fmt"
	"net/http"

	zone "zone/handlers"
)

/*
   This Program is supposed to fetch data from a public API and displays information about music artists and bands.
*/

// The main function
func main() {

	http.HandleFunc("/static/", zone.HandleStatic)
	http.HandleFunc("/", zone.HandlerHome)
	http.HandleFunc("/artist/", zone.HandlerArtist)
	http.HandleFunc("/filter", zone.HandleFilter)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
