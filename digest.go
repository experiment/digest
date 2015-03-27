package main

import (
	"log"
	"regexp"
	"time"
)

type digest struct {
	config *config
	client *client
}

// DigestNoteRegex is the regex used to find notes to add to the digest
var DigestNoteRegex = regexp.MustCompile(`COMM:(.*)(\n|$)`)

func (d *digest) get(since time.Time) {
	PRs := d.client.getPullRequestsMergedSince(since)

	for _, PR := range PRs {
		if PR.Body != nil {
			if match := DigestNoteRegex.FindSubmatch([]byte(*PR.Body)); match != nil {
				// Log match
				log.Print(string(match[1]))
			}
		}
	}
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

	digest.get(since)
}
