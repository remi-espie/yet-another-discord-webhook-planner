package main

import (
	ics "github.com/arran4/golang-ical"
	"github.com/gtuk/discordwebhook"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// URL of the .ics file in .env
	icsUrl := os.Getenv("ICS_URL")

	// webhook URL in .env
	webhook := os.Getenv("WEBHOOK_URL")

	//set the username and avatar of the bot
	username := "📆 Planning Bot"
	avatar := os.Getenv("AVATAR_URL")

	// set the time location to Europe/Paris
	location, err := time.LoadLocation("Europe/Paris")

	// Prepare information for the CROUS meal
	sendMenu := false
	crousRestaurantId := os.Getenv("CROUS_RESTAURANT_ID")

	if err != nil {
		log.Fatalf("Error loading time")
	}

	today := time.Now().In(location)
	tomorrow := today.AddDate(0, 0, 1)

	var todayCourse []*ics.VEvent
	var todayEvent []*ics.VEvent

	var tomorrowCourse []*ics.VEvent
	var tomorrowEvent []*ics.VEvent

	weatherChan := make(chan []discordwebhook.Field)

	go getWeather(weatherChan)

	// Fetch the calendar
	cal := getCal(icsUrl)

	// store the events of today and tomorrow
	for _, event := range cal.Events() {
		// Print the event as JSON
		at, err := event.GetStartAt()
		if err != nil {
			return
		}

		if DateEqual(at, today) {
			category := event.GetProperty("CATEGORIES").Value
			if category == "Cours" {
				todayCourse = append(todayCourse, event)
			}

			if category == "Important" {
				todayEvent = append(todayEvent, event)
			}
		}

		if DateEqual(at, tomorrow) {
			category := event.GetProperty("CATEGORIES").Value
			if category == "Cours" {
				tomorrowCourse = append(tomorrowCourse, event)
			}

			if category == "Important" {
				tomorrowEvent = append(tomorrowEvent, event)
			}
		}
	}

	var embeds []discordwebhook.Embed
	weather := <-weatherChan

	if len(todayCourse) > 0 || len(todayEvent) > 0 {
		embeds = append(embeds, getEmbed(todayCourse, todayEvent, "today", weather[0]))
		sendMenu = true
	}

	if len(tomorrowCourse) > 0 || len(tomorrowEvent) > 0 {
		embeds = append(embeds, getEmbed(tomorrowCourse, tomorrowEvent, "tomorrow", weather[1]))
	}

	if sendMenu && crousRestaurantId != "" {
		id, err := strconv.Atoi(crousRestaurantId)
		if err != nil {
			log.Fatalf("Error converting crousRestaurantId to integer")
		}
		embeds = append(embeds, getMenuEmbed(id))
	}

	sendMessage(webhook, username, avatar, embeds)
}

func sendMessage(webhook string, username string, avatar string, embed []discordwebhook.Embed) {
	message := discordwebhook.Message{
		Username:  &username,
		AvatarUrl: &avatar,
		Embeds:    &embed,
	}

	err := discordwebhook.SendMessage(webhook, message)
	if err != nil {
		log.Fatal(err)
	}
}

func getCal(icsUrl string) *ics.Calendar {
	resp, err := http.Get(icsUrl)
	if err != nil {
		log.Fatal("Error fetching .ics file:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Error closing response body:", err)
		}
	}(resp.Body)

	// Read the .ics file contents
	icsData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading .ics file:", err)
	}

	// Parse the .ics data
	cal, err := ics.ParseCalendar(strings.NewReader(string(icsData)))
	if err != nil {
		log.Fatal("Error parsing .ics data:", err)
	}

	return cal
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return d1 == d2 && m1 == m2 && y1 == y2
}
