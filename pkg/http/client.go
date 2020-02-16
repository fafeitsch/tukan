package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api"
	"log"
	"net"
	"net/http"
	"strings"
)

type PhoneClient struct {
	Client   *http.Client
	Port     int
	Login    string
	Password string
}

func (p *PhoneClient) Scan(ip string, number int) error {
	forEach := func(ip string, token string) {
		p.logout(ip, token)
	}
	p.forEachPhoneIn(ip, number, forEach)
	return nil
}

func (p *PhoneClient) UploadPhoneBook(ip string, number int, payload string, delimiter string) {
	todo := func(ip string, token string) {
		defer p.logout(ip, token)
		url := fmt.Sprintf("http://%s:%d/LocalPhonebook", ip, p.Port)
		req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
		req.Header.Add("Authorization", "Bearer "+token)
		multipartHeader := fmt.Sprintf("multipart/form-data; boundary=%s", delimiter)
		req.Header.Add("Content-Type", multipartHeader)
		resp, err := p.Client.Do(req)
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
		token := p.fetchToken(currentIp.String())
		if token == nil {
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
	token := p.fetchToken(ip)
	if token == nil {
		log.Printf("fetching token for %s failed", ip)
		return ""
	}
	defer p.logout(ip, *token)
	url := fmt.Sprintf("http://%s:%d/SaveLocalPhonebook", ip, p.Port)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+*token)
	resp, err := p.Client.Do(req)
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

func (p *PhoneClient) fetchToken(ip string) *string {
	url := fmt.Sprintf("http://%s:%d/Login", ip, p.Port)
	log.Printf("fetching token for %s (%s) …", ip, url)
	credentials := api.Credentials{
		Login:    p.Login,
		Password: p.Password,
	}
	payload, _ := json.Marshal(credentials)
	reader := bytes.NewBuffer(payload)
	resp, err := p.Client.Post(url, "application/json", reader)
	err = checkResponse(resp, err)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	defer resp.Body.Close()
	tokenResp := api.TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		log.Printf("could not unmarshal token from %s: %v", ip, err)
		return nil
	}
	log.Printf("Fetching token for %s successful", ip)
	return &tokenResp.Token
}

func (p *PhoneClient) logout(ip string, token string) {
	url := fmt.Sprintf("http://%s:%d/Logout", ip, p.Port)
	log.Printf("logging out of %s (%s) …", ip, url)
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+token)
	resp, err := p.Client.Do(request)
	err = checkResponse(resp, err)
	if err != nil {
		log.Printf("could not logout from %s: %v", ip, err)
	} else {
		log.Printf("logout of %s successful", ip)
	}
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
