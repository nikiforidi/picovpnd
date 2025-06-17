package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Daemon struct {
	Address     string `json:"address"`
	Port        string `json:"port"`
	Certificate string `json:"certificate"`
}

func RegisterSelf(daemon Daemon) {
	// This function would typically register the daemon with a service registry
	// or perform some initialization logic.
	// For now, we will just return nil to indicate success.
	b, err := json.Marshal(daemon)
	if err != nil {
		panic("failed to marshal daemon: " + err.Error())
	}

	caCert, err := os.ReadFile("/etc/ssl/certs/PicoVPNAPI.pem")
	if err != nil {
		panic("failed to read CA certificate: " + err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	resp, err := client.Post("https://picovpn.ru/api/daemon", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("failed to register daemon:", err)
		time.Sleep(5 * time.Second) // Retry after 5 seconds
	} else if resp.StatusCode != http.StatusOK {
		panic("failed to register daemon: " + resp.Status)
	}
	defer resp.Body.Close()
}
