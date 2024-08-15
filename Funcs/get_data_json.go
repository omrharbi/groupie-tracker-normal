package Groupie_tracker

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	// "strconv"
	"text/template"
)

type ArtistWithLocation struct {
	JsonData interface{}
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
	artisData, err1 := changeJsonToStruct()
	if err1 != nil {
		HandleErrors(w, errors.BadRequest, errors.DescriptionBadRequest, http.StatusBadRequest)
		return
	}
	err := tmpl.ExecuteTemplate(w, "index.html", artisData)
	if err != nil {
		HandleErrors(w, errors.InternalError, errors.DescriptionInternalError, http.StatusInternalServerError)
		return
	}
}

func HandlerShowRelation(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	// id, err := strconv.Atoi(idParam)
	// if err != nil || id <= 0 {
	// 	HandleErrors(w, errors.BadRequest, errors.DescriptionBadRequest, http.StatusBadRequest)
	// 	return
	// }
	// if r.URL.Path != "/Artist/"+idParam {
	// 	HandleErrors(w, errors.NotFound, errors.DescriptionNotFound, http.StatusNotFound)
	// 	return
	// }
	if r.Method != http.MethodGet {
		HandleErrors(w, errors.InternalError, errors.DescriptionInternalError, http.StatusMethodNotAllowed)
		return
	}
	fmt.Println(r.Method)
	artist, err := FetchDataRelationFromId(idParam)
	if err != nil {
		HandleErrors(w, errors.InternalError, errors.DescriptionInternalError, http.StatusInternalServerError)
		return
	}
	errs := tmpl.ExecuteTemplate(w, "InforArtis.html", artist)
	if errs != nil {
		HandleErrors(w, errors.InternalError, errors.DescriptionInternalError, http.StatusInternalServerError)

		return
	}
}

func HandleStyle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/styles"):]
	fullpath := filepath.Join("src", path)
	fileinfo, err := os.Stat(fullpath)

	if !os.IsNotExist(err) && !fileinfo.IsDir() {
		http.StripPrefix("/styles", http.FileServer(http.Dir("src"))).ServeHTTP(w, r)
	} else {
		HandleErrors(w, errors.InternalError, errors.DescriptionInternalError, http.StatusInternalServerError)
		return
	}
}

func HandleErrors(w http.ResponseWriter, message, description string, code int) {
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
