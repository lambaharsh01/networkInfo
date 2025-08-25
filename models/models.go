package models

type LocalHostInfo struct {
	HostName string   `json:"hostName"`
	OS       string   `json:"os"`
	UserName string   `json:"userName"`
	IPs      []string `json:"ips"`
}

type NetworkDevicesInfo struct {
	IP       string `json:"ip"`
	HostName string `json:"hostName"`
}

type JSONReport struct {
	LocalHost      LocalHostInfo        `json:"localHost"`
	NetworkDevices []NetworkDevicesInfo `json:"networkDevices"`
}