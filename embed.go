package main

import (
	ics "github.com/arran4/golang-ical"
	"github.com/gtuk/discordwebhook"
	"sort"
	"strconv"
	"time"
	"unicode"
)

func getEmbed(course, event []*ics.VEvent, day string, weather *discordwebhook.Field) discordwebhook.Embed {

	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return discordwebhook.Embed{}
	}

	dayTitle := []rune(day)
	dayTitle[0] = unicode.ToUpper(dayTitle[0])
	title := "📅 " + string(dayTitle) + "'s planning !"

	content := "No courses today !"

	if len(course) > 0 {
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

		firstCourse, err := course[0].GetStartAt()
		if err != nil {
			return discordwebhook.Embed{}
		}
		endCourse, err := course[len(course)-1].GetEndAt()
		if err != nil {
			return discordwebhook.Embed{}
		}
		content = "🔋 Start of " + day + " : **" + firstCourse.In(location).Format("15:04") + "**\n\n" +
			"\U0001FAAB End of " + day + " : **" + endCourse.In(location).Format("15:04") + "**"
	}

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
		// add location if available
		if event.GetProperty("LOCATION") != nil {
			courseContent += "\n**" + start.In(location).Format("15:04") + "** → **" + end.In(location).Format("15:04") + "** : " + event.GetProperty("SUMMARY").Value + " | " + event.GetProperty("LOCATION").Value
		} else {
			courseContent += "\n**" + start.In(location).Format("15:04") + "** → **" + end.In(location).Format("15:04") + "** : " + event.GetProperty("SUMMARY").Value
		}
	}

	inline := true

	thumbnail := ""

	baseUrl := "https://calengo.espie.dev/calendar?locale=" + location.String() + "&size=200"

	if day == "today" {
		thumbnail = baseUrl + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10)
	} else {
		thumbnail = baseUrl + "&timestamp=" + strconv.FormatInt(time.Now().AddDate(0, 0, 1).Unix(), 10)
	}

	thumbnailField := discordwebhook.Thumbnail{Url: &thumbnail}

	courseField := discordwebhook.Field{
		Name:   &courseTitle,
		Value:  &courseContent,
		Inline: &inline,
	}

	fields := []discordwebhook.Field{
		titleField,
	}

	if importantField.Value != nil {
		fields = append(fields, importantField)
	}

	fields = append(fields, courseField)

	if weather != nil {
		fields = append(fields, *weather)
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
