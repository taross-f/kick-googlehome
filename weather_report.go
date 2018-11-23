package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./envfiles/develop.env")
	if err != nil {
		print("Error loading .env file")
	}
	wr := NewWeatherReport()
	wr.Report()
}

func loadEnv() {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "develop")
	}
	err := godotenv.Load("./envfiles/%s.env", os.Getenv("GO_ENV"))
	if err != nil {
		log.Fatal("Can't load env")
	}
}

const URL = "https://map.yahooapis.jp/weather/V1/place"

// WeatherReport reports weather
type WeatherReport struct {
	url      string
	appID    string
	latLon   string
	rainFile string
}

// NewWeatherReport creates a new Weather report
func NewWeatherReport() *WeatherReport {
	wr := &WeatherReport{
		appID:    os.Getenv("APP_ID"),
		url:      URL,
		latLon:   os.Getenv("LATLON"),
		rainFile: os.Getenv("RAIN_FILE"),
	}
	return wr
}

// Report weather
func (p *WeatherReport) Report() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprint(URL, "?output=json&coordinates=", p.latLon, "&appid=", p.appID), nil)
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
	p.action(result)
}

// State represents raining status
type State int

// State type
const (
	Nothing State = iota
	StartRaining
	StopRaining
)

func (p *WeatherReport) action(result *YahooWeatherResult) State {
	nowRainfall := result.Feature[0].Property.WeatherList.Weather[0].Rainfall
	if p.isRaining() {
		if nowRainfall == 0 {
			os.Remove(p.rainFile)
			p.alertStoppedRaining()
			return StopRaining
		}
	} else {
		if nowRainfall > 0 {
			os.Create(p.rainFile)
			p.alertStartedRaining()
			return StartRaining
		}
	}
	return Nothing
}

func (p *WeatherReport) alertStoppedRaining() {
	print("stop raining")
}

func (p *WeatherReport) alertStartedRaining() {
	print("start raining")
}

func (p *WeatherReport) isRaining() bool {
	_, err := os.Stat(p.rainFile)
	return err == nil
}
