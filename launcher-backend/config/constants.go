package config

const CfgServerUrl = "http://178.128.203.194:9411"

var serverId string

func SetServerId(id string) {
	serverId = id
}
