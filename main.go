package main

import (
	"encoding/json"
	"github.com/kkyr/fig"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Config holds data from the config.yml
type Config struct {
	Token    string        `fig:"token" validate:"required"`
	Endpoint string        `fig:"endpoint" validate:"required"`
	Domain   string        `fig:"domain" validate:"required"`
	Timeout  time.Duration `fig:"timeout" validate:"required"`
}

// UserLang is the dumb structure used to unmarshall the request from the router
type UserLang struct {
	WanIP4Addr string `json:"wan_ip4_addr"`
}

var (
	cfg      Config
	replacer = strings.NewReplacer("[", "{", "]", "}", "{", "", "}", "")
)

func init() {
	// Loads the config file
	err := fig.Load(&cfg, fig.File("config.yml"))
	if err != nil {
		panic(err)
	}
}

func main() {
	var ip, newIP string

	for {
		newIP = getIP()
		if newIP != ip {
			ip = newIP
			updateDuckDNS(&ip)
		}

		time.Sleep(cfg.Timeout)
	}
}

// Gets the IP from the modem
func getIP() string {
	var out UserLang

	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.Endpoint, nil)
	// Add Accept-Language header, otherwise the modem will throw bad requests at us
	req.Header.Set("Accept-Language", "it-IT")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error while requesting ip: " + err.Error())
		return ""
	}

	b, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	// The JSON is given to us in an array. We parse that and remove the brackets, and add them only at the end
	_ = json.Unmarshal([]byte(replacer.Replace(string(b))), &out)

	return out.WanIP4Addr
}

// Updated the IP with the given one
func updateDuckDNS(ip *string) {
	_, err := http.Get("https://www.duckdns.org/update?domains=" + cfg.Domain + "&token=" + cfg.Token + "&ip=" + *ip)
	if err != nil {
		log.Println("Error while updating Duck DNS: " + err.Error())
		// Doing so forces another updateDuckDNS
		*ip = ""
	}
}
