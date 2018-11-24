package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ikasamah/homecast"
	"github.com/joho/godotenv"
)

func main() {
	err := loadEnv()
	if err != nil {
		print("Error loading .env file")
	}
	wr := NewWeatherReport()
	wr.Report()
}

func loadEnv() error {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "develop")
	}
	err := godotenv.Load(fmt.Sprintf("./envfiles/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		log.Fatalf("Can't load env, %q", err)
	}
	return err
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
		log.Println(err)
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

func (s State) String() string {
	switch s {
	case Nothing:
		return "Nothing"
	case StartRaining:
		return "StartRaining"
	case StopRaining:
		return "StopRaining"
	}
	return "Unknown"
}

func (p *WeatherReport) action(result *YahooWeatherResult) State {
	nowRainfall := result.Feature[0].Property.WeatherList.Weather[0].Rainfall
	if p.isRaining() {
		if nowRainfall == 0 {
			err := os.Remove(p.rainFile)
			if err != nil {
				log.Printf("Cannot remove rainfile, %q", err)
				return Nothing
			}
			speak("このあたりのあめがやみました")
			return StopRaining
		}
	} else {
		if nowRainfall > 0 {
			_, err := os.Create(p.rainFile)
			if err != nil {
				log.Printf("Cannot create rainfile, %q", err)
				return Nothing
			}
			speak("このあたりであめがふりはじめました")
			return StartRaining
		}
	}
	return Nothing
}

func speak(s string) {
	if os.Getenv("GO_ENV") == "test" {
		log.Printf("speak: %s", s)
		return
	}
	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)

	for _, device := range devices {
		log.Printf("Device: [%s:%d]%s", device.AddrV4, device.Port, device.Name)

		// if err := device.Speak(ctx, s, "ja"); err != nil {
		// 	log.Printf("Failed to speak: %v", err)
		// }
	}
}

func (p *WeatherReport) isRaining() bool {
	_, err := os.Stat(p.rainFile)
	return err == nil
}
