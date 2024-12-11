package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gtuk/discordwebhook"
	"github.com/hectormalot/omgo"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func getWeather(weather chan []discordwebhook.Field, e chan error) {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		e <- fmt.Errorf("Error loading .env file: %v", err)
		weather <- nil
		return
	}

	long, err := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	if err != nil {
		e <- fmt.Errorf("Error converting longitude:  %v", err)
		weather <- nil
		return
	}

	lat, err := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	if err != nil {
		e <- fmt.Errorf("Error converting latitude: %v", err)
		weather <- nil
		return
	}

	weatherDescriptions := map[int]string{}

	// Read the JSON file
	file, err := os.ReadFile("weather_description.json")
	if err != nil {
		e <- fmt.Errorf("Error reading JSON file: %v", err)
		weather <- nil
		return
	}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(file, &weatherDescriptions)
	if err != nil {
		e <- fmt.Errorf("Error unmarshalling JSON data: %v", err)
		weather <- nil
		return
	}

	c, err := omgo.NewClient()
	if err != nil {
		e <- fmt.Errorf("Error creating OpenMeteo client: %v", err)
		weather <- nil
		return
	}

	loc, err := omgo.NewLocation(lat, long)
	if err != nil {
		e <- fmt.Errorf("Error creating location: %v", err)
		weather <- nil
		return
	}

	opts := omgo.Options{
		DailyMetrics: []string{"temperature_2m_max", "temperature_2m_min", "weathercode"},
	}

	res, err := c.Forecast(context.Background(), loc, &opts)
	if err != nil {
		e <- fmt.Errorf("Error getting forecast: %v", err)
		weather <- nil
		return
	}

	weather <- []discordwebhook.Field{
		parseWeather(res, 0, weatherDescriptions),
		parseWeather(res, 1, weatherDescriptions),
	}
	e <- nil
}

func parseWeather(forecast *omgo.Forecast, day int, weatherDesc map[int]string) discordwebhook.Field {
	weatherType := int(forecast.DailyMetrics["weathercode"][day])
	weatherEmoji := ""
	switch weatherType {
	case 0, 1:
		weatherEmoji = "☀️"
		break
	case 2, 3:
		weatherEmoji = "☁️"
		break
	case 45, 48:
		weatherEmoji = "🌁"
		break
	case 51, 53, 55, 56, 57, 61, 63, 65, 66, 67:
		weatherEmoji = "🌧️"
		break
	case 71, 73, 75, 77, 85, 86:
		weatherEmoji = "❄️"
		break
	case 98, 96, 99:
		weatherEmoji = "⛈️"
		break
	}

	weatherDescription := weatherDesc[weatherType]

	title := weatherEmoji + " " + weatherDescription

	maxTemp := forecast.DailyMetrics["temperature_2m_max"][day]
	minTemp := forecast.DailyMetrics["temperature_2m_min"][day]
	averageTemp := (maxTemp + minTemp) / 2

	content := "\n🌡️ " + fmt.Sprintf("%.2f", averageTemp) + "°C"
	content += "\n🔺 " + fmt.Sprintf("%.2f", maxTemp) + "°C"
	content += "\n🔻 " + fmt.Sprintf("%.2f", minTemp) + "°C"

	inline := true

	field := discordwebhook.Field{
		Name:   &title,
		Value:  &content,
		Inline: &inline,
	}

	return field
}
