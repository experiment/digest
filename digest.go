package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type digest struct {
	config *config
	client *client
}

// DigestNoteRegex is the regex used to find notes to add to the digest
var DigestNoteRegex = regexp.MustCompile(`COMM:(.*)(\n|$)`)

func postToSlack(room string, message string) {
	message = url.QueryEscape(message)

	log.Print("Posting to slack: ", message)

	_, err := http.Get("http://gon-bot.herokuapp.com/hubot/say?room=" + room + "&message=" + message)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *digest) get(since time.Time) []string {
	PRs := d.client.getPullRequestsMergedSince(since)

	var notes []string

	for _, PR := range PRs {
		if PR.Body != nil {
			if match := DigestNoteRegex.FindSubmatch([]byte(*PR.Body)); match != nil {
				notes = append(notes, string(match[1]))
			}
		}
	}

	return notes
}

func (d *digest) send(since time.Time) {
	messages := []string{
		"Hey @channel, here's what we've changed in the product in the last week:",
	}
	messages = append(messages, d.get(since)...)

	// Post digest to slack general channel
	postToSlack("general", strings.Join(messages, "\n *"))
}

func main() {
	conf := new(config)
	conf.load()

	digest := &digest{
		config: conf,
		client: &client{config: conf},
	}

	// TODO, get last time from somewhere (redis?)
	// Currently hardcoded to be 7 days ago
	since := time.Now().Add(time.Hour * 24 * 7 * -1)

	digest.send(since)
}
