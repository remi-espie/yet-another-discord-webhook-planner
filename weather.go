package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"unicode"
)

func getWeather() string {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("OWM_KEY")

	forecast, err := owm.NewOneCall("C", "EN", apiKey, []string{owm.ExcludeHourly, owm.ExcludeCurrent})
	if err != nil {
		log.Fatal("Error get api data:", err)
	}

	long, err := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	if err != nil {
		log.Fatal("Error converting longitude:", err)
	}

	lat, err := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	if err != nil {
		log.Fatal("Error converting latitude:", err)
	}

	coord := &owm.Coordinates{
		Longitude: long,
		Latitude:  lat,
	}

	err = forecast.OneCallByCoordinates(coord)
	if err != nil {
		log.Fatal(err)
	}

	return parseWeather(forecast.Daily[0])
}

func parseWeather(forecast owm.OneCallDailyData) string {
	weatherType := forecast.Weather[0].Main
	weatherEmoji := ""
	switch weatherType {
	case "Clear":
		weatherEmoji = "☀️"
		break
	case "Clouds":
		weatherEmoji = "☁️"
		break
	case "Mist":
		weatherEmoji = "🌁"
		break
	case "Snow":
		weatherEmoji = "❄️"
		break
	case "Rain":
		weatherEmoji = "🌧️"
		break
	case "Drizzle":
		weatherEmoji = "🌧️"
		break
	case "Thunderstorm":
		weatherEmoji = "⛈️"
		break
	}

	weatherDescription := []rune(forecast.Weather[0].Description)
	weatherDescription[0] = unicode.ToUpper(weatherDescription[0])

	output := "## 🛰️ Forecasted weather:\n" + weatherEmoji + " " + string(weatherDescription)
	output += "\n🌡️ " + fmt.Sprintf("%.2f", forecast.Temp.Day) + "°C"
	output += "\n🔺 " + fmt.Sprintf("%.2f", forecast.Temp.Max) + "°C"
	output += "\n🔻 " + fmt.Sprintf("%.2f", forecast.Temp.Min) + "°C"

	return output
}
