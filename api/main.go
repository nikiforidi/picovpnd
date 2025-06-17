package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Daemon struct {
	Address     string `json:"address"`
	Port        string `json:"port"`
	Certificate string `json:"certificate"`
}

func RegisterSelf(daemon Daemon) {
	status := 0
	for status != http.StatusOK {
		// This function would typically register the daemon with a service registry
		// or perform some initialization logic.
		// For now, we will just return nil to indicate success.
		b, err := json.Marshal(daemon)
		if err != nil {
			break
		}

		// caCert, err := os.ReadFile("/etc/ssl/certs/PicoVPNAPI.pem")
		// if err != nil {
		// 	break
		// }
		// caCertPool := x509.NewCertPool()
		// caCertPool.AppendCertsFromPEM(caCert)

		// client := &http.Client{
		// 	Transport: &http.Transport{
		// 		TLSClientConfig: &tls.Config{
		// 			RootCAs: caCertPool,
		// 		},
		// 	},
		// }
		// client := http.Client{
		// 	Transport: &http.Transport{
		// 		TLSClientConfig: &tls.Config{
		// 			InsecureSkipVerify: true, // Skip verification for self-signed certs
		// 			// RootCAs: caCertPool, // Uncomment if you have a CA cert pool
		// 		},
		// 	},
		// }

		resp, err := http.Post("https://picovpn.ru/api/daemon", "application/json", bytes.NewBuffer(b))
		if err != nil {
			log.Println("failed to register daemon:", err)
			time.Sleep(5 * time.Second) // Retry after 5 seconds
		} else if resp.StatusCode != http.StatusOK {
			break
		}
	}
	// defer resp.Body.Close()
}
