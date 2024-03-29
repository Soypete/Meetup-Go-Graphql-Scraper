# Meetup-API-Scraper

Generate monthly analytics reports for YTD meetup events held, rsvps per event, number speakers, etc

[![Actions Status](https://github.com/soypete/meetup-go-graphql-scraper/workflows/build/badge.svg)](https://github.com/soypete/meetup-go-graphql-scraper/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/soypete/Meetup-Go-Graphql-Scraper@v0.1.1.svg)](https://pkg.go.dev/github.com/soypete/Meetup-Go-Graphql-Scraper@v0.1.1)

# Install

## Go env

Run this command when you have working go environment.

```bash
go install github.com/soypete/Meetup-Go-Graphql-Scraper@latest
```

## No Go env

Instructions coming soon

# Setup

This repo takes a config file will all the necessary information to connect to meetup using the [kolla SDK](). The config file is in the json format and should look like the following

```
{
  "pro_account": "go",
  "kolla_key": {kolla.secret},
  "connector_id": {kolla.account},
  "consumer_id": {kolla.key}
}
```

The path for the connector key will default to `config.json` but you can pass a custom path wiht the flag `--config-file`.
