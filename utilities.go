package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func getIP() string {
	var (
		out      StationJSON
		replacer = strings.NewReplacer("[", "{", "]", "}", "{", "", "}", "")
	)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.Endpoint, nil)
	// Add Accept-Language header, otherwise the modem will throw bad requests at us
	req.Header.Set("Accept-Language", "it-IT")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error while requesting ip: " + err.Error())
		return ""
	}

	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	// The JSON is given to us in an array. We parse that and remove the brackets, and add them only at the end
	_ = json.Unmarshal([]byte(replacer.Replace(string(b))), &out)

	return out.WanIP4Addr
}
