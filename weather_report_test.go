package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func setUp() {
	os.Setenv("GO_ENV", "test")
}
func TestAction(t *testing.T) {
	setUp()
	wr := NewWeatherReport()
	result := createRainingResult()
	fmt.Printf("%v\n", result)
	out := wr.action(result)
	expected := StartRaining
	if out != expected {
		log.Fatal("error")
	}
}

func createRainingResult() *YahooWeatherResult {
	result := new(YahooWeatherResult)
	raw, err := ioutil.ReadFile("./testfiles/raining.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(raw, result)
	if err != nil {
		print(err)
		return nil
	}
	return result
}

func createNotRainingResult() *YahooWeatherResult {
	result := new(YahooWeatherResult)
	raw, err := ioutil.ReadFile("./testfiles/notraining.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(raw, result)
	if err != nil {
		print(err)
		return nil
	}
	return result
}
