package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Daemon struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	Certificate []byte `json:"certificate"`
}

func RegisterSelf(daemon Daemon) {
	b, err := json.Marshal(daemon)
	if err != nil {
		log.Println("failed to marshal daemon:", err)
		return
	}

	resp, err := http.Post("https://picovpn.ru/api/daemons", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("failed to send request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to register daemon, status code: %d\n", resp.StatusCode)
		return
	}
}
