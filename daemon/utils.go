package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type IP struct {
    Query string
}

func GetPublicIP() string {
    req, err := http.Get("http://ip-api.com/json/")
    if err != nil {
        return err.Error()
    }
    defer req.Body.Close()

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return err.Error()
    }

    var ip IP
    json.Unmarshal(body, &ip)

    return ip.Query
}


func GetLocalHostname() string{
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}