package main

import (
	ics "github.com/arran4/golang-ical"
	"github.com/gtuk/discordwebhook"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
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
	url := os.Getenv("ICS_URL")

	// webhook URL in .env
	webhook := os.Getenv("WEBHOOK_URL")

	// Fetch the .ics file
	resp, err := http.Get(url)
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

	// set the time location to Europe/Paris
	location, err := time.LoadLocation("Europe/Paris")

	today := time.Now().In(location)
	tomorrow := today.AddDate(0, 0, 1)

	var todayCourse []*ics.VEvent
	var todayEvent []*ics.VEvent

	var tomorrowCourse []*ics.VEvent
	var tomorrowEvent []*ics.VEvent

	// store the events of today and tomorrow
	for _, event := range cal.Events() {
		// Print the event as JSON
		at, err := event.GetStartAt()
		if err != nil {
			return
		}
		if DateEqual(at, today) {
			if event.GetProperty("CATEGORIES").Value == "Cours" {
				todayCourse = append(todayCourse, event)
			}

			if event.GetProperty("CATEGORIES").Value == "Important" {
				todayEvent = append(todayEvent, event)
			}
		}

		if DateEqual(at, tomorrow) {
			if event.GetProperty("CATEGORIES").Value == "Cours" {
				tomorrowCourse = append(tomorrowCourse, event)
			}

			if event.GetProperty("CATEGORIES").Value == "Important" {
				tomorrowEvent = append(tomorrowEvent, event)
			}
		}
	}

	//set the username and avatar of the bot
	username := "📆 Planning Bot"
	avatar := os.Getenv("AVATAR_URL")

	if len(todayCourse) > 0 || len(todayEvent) > 0 {
		content1 := getOrderedPlanning(todayCourse, todayEvent, "today")
		sendMessage(webhook, username, avatar, content1)
	}

	if len(tomorrowCourse) > 0 || len(tomorrowEvent) > 0 {
		content2 := getOrderedPlanning(tomorrowCourse, tomorrowEvent, "tomorrow")
		sendMessage(webhook, username, avatar, content2)
	}
}

func sendMessage(webhook string, username string, avatar string, content string) {
	message := discordwebhook.Message{
		Username:  &username,
		Content:   &content,
		AvatarUrl: &avatar,
	}

	err := discordwebhook.SendMessage(webhook, message)
	if err != nil {
		log.Fatal(err)
	}
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func getOrderedPlanning(course, event []*ics.VEvent, day string) string {

	sort.Slice(course, func(i, j int) bool {
		first, err := course[i].GetStartAt()
		if err != nil {
			return false
		}
		second, err := course[j].GetStartAt()
		if err != nil {
			return false
		}
		return first.Before(second)
	})

	location, err := time.LoadLocation("Europe/Paris")

	content := "----------------------\n📅 " + strings.Title(day) + "'s planning !\n\n"

	firstCourse, err := course[0].GetStartAt()
	if err != nil {
		return ""
	}
	endCourse, err := course[len(course)-1].GetEndAt()
	if err != nil {
		return ""
	}
	content += "🔋 Start of " + day + " : **" + firstCourse.In(location).Format("15:04") + "**\n\n" +
		"\U0001FAAB End of " + day + " : **" + endCourse.In(location).Format("15:04") + "**"

	if len(event) > 0 {
		content += "\n\n⚠️ **" + strconv.Itoa(len(event)) + "** important events " + day + " !"
		for _, event := range event {
			start, err := event.GetStartAt()
			if err != nil {
				return ""
			}
			end, err := event.GetEndAt()
			if err != nil {
				return ""
			}
			content += "\n**" + start.In(location).Format("15:04") + "** → **" + end.In(location).Format("15:04") + "** : " + event.GetProperty("SUMMARY").Value
		}
	}

	content += "\n\n**" + strconv.Itoa(len(course)) + "** courses " + day + " !"

	for _, event := range course {
		start, err := event.GetStartAt()
		if err != nil {
			return ""
		}
		end, err := event.GetEndAt()
		if err != nil {
			return ""
		}
		content += "\n**" + start.In(location).Format("15:04") + "** → **" + end.In(location).Format("15:04") + "** : " + event.GetProperty("SUMMARY").Value
	}

	return content
}
