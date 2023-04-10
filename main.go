package main

import (
	"log"

	"github.com/soypete/Metup-Go-Graphql-Scraper/auth"
	"github.com/soypete/Metup-Go-Graphql-Scraper/meetup"
)

func main() {
	bearerToken, err := auth.GetBearerToken()
	if err != nil {
		log.Fatal(err)
	}
	meetupClient := meetup.Setup(bearerToken)

	// TODO(soypete): create csv with relevant data.
	// for list of data in README.md
	err = meetupClient.BuildDataSet()
	if err != nil {
		log.Fatal(err)
	}
}
