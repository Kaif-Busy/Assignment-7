package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

func getNumber(input string) int64 {
	number := ""
	for _, character := range input {
		if character >= '0' && character <= '9' {
			number += string(character)
		}
	}
	finalNumber, _ := strconv.ParseInt(number, 10, 64)
	return finalNumber
}

type WeatherData struct {
	Name    string   `json:"name"`
	Weather string   `json:"weather"`
	Status  []string `json:"status"`
}

type APIResponse struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	Total      int           `json:"total"`
	TotalPages int           `json:"total_pages"`
	Data       []WeatherData `json:"data"`
}

func getRequest(url string) APIResponse {
	var apiResponse APIResponse
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return apiResponse
	}
	//defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", response.Status)
		return apiResponse
	}

	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return apiResponse
	}
	return apiResponse
}

func getWeather(url string) [][]interface{} {
	var finalResult [][]interface{}
	apiResponse := getRequest(url)
	urlPage := "&page="

	var wg sync.WaitGroup
	var mu sync.Mutex

	for page := 1; page <= apiResponse.TotalPages; page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			newUrl := url + urlPage + strconv.Itoa(page)
			apiResponse := getRequest(newUrl)

			mu.Lock()
			defer mu.Unlock()

			for _, weatherData := range apiResponse.Data {
				var cityData []interface{}

				name := weatherData.Name
				temperature := getNumber(weatherData.Weather)
				wind := getNumber(weatherData.Status[0])
				humidity := getNumber(weatherData.Status[1])

				cityData = append(cityData, name, temperature, wind, humidity)

				finalResult = append(finalResult, cityData)
			}
		}(page)
	}

	wg.Wait()
	return finalResult
}

func main() {
	var searchTerm string
	weatherAPIURL := "https://jsonmock.hackerrank.com/api/weather/search?name="
	fmt.Println("Enter the search term")
	fmt.Scan(&searchTerm)
	weatherAPIURL += searchTerm
	weatherDetails := getWeather(weatherAPIURL)
	for _, weatherDetail := range weatherDetails {
		fmt.Println(weatherDetail)
	}

}
