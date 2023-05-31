# Yet Another Discord WebHook Planner

## What is this?

This is a simple Discord WebHook Planner, which allows you to fetch your events from a .ics online and display them in a Discord through a webhook.

It displays 2 categories of events:
- courses, which are displayed as a list of events for the current day. It also displays the hour of the first and last event of this category.
- important, which are important events that are displayed as a list of events for the current day.

Moreover, it displays the events of today and tomorrow, but only if there are events for these days.  
Else, it doesn't send a message.

It uses a `.env` file which contains your personal data, see .env.example:

## How to use it?

### Get a good .ics

This little script is made to work with specifics .ics file; it assumes that your relevant events: 
- are in the `VEVENT` section of the file, and that the `SUMMARY` field contains the name of the event.
- have the property `CATEGORIES` set to `Cours` if it is a course, or `Important` if it is an important event.

### Get Go or see releases

First import the dependencies with `go get -d ./...`.

Then update the `.env` file with your own data:
```
WEBHOOK_URL= # Your webhook URL
ICS_URL= # Your .ics URL
AVATAR_URL= # Your avatar URL
```

Finally just run `go run main.go` and it will fetch the set ics and sent the data to your webhook.

You could also build it with `go build main.go` and run it with `./main`.

#### The only release available was built for Linux x86_64.
