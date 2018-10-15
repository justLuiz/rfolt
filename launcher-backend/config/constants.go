package config

const CfgServerUrl = "http://127.0.0.1:9411"

var serverId string

func SetServerId(id string) {
	serverId = id
}
