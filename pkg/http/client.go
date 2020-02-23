package http

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type PhoneClient struct {
	client  *http.Client
	port    int
	tokener tokener
}

func BuildPhoneClient(port int, login string, password string) PhoneClient {
	tokener := tokenerImpl{port: port, login: login, password: password}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	return PhoneClient{port: port, client: client, tokener: &tokener}
}

func (p *PhoneClient) Scan(ip string, number int) error {
	forEach := func(ip string, token string) {
		p.tokener.logout(ip, token)
	}
	p.forEachPhoneIn(ip, number, forEach)
	return nil
}

func (p *PhoneClient) UploadPhoneBook(ip string, number int, payload string, delimiter string) {
	todo := func(ip string, token string) {
		defer p.tokener.logout(ip, token)
		url := fmt.Sprintf("http://%s:%d/LocalPhonebook", ip, p.port)
		req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
		req.Header.Add("Authorization", "Bearer "+token)
		multipartHeader := fmt.Sprintf("multipart/form-data; boundary=%s", delimiter)
		req.Header.Add("Content-Type", multipartHeader)
		resp, err := p.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
		err = checkResponse(resp, err)
		if err != nil {
			log.Printf("could not upload phonebook to %s: %v", ip, err)
		}
	}
	p.forEachPhoneIn(ip, number, todo)
}

func (p *PhoneClient) forEachPhoneIn(ip string, number int, todo func(string, string)) {
	currentIp := net.ParseIP(ip)
	for i := 0; i < number; i++ {
		log.Printf("Try to fetch token for IP %s …", currentIp)
		token, err := p.tokener.fetchToken(currentIp.String())
		if err != nil {
			log.Printf("fetching token for %s failed; skip %s", currentIp, currentIp)
			continue
		}
		todo(currentIp.String(), *token)
		incrementIP(currentIp)
	}
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (p *PhoneClient) DownloadPhoneBook(ip string) string {
	token, err := p.tokener.fetchToken(ip)
	if err != nil {
		log.Printf("fetching token for %s failed", ip)
		return ""
	}
	defer p.tokener.logout(ip, *token)
	url := fmt.Sprintf("http://%s:%d/SaveLocalPhonebook", ip, p.port)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+*token)
	resp, err := p.client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	err = checkResponse(resp, err)
	if err != nil {
		log.Printf("could not get phonebook from %s: %v", ip, err)
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	return buf.String()
}

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication error, status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	if resp.StatusCode >= 299 {
		return fmt.Errorf("unexpected status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	return nil
}
