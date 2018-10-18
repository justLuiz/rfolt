package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rfolt/configuration"
)

func get(kind string, v interface{}) {
	resp, err := http.Get(fmt.Sprintf("%v/%v/%v", ControlServerUrl, serverId, kind))
	if err != nil {
		v = nil
		return
	}
	defer resp.Body.Close()

	cfg, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(cfg, &v)
	if err != nil {
		v = nil
		return
	}
}

func GetMain() configuration.Main {
	var cfg configuration.Main
	get("main", &cfg)
	return cfg
}

func GetNetwork() configuration.Network {
	var cfg configuration.Network
	get("network", &cfg)
	return cfg
}
