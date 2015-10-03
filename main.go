package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

var Config *ini.File

func main() {
	// Define buttons and action function
	DashMacs["74:75:48:10:33:ec"] = ToggleWorkshopLights

	// Load config file
	var conf = flag.String("conf", "./go-dash-button.ini", "Configuration file required for button events.")
	flag.Parse()

	if *conf != "" {
		if _, err := os.Stat(*conf); os.IsNotExist(err) {
			log.Printf("Can't find config file at: %s", *conf)
			os.Exit(0)
		}

		var err error
		Config, err = ini.Load(*conf)
		if err != nil {
			log.Printf("Unable to parse config file.")
			os.Exit(0)
		}
	}

	// Kick it off!
	SnifferStart()
}

func ToggleWorkshopLights() {
	HueSetup(Config.Section("hue").Key("baseStationIP").String(), Config.Section("hue").Key("baseStationUser").String())

	lights, err := HueGetList()
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return
	}

	for k, v := range lights {
		if strings.Contains(v.Name, "Workshop") {
			toggle := true
			if v.State.On {
				toggle = false
			}

			go HueSetLight(k, HueLightState{On: toggle, Bri: 200})
		}
	}
}
