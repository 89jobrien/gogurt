package tools

import "fmt"

type WeatherInput struct {
	City string `json:"city"`
}

func GetWeather(input WeatherInput) (string, error) {
	if input.City == "New York" {
		return "The weather in New York is sunny.", nil
	}
	return fmt.Sprintf("I don't know the weather for %s", input.City), nil
}