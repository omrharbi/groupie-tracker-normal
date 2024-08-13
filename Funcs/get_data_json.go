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
	mpt, err := template.ParseFiles("templats/index.html")
	if err != nil {
		http.Error(w, "Error 500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	mpt.Execute(w, artisData)
}

func ShowArtistHandler(w http.ResponseWriter, r *http.Request) {
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
	tmpl := template.Must(template.ParseFiles("templats/InforArtis.html"))
	err = tmpl.Execute(w, artist)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func HandleSearsh(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	artisData := changeJsonToStruct()
	result := Searsh_data(query, artisData)
	tmpl := template.Must(template.ParseFiles("templats/index.html"))
	err := tmpl.Execute(w, result)
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
	}else{
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}
