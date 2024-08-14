package Groupie_tracker

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type ArtistWithLocation struct {
	JsonData   interface{}
	ArtistData Artist
}

func GetDataFromJson(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	artisData := changeJsonToStruct()
	mpt, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error 500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	mpt.Execute(w, artisData)
}

func Handler_Show_Relation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if r.URL.Path != "/Artist/"+id {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	artist, err := Fetch_Data_Relation_FromId(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Artist not found: %v", err), http.StatusNotFound)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/InforArtis.html"))
	err = tmpl.Execute(w, artist)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func HandleStyle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/styles"):]
	fullpath := filepath.Join("src", path)
	fileinfo, err := os.Stat(fullpath)

	if !os.IsNotExist(err) && !fileinfo.IsDir() {
		http.StripPrefix("/styles", http.FileServer(http.Dir("src"))).ServeHTTP(w, r)
	} else {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}
