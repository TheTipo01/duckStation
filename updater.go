package main

import (
	"github.com/bwmarrin/lit"
	"log"
	"net/http"
)

// Updated the IP with the given one
func updateDuckDNS(ip string) error {
	_, err := http.Get("https://www.duckdns.org/update?domains=" + cfg.DDDomain + "&token=" + cfg.DDToken + "&ip=" + ip)
	if err != nil {
		log.Println("Error while updating DuckDNS: " + err.Error())
		return err
	}

	return nil
}

func updateDNSZones(ip string) (err error) {
	for _, r := range records {
		r.Content = ip

		err = api.UpdateDNSRecord(ctx, cfg.ZoneID, r.ID, r)
		if err != nil {
			lit.Error("Error while updating DNS record " + r.Name + ": " + err.Error())
		}
	}

	return err
}
