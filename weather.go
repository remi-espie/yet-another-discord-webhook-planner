package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gtuk/discordwebhook"
	"github.com/hectormalot/omgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func getWeather(weather chan []discordwebhook.Field) {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	long, err := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	if err != nil {
		log.Fatal("Error converting longitude:", err)
	}

	lat, err := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	if err != nil {
		log.Fatal("Error converting latitude:", err)
	}

	location := os.Getenv("LOCATION")

	weatherDescriptions := map[int]string{}

	// Read the JSON file
	file, err := os.ReadFile("weather_description.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(file, &weatherDescriptions)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}

	c, err := omgo.NewClient()
	if err != nil {
		log.Fatalf("Error creating OpenMeteo client: %v", err)
	}

	loc, err := omgo.NewLocation(lat, long)
	if err != nil {
		log.Fatalf("Error creating location: %v", err)
	}

	println(location)

	opts := omgo.Options{
		Timezone:     location,
		DailyMetrics: []string{"temperature_2m_max", "temperature_2m_min", "weathercode"},
	}

	res, err := c.Forecast(context.Background(), loc, &opts)
	if err != nil {
		log.Fatalf("Error getting forecast: %v", err)
	}

	weather <- []discordwebhook.Field{
		parseWeather(res, 0, location, weatherDescriptions),
		parseWeather(res, 1, location, weatherDescriptions),
	}
}

func parseWeather(forecast *omgo.Forecast, day int, location string, weatherDesc map[int]string) discordwebhook.Field {
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
