package Groupie_tracker

import (
	"fmt"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println(id)
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Error 500 Internal Server Error", http.StatusInternalServerError)
		}
	}()
	if r.URL.Path == "/" {
		GetDataFromJson(w, r)
	} else if r.URL.Path == "/Artist/"+id {
		ShowArtistHandler(w, r)
	} else {
		http.Error(w, "404 not Found", http.StatusNotFound)
	}
}
