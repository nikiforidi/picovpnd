package ip

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.IP == "" {
		return "", fmt.Errorf("no IP found in response")
	}
	return result.IP, nil
}
