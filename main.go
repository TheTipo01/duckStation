package main

import (
	"context"
	"github.com/bwmarrin/lit"
	"github.com/cloudflare/cloudflare-go"
	"github.com/kkyr/fig"
	"strings"
	"sync"
	"time"
)

var (
	cfg       Config
	zones     map[string]bool
	records   []cloudflare.DNSRecord
	api       *cloudflare.API
	ctx       = context.Background()
	errorFlag bool
	wg        sync.WaitGroup
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

	// Create the file lastip if it doesn't exist
	if !fileExists("lastip") {
		writeIP("")
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

	// Reads the last saved ip
	ip = readIP()

	// Main loop: checks for a new ip change every cfg.Timeout
	for {
		newIP, err = getIP()

		if err == nil && newIP != ip {
			lit.Info("IP changed from " + ip + " to " + newIP)

			wg.Add(2)
			go updateDuckDNS(newIP)
			go updateDNSZones(newIP)
			wg.Wait()

			// If we don't get any errors, we save the new ip
			if !errorFlag {
				ip = newIP
				writeIP(newIP)
			} else {
				errorFlag = false
			}
		}

		time.Sleep(cfg.Timeout)
	}
}
