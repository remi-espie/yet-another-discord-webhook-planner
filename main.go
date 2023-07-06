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
	"unicode"
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

	if err != nil {
		log.Fatalf("Error loading time")
	}

	today := time.Now().In(location)
	tomorrow := today.AddDate(0, 0, 1)

	var todayCourse []*ics.VEvent
	var todayEvent []*ics.VEvent

	var tomorrowCourse []*ics.VEvent
	var tomorrowEvent []*ics.VEvent

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

	weather := getWeather()

	var embeds []discordwebhook.Embed

	if len(todayCourse) > 0 || len(todayEvent) > 0 {
		embeds = append(embeds, getEmbed(todayCourse, todayEvent, "today", weather[0]))
	}

	if len(tomorrowCourse) > 0 || len(tomorrowEvent) > 0 {
		embeds = append(embeds, getEmbed(tomorrowCourse, tomorrowEvent, "tomorrow", weather[1]))
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

func getEmbed(course, event []*ics.VEvent, day string, weather discordwebhook.Field) discordwebhook.Embed {

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

	dayTitle := []rune(day)
	dayTitle[0] = unicode.ToUpper(dayTitle[0])
	title := "📅 " + string(dayTitle) + "'s planning !"

	firstCourse, err := course[0].GetStartAt()
	if err != nil {
		return discordwebhook.Embed{}
	}
	endCourse, err := course[len(course)-1].GetEndAt()
	if err != nil {
		return discordwebhook.Embed{}
	}
	content := "🔋 Start of " + day + " : **" + firstCourse.In(location).Format("15:04") + "**\n\n" +
		"\U0001FAAB End of " + day + " : **" + endCourse.In(location).Format("15:04") + "**"

	titleField := discordwebhook.Field{
		Name:   &title,
		Value:  &content,
		Inline: nil,
	}

	importantField := discordwebhook.Field{}
	importantTitle := ""

	if len(event) > 0 {
		if len(event) > 1 {
			importantTitle = "⚠️ **" + strconv.Itoa(len(event)) + "** important events " + day + " !"
		} else {
			importantTitle = "⚠️ **1** important event " + day + " !"
		}
		importantContent := ""
		for _, event := range event {
			start, err := event.GetStartAt()
			if err != nil {
				return discordwebhook.Embed{}
			}
			end, err := event.GetEndAt()
			if err != nil {
				return discordwebhook.Embed{}
			}
			importantContent += "\n**" + start.In(location).Format("15:04") + "** → **" + end.In(location).Format("15:04") + "** : " + event.GetProperty("SUMMARY").Value
		}

		importantField = discordwebhook.Field{
			Name:   &importantTitle,
			Value:  &importantContent,
			Inline: nil,
		}
	}

	courseTitle := "**" + strconv.Itoa(len(course)) + "** courses " + day + " !"
	courseContent := ""

	for _, event := range course {
		start, err := event.GetStartAt()
		if err != nil {
			return discordwebhook.Embed{}
		}
		end, err := event.GetEndAt()
		if err != nil {
			return discordwebhook.Embed{}
		}
		courseContent += "\n**" + start.In(location).Format("15:04") + "** → **" + end.In(location).Format("15:04") + "** : " + event.GetProperty("SUMMARY").Value
	}

	inline := true

	thumbnail := ""

	if day == "today" {
		thumbnail = "https://calendar.cluster-2022-2.dopolytech.fr/calendar?locale=" + location.String() + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + "&size=500"
	} else {
		thumbnail = "https://calendar.cluster-2022-2.dopolytech.fr/calendar?locale=" + location.String() + "&timestamp=" + strconv.FormatInt(time.Now().AddDate(0, 0, 1).Unix(), 10) + "&size=500"
	}

	thumbnailField := discordwebhook.Thumbnail{Url: &thumbnail}

	courseField := discordwebhook.Field{
		Name:   &courseTitle,
		Value:  &courseContent,
		Inline: &inline,
	}

	var fields []discordwebhook.Field

	if importantField.Value != nil {
		fields = []discordwebhook.Field{
			titleField,
			importantField,
			courseField,
			weather,
		}
	} else {
		fields = []discordwebhook.Field{
			titleField,
			courseField,
			weather,
		}
	}

	footerText := "Made with ❤️ by @luckmk1 | " + time.Now().In(location).Format("2006-01-02 15h04:05 Z0700 MST")

	footer := discordwebhook.Footer{
		Text:    &footerText,
		IconUrl: &thumbnail,
	}

	embed := discordwebhook.Embed{
		Fields:    &fields,
		Thumbnail: &thumbnailField,
		Footer:    &footer,
	}

	return embed
}
