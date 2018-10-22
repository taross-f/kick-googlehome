package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./envfiles/develop.env")
	if err != nil {
		print("Error loading .env file")
	}
	result := request()
	action(result)
}

const URL = "https://map.yahooapis.jp/weather/V1/place"
const LATLON = "139.732293,35.663613"
const RAIN_FILE = "state/rain"

func request() *YahooWeatherResult {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprint(URL, "?output=json&coordinates=", LATLON, "&appid=", os.Getenv("APP_ID")), nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	result := new(YahooWeatherResult)
	err = json.Unmarshal(body, result)
	if err != nil {
		print(err)
	}

	fmt.Printf("%+v\n", result.Feature)
	return result
}

func action(result *YahooWeatherResult) {
	nowRainfall := result.Feature[0].Property.WeatherList.Weather[0].Rainfall
	print(nowRainfall)
	if isRaining() {
		if nowRainfall == 0 {
			os.Remove(RAIN_FILE)
			alertStoppedRaining()
		}
	} else {
		if nowRainfall > 0 {
			os.Create(RAIN_FILE)
			alertStartedRaining()
		}
	}
}

func alertStoppedRaining() {
	print("stop raining")
}

func alertStartedRaining() {
	print("start raining")
}

func isRaining() bool {
	return isExist(RAIN_FILE)
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
