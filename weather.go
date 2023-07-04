package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/gtuk/discordwebhook"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
)

func getWeather() []discordwebhook.Field {

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

	location := os.Getenv("LOCATION")

	coord := &owm.Coordinates{
		Longitude: long,
		Latitude:  lat,
	}

	err = forecast.OneCallByCoordinates(coord)
	if err != nil {
		log.Fatal(err)
	}

	return []discordwebhook.Field{
		parseWeather(forecast.Daily[0], forecast.Alerts, location),
		parseWeather(forecast.Daily[1], forecast.Alerts, location),
	}
}

func parseWeather(forecast owm.OneCallDailyData, alerts []owm.OneCallAlertData, location string) discordwebhook.Field {
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

	title := weatherEmoji + " " + string(weatherDescription)

	content := "\n🌡️ " + fmt.Sprintf("%.2f", forecast.Temp.Day) + "°C"
	content += "\n🔺 " + fmt.Sprintf("%.2f", forecast.Temp.Max) + "°C"
	content += "\n🔻 " + fmt.Sprintf("%.2f", forecast.Temp.Min) + "°C"

	if len(alerts) > 0 {

		location, err := time.LoadLocation(location)
		if err != nil {
			log.Fatalf("Error loading time")
		}

		for _, alert := range alerts {
			content += "\n🚨 " + alert.Event
			content += " from " + time.Unix(int64(alert.Start), 0).In(location).Format("15:04")
			content += " to " + time.Unix(int64(alert.End), 0).In(location).Format("15:04")
			content += " !"
		}
	}

	flag := true

	field := discordwebhook.Field{
		Name:   &title,
		Value:  &content,
		Inline: &flag,
	}

	return field
}
