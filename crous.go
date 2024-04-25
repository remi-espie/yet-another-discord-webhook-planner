package main

import (
	"encoding/json"
	"github.com/gtuk/discordwebhook"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getMenuEmbed(crousRestaurantId int) discordwebhook.Embed {
	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return discordwebhook.Embed{}
	}

	thumbnail := "https://calengo.espie.dev/calendar?locale=" + location.String() + "&size=200" + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10)

	title := "🍔 Today's meal!"
	content := getMenu(crousRestaurantId)

	emptyString := ""

	fields := []discordwebhook.Field{
		{
			&title,
			&emptyString,
			nil,
		},
	}

	if content == nil {
		noMenu := "No menu available for today!"
		content = []discordwebhook.Field{
			{
				&noMenu,
				&emptyString,
				nil,
			},
		}
	}

	fields = append(fields, content...)

	footerText := "Made with ❤️ by @luckmk1 featuring HackTheCrous! | " + time.Now().In(location).Format("2006-01-02 15h04:05 Z0700 MST")

	footer := discordwebhook.Footer{
		Text:    &footerText,
		IconUrl: &thumbnail,
	}

	embed := discordwebhook.Embed{
		Fields:    &fields,
		Thumbnail: nil,
		Footer:    &footer,
	}

	return embed
}

func getMenu(restaurantId int) []discordwebhook.Field {
	url := "https://api.hackthecrous.com/restaurants/meals/" + strconv.Itoa(restaurantId)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching menu:", err)
	}
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal("Error closing body:", err)
			}
		}(resp.Body)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var menu crousMenu
	unmarshalErr := json.Unmarshal(body, &menu)
	if unmarshalErr != nil {
		log.Fatal("Error unmarshall menu: ", unmarshalErr)
	}

	var curatedMenu []foodie

	for _, restaurant := range menu {
		for _, food := range restaurant.Foodies {
			flag := false
			for _, content := range food.Content {
				if content == "menu non communiqué" {
					flag = true
					break
				}
			}
			if !flag {
				curatedMenu = append(curatedMenu, food)
			}
		}
	}

	var fields []discordwebhook.Field

	inline := true

	for _, foodie := range curatedMenu {
		title := "🍽️ " + foodie.Type
		content := ""
		for _, food := range foodie.Content {
			content += food + "\n"
		}
		fields = append(fields, discordwebhook.Field{
			Name:   &title,
			Value:  &content,
			Inline: &inline,
		})
	}

	return fields
}

type crousMenu []struct {
	ID      int
	Type    string
	Day     string
	Foodies []foodie
}

type foodie struct {
	Content []string `json:"content"`
	Type    string   `json:"type"`
}
