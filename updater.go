package main

import (
	"github.com/bwmarrin/lit"
	"log"
	"net/http"
)

// Updates DuckDNS
func updateDuckDNS(ip string) {
	_, err := http.Get("https://www.duckdns.org/update?domains=" + cfg.DDDomain + "&token=" + cfg.DDToken + "&ip=" + ip)
	if err != nil {
		log.Println("Error while updating DuckDNS: " + err.Error())
		errorFlag = true
	}

	wg.Done()
}

// Updates Cloudflare given zones
func updateDNSZones(ip string) {
	for _, r := range records {
		r.Content = ip

		err := api.UpdateDNSRecord(ctx, cfg.ZoneID, r.ID, r)
		if err != nil {
			lit.Error("Error while updating DNS record " + r.Name + ": " + err.Error())
			errorFlag = true
		}
	}

	wg.Done()
}
