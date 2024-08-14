package Groupie_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func changeJsonToStruct() []JsonData {
	var artistData []JsonData
	url := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error from Response")
	}

	if response.StatusCode == http.StatusOK {
		boderesponse, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Error from reading Response")
		}
		err = json.Unmarshal(boderesponse, &artistData)
		if err != nil {
			log.Fatal("Error To Unmarshal Data ")
		}

	}
	return artistData
}

func Fetch_Data_Relation_FromId(id string) (Artist, error) {
	artistsJson, err := Get_Data_From_Artists_With_ID(id)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from artis data %v", err)
	}
	url := "https://groupietrackers.herokuapp.com/api/relation/" + id
	ResponBody, err := http.Get(url)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from URL: %v", err)
	}
	defer ResponBody.Body.Close()
	var artist Artist
	body, err := io.ReadAll(ResponBody.Body)
	if err != nil {
		return Artist{}, fmt.Errorf("error reading response body: %v", err)
	}
	err = json.Unmarshal(body, &artist)
	if err != nil {
		return Artist{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	artist = Artist{
		ID:             artist.ID,
		Image:          artistsJson.Image,
		Name:           artistsJson.Name,
		DatesLocations: artist.DatesLocations,
		Members:        artistsJson.Members,
	}
	formatlocations := make(map[string][]string)

	for locatios, dates := range artist.DatesLocations {
		formaloca := CaptalizdString(locatios)
		formatlocations[formaloca] = dates
	}
	artist.DatesLocations = formatlocations
	return artist, nil
}

func Get_Data_From_Artists_With_ID(id string) (JsonData, error) {
	urlartists := "https://groupietrackers.herokuapp.com/api/artists/" + id
	respoceArtists, err := http.Get(urlartists)
	if err != nil {
		return JsonData{}, fmt.Errorf("error fetching data from URL: %v", err)
	}
	defer respoceArtists.Body.Close()

	var artistsJson JsonData
	bodyResponse_artists, err := io.ReadAll(respoceArtists.Body)
	if err != nil {
		return JsonData{}, fmt.Errorf("error reading response body: %v", err)
	}
	err = json.Unmarshal(bodyResponse_artists, &artistsJson)
	if err != nil {
		return JsonData{}, fmt.Errorf("error decoding artist data: %v", err)
	}
	return artistsJson, nil
}

func Searsh_data(search string, artists []JsonData) []JsonData {
	var result []JsonData
	for _, item := range artists {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(item.FirstAlbum), strings.ToLower(search)) {
			result = append(result, item)
		}
	}
	return result
}

func CaptalizdString(s string) string {
	s = strings.Replace(s, "-", " ", -1)
	s = strings.Replace(s, "_", " ", -1)
	return strings.Title(s)
}

func ErrorsMessage() AllMessageErrors {
	return AllMessageErrors{
		NotFound:                 "Page Not Found",
		BadRequest:               "Bad Request",
		InternalError:            "Internal Server Error",
		DescriptionNotFound:      "Sorry, the page you are looking for does not exist. It might have been moved or deleted. Please check the URL or return to the homepage.",
		DescriptionBadRequest:    "The request could not be understood by the server due to malformed syntax. Please verify your input and try again.",
		DescriptionInternalError: "An unexpected error occurred on the server. We are working to resolve this issue. Please try again later.",
	}
}
