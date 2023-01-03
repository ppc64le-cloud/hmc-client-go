package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ppc64le-cloud/hmc-client-go"
)

func main() {
	client := hmc.NewClient(&hmc.BasicAuthenticator{
		UserName: "hscroot",
		Password: "",
		BaseURL:  "https://192.168.0.100:12443",
	}, hmc.ClientParams{Insecure: true, LogLevel: log.DebugLevel})

	systems, _, err := client.ManagedSystems().GET()
	if err != nil {
		log.Fatalf("errored: %v", err)
	}
	for _, system := range systems.Entry {
		fmt.Printf("ManagedSystem: %+v and entry: %+v\n", system.Content.ManagedSystem, system.Entry)
	}
}
