package Groupie_tracker

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

type ArtistWithLocation struct {
	JsonData   interface{}
	// ArtistData Artist
}

var (
	tmpl   *template.Template
	errors AllMessageErrors
)

// Initialize the global template variable
func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	errors = ErrorsMessage()
}

func GetDataFromJson(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleErrors(w, errors.NotFound, errors.DescriptionNotFound, http.StatusNotFound)
		return
	}
	artisData := changeJsonToStruct()
	mpt, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleErrors(w, errors.InternalError, errors.DescriptionInternalError, http.StatusInternalServerError)
		return
	}
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
	artist, err := Fetch_Data_Relation_FromId(idParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func HandleSearsh(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	artisData := changeJsonToStruct()
	result := Searsh_data(query, artisData)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func HandleStyle(w http.ResponseWriter, r *http.Request) {
	errors := ErrorsMessage()
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

func HandleErrors(w http.ResponseWriter, message, description string, code int) {
	// tmp1 := Tmp{}
	errorsMessage := Errors{
		Message:     message,
		Description: description,
		Code:        code,
	}
	w.WriteHeader(code)
	err := tmpl.ExecuteTemplate(w, "errors.html", errorsMessage)
	if err != nil {
		http.Error(w, "Error 500 Internal Server Error", http.StatusInternalServerError)
	}
}
