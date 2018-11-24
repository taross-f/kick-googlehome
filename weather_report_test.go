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
	fmt.Println("setUp()")
	os.Setenv("GO_ENV", "test")
	loadEnv()
}

func tearDown() {
	fmt.Println("tearDown()")
	os.Remove(os.Getenv("RAIN_FILE"))
}

func TestMain(m *testing.M) {
	setUp()
	result := m.Run()
	tearDown()
	os.Exit(result)
}

func TestAction_StartRaining(t *testing.T) {
	wr := NewWeatherReport()
	result := createRainingResult()

	got := wr.action(result)
	want := StartRaining
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestAction_KeepRaining(t *testing.T) {
	wr := NewWeatherReport()
	result := createRainingResult()

	wr.action(result)
	got := wr.action(result)
	want := Nothing
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestAction_StopRaining(t *testing.T) {
	wr := NewWeatherReport()
	result := createRainingResult()

	wr.action(result)
	result2 := createNotRainingResult()
	got := wr.action(result2)
	want := StopRaining
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestAction_KeepNotRaining(t *testing.T) {
	wr := NewWeatherReport()
	result := createNotRainingResult()

	wr.action(result)
	got := wr.action(result)
	want := Nothing
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestReadJson(t *testing.T) {
	createRainingResult()
}

func createRainingResult() *YahooWeatherResult {
	result := new(YahooWeatherResult)
	raw, err := ioutil.ReadFile("./testfiles/raining.json")
	print(raw)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(raw, result)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	return result
}

func createNotRainingResult() *YahooWeatherResult {
	result := new(YahooWeatherResult)
	raw, err := ioutil.ReadFile("./testfiles/notraining.json")
	print(raw)
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
