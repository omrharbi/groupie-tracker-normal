package Groupie_tracker

import (
	"encoding/json"
	"fmt"
	"net/http"

)

func changeJsonToStruct() ([]JsonData, error) {
	var artistData []JsonData

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, fmt.Errorf("khata2 f jib l'data: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code machi normal: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&artistData)
	if err != nil {
		return nil, fmt.Errorf("khata2 f t7wil JSON: %v", err)
	}

	return artistData, nil
}

func FetchDataRelationFromId(id string) (Artist, error) {
    artistsJson, err := GetDataFromArtistsWithID(id)
    if err != nil {
        return Artist{}, fmt.Errorf("error fetching data from artist data: %w", err)
    }

    respBody, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
    if err != nil {
        return Artist{}, fmt.Errorf("error fetching data from URL: %w", err)
    }
    defer respBody.Body.Close()

    var artist Artist
    err = json.NewDecoder(respBody.Body).Decode(&artist)
    if err != nil {
        return Artist{}, fmt.Errorf("error decoding response: %w", err)
    }

    artist.Image = artistsJson.Image
    artist.Name = artistsJson.Name
    artist.Members = artistsJson.Members

    if artist.ID == 0 {
        return Artist{}, fmt.Errorf("invalid artist ID")
    }
	
    return artist, nil
}

func GetDataFromArtistsWithID(id string) (JsonData, error) {
    urlArtists := "https://groupietrackers.herokuapp.com/api/artists/" + id
    resp, err := http.Get(urlArtists)
    if err != nil {
        return JsonData{}, fmt.Errorf("error fetching data from URL: %w", err)
    }
    defer resp.Body.Close()

    var artistsJson JsonData
    err = json.NewDecoder(resp.Body).Decode(&artistsJson)
    if err != nil {
        return JsonData{}, fmt.Errorf("error decoding artist data: %w", err)
    }
    return artistsJson, nil
}

// func formatLocations(locations map[string][]string) map[string][]string {
//     formatted := make(map[string][]string, len(locations))
//     for location, dates := range locations {
//         formattedLoc := strings.Title(strings.NewReplacer("-", " ", "_", " ").Replace(location))
//         formatted[formattedLoc] = dates
//     }
//     return formatted
// }

// func CapitalizeString(s string) string {
//     return strings.Title(strings.NewReplacer("-", " ", "_", " ").Replace(s))
// }

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
