package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var HueBaseStationIP string
var HueUserName string

type HueLightState struct {
	Alert     string `json:"alert,omitempty"`
	Bri       int    `json:"bri,omitempty"`
	On        bool   `json:"on"`
	Reachable bool   `json:"reachable,omitempty"`
}

type HueLight struct {
	State            HueLightState `json:"state"`
	Type             string        `json:"type"`
	Name             string        `json:"name"`
	Modelid          string        `json:"modelid"`
	Manufacturername string        `json:"manufacturername"`
	Uniqueid         string        `json:"uniqueid"`
	Swversion        string        `json:"swversion"`
	Pointsymbol      struct {
		One   string `json:"1"`
		Two   string `json:"2"`
		Three string `json:"3"`
		Four  string `json:"4"`
		Five  string `json:"5"`
		Six   string `json:"6"`
		Seven string `json:"7"`
		Eight string `json:"8"`
	} `json:"pointsymbol"`
}

type HueLightList struct {
	One   HueLight `json:"1,omitempty"`
	Two   HueLight `json:"2,omitempty"`
	Three HueLight `json:"3,omitempty"`
	Four  HueLight `json:"4,omitempty"`
	Five  HueLight `json:"5,omitempty"`
	Six   HueLight `json:"6,omitempty"`
	Seven HueLight `json:"7,omitempty"`
	Eight HueLight `json:"8,omitempty"`
	Nine  HueLight `json:"9,omitempty"`
	Ten   HueLight `json:"10,omitempty"`
}

func HueSetup(baseStationIP string, userName string) {
	HueBaseStationIP = baseStationIP
	HueUserName = userName
}

func HueGetList() (map[string]HueLight, error) {
	response, err := http.Get("http://" + HueBaseStationIP + "/api/" + HueUserName + "/lights")
	if err != nil {
		return map[string]HueLight{}, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return map[string]HueLight{}, err
		}

		// Because of Hue's weird API response we're going to unmarshal
		// to a map rather than a struct

		var lightTempMap map[string]*json.RawMessage
		err = json.Unmarshal(contents, &lightTempMap)
		if err != nil {
			return map[string]HueLight{}, err
		}

		lightMap := make(map[string]HueLight)
		for k, _ := range lightTempMap {
			var m HueLight
			_ = json.Unmarshal(*lightTempMap[k], &m)

			lightMap[k] = m
		}

		return lightMap, nil
	}
}

func HueGetLight(id string) (HueLight, error) {
	response, err := http.Get("http://" + HueBaseStationIP + "/api/" + HueUserName + "/lights/" + id)
	if err != nil {
		return HueLight{}, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return HueLight{}, err
		}

		var light HueLight
		err = json.Unmarshal(contents, &light)
		if err != nil {
			return HueLight{}, err
		}

		return light, nil
	}
}

func HueSetLight(id string, options HueLightState) error {
	url := "http://" + HueBaseStationIP + "/api/" + HueUserName + "/lights/" + id + "/state"

	jsonStr, err := json.Marshal(options)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(contents), "error") {
		return errors.New(string(contents))
	}

	return nil
}
