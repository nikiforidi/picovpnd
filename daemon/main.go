package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/anatolio-deb/picovpnd/common"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(os.Getenv("AUTOCERT_DOMAIN")), // Use your email or domain here
		Cache:      autocert.DirCache(os.Getenv("AUTOCERT_DIR")),         // Directory to cache certificates
	}
	config := &tls.Config{GetCertificate: certManager.GetCertificate, MinVersion: tls.VersionTLS12}
	l, err := tls.Listen("tcp", ":0", config)
	if err != nil {
		logrus.Fatal(err)
	}
	defer l.Close()
	go postReport(
		l.Addr().(*net.TCPAddr).Port,
		readCertFromAutocertDir(os.Getenv("AUTOCERT_DIR"),
			os.Getenv("AUTOCERT_DOMAIN")))

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("accepted connection from %s\n", conn.RemoteAddr())

		go handler(conn)
	}

}

func handler(connection net.Conn) {
	defer connection.Close()
	var err error
	resp := common.Response{}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		resp.Code++
		resp.Error = err.Error()
	} else {
		req := common.Request{}
		err = json.Unmarshal(buffer[:mLen], &req)
		if err != nil {
			resp.Code++
			resp.Error = err.Error()
		} else {
			err = common.PayloadDispatcher(req)
			if err != nil {
				resp.Code++
				resp.Error = err.Error()
			}
		}
	}

	b, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
	}
	_, err = connection.Write(b)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debug("leaving handler")
}

func postReport(port int, cert string) {
	resp := &http.Response{}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
	report := Report{
		PublicIP:   getPublicIP(),
		DaemonPort: port,
		Cert:       cert,
	}

	b, err := json.Marshal(report)
	if err != nil {
		logrus.Error("Failed to marshal report:", err)
	}
	for resp.StatusCode != http.StatusOK {
		resp, err = client.Post(os.Getenv("REST_URL"), "application/json", bytes.NewBuffer(b))
		if err != nil {
			logrus.Error("Failed to send report:", err)
			time.Sleep(10 * time.Second) // Retry after 10 seconds
		} else {
			logrus.Infof("Report sent, status code: %d", resp.StatusCode)
		}
	}
}

func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		logrus.Error("Failed to get public IP:", err)
		return ""
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.Error("Failed to decode public IP response:", err)
		return ""
	}

	return result["ip"]
}

// Reads the certificate PEM from the autocert directory for the given domain.
// Returns an empty string if not found or on error.
func readCertFromAutocertDir(dir, domain string) string {
	if dir == "" || domain == "" {
		return ""
	}
	certPath := dir + "/" + domain
	b, err := os.ReadFile(certPath)
	if err != nil {
		logrus.Warnf("Could not read cert from autocert dir: %v", err)
		return ""
	}
	return string(b)
}
