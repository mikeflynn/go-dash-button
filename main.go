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
	DashMacs["f0:27:2d:6d:aa:de"] = ToggleMovieLights
	DashMacs["74:75:48:68:5f:8c"] = ToggleComputerLights

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

func NoAction() {
	log.Println("No action on click.")
}

func ToggleComputerLights() {
	config := map[string]interface{}{
		"front_on":  false,
		"front_bri": 200,
		"back_on":   true,
		"back_bri":  127,
	}

	ToggleWorkshopConfig(config)
}

func ToggleMovieLights() {
	config := map[string]interface{}{
		"front_on":  true,
		"front_bri": 100,
		"back_on":   false,
		"back_bri":  200,
	}

	ToggleWorkshopConfig(config)
}

func ToggleWorkshopConfig(config map[string]interface{}) {
	HueSetup(Config.Section("hue").Key("baseStationIP").String(), Config.Section("hue").Key("baseStationUser").String())

	lights, err := HueGetList()
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return
	}

	inConfig := true

	// Check to see if they are in the configuration already
	for _, v := range lights {
		if strings.Contains(v.Name, "Workshop Front") {
			if v.State.On != config["front_on"].(bool) {
				inConfig = false
				break
			}
		} else if strings.Contains(v.Name, "Workshop Back") {
			if v.State.On != config["back_on"].(bool) || v.State.Bri != config["back_bri"].(int) {
				inConfig = false
				break
			}
		}
	}

	for k, v := range lights {
		if inConfig == false {
			// Turn the front lights off and the back lights down to half brightness.
			if strings.Contains(v.Name, "Workshop Front") {
				go func(idx string) {
					for i := 0; i < 3; i++ {
						err := HueSetLight(idx, HueLightState{On: config["front_on"].(bool), Bri: config["front_bri"].(int)})
						if err != nil {
							log.Printf("HUE ERROR: %v", err.Error())
						} else {
							break
						}
					}
				}(k)
			} else if strings.Contains(v.Name, "Workshop Back") {
				go func(idx string) {
					for i := 0; i < 3; i++ {
						err := HueSetLight(idx, HueLightState{On: config["back_on"].(bool), Bri: config["back_bri"].(int)})
						if err != nil {
							log.Printf("HUE ERROR: %v", err.Error())
						} else {
							break
						}
					}
				}(k)
			}
		} else {
			// Turn the lights back on with regular brightness.
			go func(idx string) {
				for i := 0; i < 3; i++ {
					err := HueSetLight(idx, HueLightState{On: true, Bri: 200})
					if err != nil {
						log.Printf("HUE ERROR: %v", err.Error())
					} else {
						break
					}
				}
			}(k)
		}
	}
}

func ToggleWorkshopLights() {
	HueSetup(Config.Section("hue").Key("baseStationIP").String(), Config.Section("hue").Key("baseStationUser").String())

	lights, err := HueGetList()
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return
	}

	toggle := true

	for _, v := range lights {
		if strings.Contains(v.Name, "Workshop") {
			if v.State.On {
				toggle = false
				break
			}
		}
	}

	for k, v := range lights {
		if strings.Contains(v.Name, "Workshop") {
			go func(idx string) {
				for i := 0; i < 3; i++ {
					err := HueSetLight(idx, HueLightState{On: toggle, Bri: 200})
					if err != nil {
						log.Printf("HUE ERROR: %v", err.Error())
					} else {
						break
					}
				}
			}(k)
		}
	}
}
