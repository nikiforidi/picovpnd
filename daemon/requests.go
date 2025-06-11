package main

type Report struct {
	PublicIP   string `json:"publicIP"`
	DaemonPort int    `json:"daemonPort"`
	Cert       string `json:"cert"`
}
