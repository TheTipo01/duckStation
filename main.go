package main

import (
	"context"
	"github.com/bwmarrin/lit"
	"github.com/cloudflare/cloudflare-go"
	"github.com/kkyr/fig"
	"strings"
	"time"
)

var (
	cfg     Config
	zones   map[string]bool
	records []cloudflare.DNSRecord
	api     *cloudflare.API
	ctx     = context.Background()
)

func init() {
	err := fig.Load(&cfg, fig.File("config.yml"))
	if err != nil {
		lit.Error(err.Error())
		return
	}

	// Set lit.LogLevel to the given value
	switch strings.ToLower(cfg.LogLevel) {
	case "logwarning", "warning":
		lit.LogLevel = lit.LogWarning

	case "loginformational", "informational":
		lit.LogLevel = lit.LogInformational

	case "logdebug", "debug":
		lit.LogLevel = lit.LogDebug
	}

	zones = make(map[string]bool, len(cfg.Zones))
	for _, zone := range cfg.Zones {
		zones[zone] = true
	}
}

func main() {
	var (
		ip, newIP string
		err       error
	)

	api, err = cloudflare.NewWithAPIToken(cfg.CFToken)
	if err != nil {
		panic(err)
	}

	// Save the records we really need to update
	allRecords, _ := api.DNSRecords(ctx, cfg.ZoneID, cloudflare.DNSRecord{})
	for _, record := range allRecords {
		if record.Type == "A" && zones[record.Name] {
			records = append(records, record)
		}
	}

	// Main loop: checks for a new ip change every cfg.Timeout
	for {
		newIP = getIP()
		if newIP != ip {
			ip = newIP
			err = updateDuckDNS(ip)
			err = updateDNSZones(ip)
			// Force an ip update
			if err != nil {
				ip = ""
			}
		}

		time.Sleep(cfg.Timeout)
	}
}
