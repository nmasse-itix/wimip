package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/chrismarget/arp"
	"github.com/spf13/viper"
)

type IPInfo struct {
	IP         string `json:"ip"`
	MacAddress string `json:"mac,omitempty"`
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	viper.SetDefault("ListenAddr", ":8080")
	viper.SetEnvPrefix("WIMIP")
	viper.AutomaticEnv()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var info IPInfo
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(500)
		}
		info.IP = ip

		arp.CacheUpdate()
		info.MacAddress = arp.Search(ip)

		out, _ := json.MarshalIndent(info, "", " ")
		w.Write(out)
	})

	log.Printf("Listening on %s...", viper.GetString("ListenAddr"))
	log.Fatal(http.ListenAndServe(viper.GetString("ListenAddr"), nil))
}
