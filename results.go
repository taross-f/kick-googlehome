package main

// YahooWeatherResult is results from yolp api
type YahooWeatherResult struct {
	Feature []struct {
		Geometry struct {
			Coordinates string `json:"Coordinates"`
			Type        string `json:"Type"`
		} `json:"Geometry"`
		ID       string `json:"Id"`
		Name     string `json:"Name"`
		Property struct {
			WeatherAreaCode int64 `json:"WeatherAreaCode"`
			WeatherList     struct {
				Weather []struct {
					Date     string  `json:"Date"`
					Rainfall float64 `json:"Rainfall"`
					Type     string  `json:"Type"`
				} `json:"Weather"`
			} `json:"WeatherList"`
		} `json:"Property"`
	} `json:"Feature"`
	ResultInfo struct {
		Copyright   string  `json:"Copyright"`
		Count       int64   `json:"Count"`
		Description string  `json:"Description"`
		Latency     float64 `json:"Latency"`
		Start       int64   `json:"Start"`
		Status      int64   `json:"Status"`
		Total       int64   `json:"Total"`
	} `json:"ResultInfo"`
}
