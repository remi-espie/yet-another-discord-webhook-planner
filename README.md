# 📅 Yet Another Discord WebHook Planner

## What is this ❓

This is a simple Discord WebHook Planner 📅, which allows you to fetch your events from a .ics online and display them in Discord through a webhook.

It displays 2 categories of events:
- courses, which are displayed as a list of events for the current day. It also displays the hour of the first and last event of this category.
- important, which are important events that are displayed as a list of events for the current day.

Moreover, it displays the events of today and tomorrow, but only if there are events for these days.  
Besides, it displays a weather forecast for the day.  
Else, it doesn't send a message.

It uses a `.env` file which contains your personal data, see [`.env.example`](.env.example):

## 🚀 Getting started

### Get a correct .ics !

This little script is made to work with specifics .ics file; it assumes that your relevant events: 
- are in the `VEVENT` section of the file, and that the `SUMMARY` field contains the name of the event.
- have the property `CATEGORIES` set to `Cours` if it is a course, or `Important` if it is an important event.

I personally used [Calendar.online](https://calendar.online/) as it is free, allow to fetch multiple data sources and can be edited by anyone from a link without login.

### OpenMeteo

The script also fetch the weather for the day and the following day from [OpenMeteo](https://open-meteo.com). As OpenMeteo is API-key free, it is not necessary to get one.

### Get Go or see [releases](https://github.com/remi-espie/yet-another-discord-webhook-planner/releases)

First import the dependencies with `go get -d ./...`.

Then update the `.env` file with your own data:
```
WEBHOOK_URL= # Your webhook URL
ICS_URL= # Your .ics URL
AVATAR_URL= # Your avatar URL
LONGITUDE= # The longitude of where you'd like to know the weather
LATITUDE= # The latitude of where you'd like to know the weather
```

Finally just run `go run main.go` and it will fetch the set ics and sent the data to your webhook.

You could also build it with `go build main.go` and run it with `./main`.

#### The only release available was built with and for Linux x86_64.

## 📝 License

This project is under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💻 Dependencies for nerds

Developed with GO 19

And the following dependencies:
- [goland-ical](https://github.com/arran4/golang-ical)
- [Discord Webhook](https://github.com/gtuk/discordwebhook)
- [GoDotEnv](https://github.com/joho/godotenv)
- [Open-Meteo-Go](https://github.com/hectormalot/omgo)
