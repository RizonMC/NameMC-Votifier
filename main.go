package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Lukaesebrot/mojango"
	"github.com/NuVotifier/go-votifier"
)

type Votes []string

func main() {
	cfg, err := ReadConfig("./config.json")
	if err != nil {
		log.Printf("error reading config: %s", err)
		return
	}

	req, err := http.NewRequest("GET", "https://api.namemc.com/server/"+cfg.NameMCAddress+"/likes", nil)
	if err != nil {
		log.Printf("error creating request: %s\n", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error doing request: %s", err)
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading body from response: %s", err)
		return
	}

	var votes Votes
	if err := json.Unmarshal(b, &votes); err != nil {
		log.Printf("error unmarshalling: %s\n", err)
	}

	sendVotes(votes, cfg)
}

func sendVotes(votes Votes, cfg *Config) {
	client := mojango.New()

	voter := votifier.NewV2Client(fmt.Sprintf("%s:%d", cfg.Votifier.Address, cfg.Votifier.Port), cfg.Votifier.Token)
	for _, rawUUID := range votes {
		uuid := strings.Replace(rawUUID, "-", "", -1)
		profile, err := client.FetchProfile(uuid, true)
		if err != nil {
			log.Printf("error fetching %s profile: %s\n", uuid, err)
			return
		}
		log.Printf("fetching %q -> %q\n", uuid, profile.Name)
		if err := voter.SendVote(votifier.NewVote("testing-namemc", profile.Name, "127.0.0.1")); err != nil {
			log.Printf("error sending vote from %s with username %s: %s", uuid, profile.Name, err)
			return
		}
	}
}
