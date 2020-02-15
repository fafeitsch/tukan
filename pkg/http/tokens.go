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

func (p *PhoneClient) Scan(cidr string) error {
	ipaddress, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("could not parse CIDR \"%s\": %v", cidr, err)
	}
	ips := make([]string, 0, 0)
	for ip := ipaddress.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	for _, ip := range ips {
		log.Printf("Checking %s …", ip)
		token, err := p.fetchToken(ip)
		if err != nil {
			log.Printf("Error getting token: %v", err)
		} else {
			log.Printf("Token obtained …")
		}
		if token != nil {
			err = p.logout(ip, *token)
			if err != nil {
				log.Printf("Logout failed: %v", err)
			} else {
				log.Printf("Logout successful")
			}
		}
	}
	return nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (p *PhoneClient) UploadPhoneBook(ip string, number int, payload string, delimiter string) {
	currentIp := net.ParseIP(ip)
	for i := 0; i < number; i++ {
		log.Printf("Try to fetch token for IP %s …", currentIp)
		token, err := p.fetchToken(currentIp.String())
		if err != nil {
			log.Printf("Fetching token for %s failed: %v", currentIp, err)
			continue
		}
		defer p.logout(ip, *token)
		url := fmt.Sprintf("http://%s:%d/LocalPhonebook", currentIp, p.Port)
		req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
		req.Header.Add("Authorization", "Bearer "+*token)
		multipartHeader := fmt.Sprintf("multipart/form-data; boundary=%s", delimiter)
		req.Header.Add("Content-Type", multipartHeader)
		resp, err := p.Client.Do(req)
		err = checkResponse(resp, err)
		if err != nil {
			log.Printf("could not upload phonebook to %s: %v", ip, err)
		}
		incrementIP(currentIp)
	}
}

func (p *PhoneClient) fetchToken(ip string) (*string, error) {
	url := fmt.Sprintf("http://%s:%d/Login", ip, p.Port)
	credentials := api.Credentials{
		Login:    p.Login,
		Password: p.Password,
	}
	payload, _ := json.Marshal(credentials)
	reader := bytes.NewBuffer(payload)
	resp, err := p.Client.Post(url, "application/json", reader)
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tokenResp := api.TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal received token: %v", err)
	}
	return &tokenResp.Token, nil
}

func (p *PhoneClient) logout(ip string, token string) error {
	url := fmt.Sprintf("http://%s:%d/Logout", ip, p.Port)
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+token)
	resp, err := p.Client.Do(request)
	return checkResponse(resp, err)
}

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication error, status code: %d, message: \"%s\"", resp.StatusCode, resp.Status)
	}
	if resp.StatusCode >= 399 {
		return fmt.Errorf("unexpected status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	return nil
}
