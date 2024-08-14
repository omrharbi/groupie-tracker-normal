package Groupie_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func changeJsonToStruct() ([]JsonData, error) {
	var artistData []JsonData

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	boderesponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(boderesponse, &artistData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return artistData, nil
}

func Fetch_Data_Relation_FromId(id string) (Artist, error) {
	artistsJson, err := Get_Data_From_Artists_With_ID(id)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from artis data %v", err)
	}

	ResponBody, err := http.Get(artistsJson.Relations)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from URL: %v", err)
	}

	defer ResponBody.Body.Close()
	var artist Artist
	body, err := io.ReadAll(ResponBody.Body)
	fmt.Println(string(body))
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
		// CreationDate:   artistsJson.ConcertDates,
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
	respoceArtists, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
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
	newartist := JsonData{
		Name:    artistsJson.Name,
		Image:   artistsJson.Image,
		Members: artistsJson.Members,
	}
	return newartist, nil
}

func CaptalizdString(s string) string {
	s = strings.Replace(s, "-", " ", -1)
	s = strings.Replace(s, "_", " ", -1)
	return strings.Title(s)
}

// func ErrorsMessage() AllMessageErrors {
// 	return AllMessageErrors{
// 		NotFound:                 "Page Not Found",
// 		BadRequest:               "Bad Request",
// 		InternalError:            "Internal Server Error",
// 		DescriptionNotFound:      "Sorry, the page you are looking for does not exist. It might have been moved or deleted. Please check the URL or return to the homepage.",
// 		DescriptionBadRequest:    "The request could not be understood by the server due to malformed syntax. Please verify your input and try again.",
// 		DescriptionInternalError: "An unexpected error occurred on the server. We are working to resolve this issue. Please try again later.",
// 	}
// }
