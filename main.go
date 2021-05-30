package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/chrismarget/arp"
	"github.com/spf13/viper"
)

type IPInfo struct {
	IP         string `json:"ip"`
	MacAddress string `json:"mac,omitempty"`
}

type IPInfoHandler struct {

}

func (h *IPInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var info IPInfo

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(500)
	}
	
	fwdAddress := r.Header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		ip = fwdAddress // If it's a single IP, then awesome!

		// If we got an array... grab the first IP
		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ip = ips[0]
		}
	}

	info.IP = ip

	arp.CacheUpdate()
	info.MacAddress = arp.Search(ip)

	out, _ := json.MarshalIndent(info, "", " ")
	w.Write(out)
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	viper.SetDefault("ListenAddr", ":8080")
	viper.SetEnvPrefix("WIMIP")
	viper.AutomaticEnv()

	var handler *IPInfoHandler
	log.Printf("Listening on %s...", viper.GetString("ListenAddr"))
	log.Fatal(http.ListenAndServe(viper.GetString("ListenAddr"), handler))
}
