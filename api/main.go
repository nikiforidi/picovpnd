package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Daemon struct {
	Address     string `json:"address"`
	Port        string `json:"port"`
	Certificate string `json:"certificate"`
}

func RegisterSelf(daemon Daemon) {
	b, err := json.Marshal(daemon)
	if err != nil {
		log.Println("failed to marshal daemon:", err)
		return
	}

	req, err := http.NewRequest("POST", "https://picovpn.ru/api/daemon", bytes.NewBuffer(b))
	if err != nil {
		log.Println("failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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
