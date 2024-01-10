# Meetup-API-Scraper

Generate monthly analytics reports for YTD meetup events held, rsvps per event, number speakers, etc

[![Actions Status](https://github.com/soypete/{}/workflows/build/badge.svg)](https://github.com/soypete/{}/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/soypete/{}/branch/master/graph/badge.svg)](https://codecov.io/gh/soypete/{})

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
