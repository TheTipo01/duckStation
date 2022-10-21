package main

import "time"

type Config struct {
	CFToken  string        `fig:"cf_token" validate:"required"`
	LogLevel string        `fig:"loglevel" default:"error"`
	Endpoint string        `fig:"endpoint" validate:"required"`
	ZoneID   string        `fig:"zone_id" validate:"required"`
	Zones    []string      `fig:"zones" validate:"required"`
	Timeout  time.Duration `fig:"timeout" validate:"required"`

	DDDomain string `fig:"dd_domain" validate:"required"`
	DDToken  string `fig:"dd_token" validate:"required"`
}

// StationJSON is the dumb structure used to unmarshall the request from the router
type StationJSON struct {
	WanIP4Addr string `json:"wan_ip4_addr"`
}
