# duckStation

[![Go Report Card](https://goreportcard.com/badge/github.com/TheTipo01/duckStation)](https://goreportcard.com/report/github.com/TheTipo01/duckStation)

A [Duck DNS](https://www.duckdns.org) and [Cloudflare](https://www.cloudflare.com) updater that gets the public
IP from
the [Vox30](https://openwrt.org/toh/vodafone/vodafone_power_station), also known as Vodafone Power Station in italy
Vodafone Wi-Fi Hub in UK or Vodafone Gigabox in ireland

It works by polling the `/data/user_lang.json` endpoint on the router every `config.Timeout` seconds, and compares
it to the old stored IP. If there's a change, the program updates the IP by calling Duck DNS.

# Usage

Grab a [release](https://github.com/TheTipo01/duckStation/releases) or build the executable, after that modify
the `example_config.yml`, rename it to `config.yml`, and run the program!