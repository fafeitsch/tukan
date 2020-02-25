package http

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type PhoneClient struct {
	client  *http.Client
	port    int
	tokener tokener
	Logger  *log.Logger
}

func BuildPhoneClient(port int, login string, password string) PhoneClient {
	tokener := tokenerImpl{port: port, login: login, password: password}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return PhoneClient{port: port, client: client, tokener: &tokener, Logger: logger}
}

func (p *PhoneClient) Scan(ip string, number int) map[string]string {
	forEach := func(ip string, token string) string {
		p.log("%s is reachable and login is possible", ip)
		return "phone is reachable, login worked"
	}
	result := p.forEachPhoneIn(ip, number, forEach)
	return result
}

func (p *PhoneClient) UploadPhoneBook(ip string, number int, payload string, delimiter string) map[string]string {
	todo := func(ip string, token string) string {
		url := fmt.Sprintf("http://%s:%d/LocalPhonebook", ip, p.port)
		req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
		req.Header.Add("Authorization", "Bearer "+token)
		multipartHeader := fmt.Sprintf("multipart/form-data; boundary=%s", delimiter)
		req.Header.Add("Content-Type", multipartHeader)
		p.log("starting upload of phone book to %s", ip)
		resp, err := p.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
		err = checkResponse(resp, err)
		if err != nil {
			p.log("could not upload phone book to %s: %s", ip, err)
			return "uploading phone book failed"
		}
		p.log("uploaded phone book successfully to %s", ip)
		return "uploading phone book successful"
	}
	return p.forEachPhoneIn(ip, number, todo)
}

func (p *PhoneClient) log(msg string, args ...interface{}) {
	if p.Logger != nil {
		p.Logger.Printf(msg, args...)
	}
}

func (p *PhoneClient) forEachPhoneIn(ip string, number int, todo func(string, string) string) map[string]string {
	currentIp := net.ParseIP(ip)
	result := make(map[string]string)
	for i := 0; i < number; i++ {
		func() {
			p.log("fetching token for %s…", currentIp.String())
			token, err := p.tokener.fetchToken(currentIp.String())
			if err != nil {
				p.log("fetching token for %s failed: %v", currentIp.String(), err)
				result[currentIp.String()] = "login failed"
				return
			}
			defer func() {
				p.log("logging out of %s…", currentIp.String())
				err := p.tokener.logout(currentIp.String(), *token)
				if err != nil {
					p.log("could not logout from %s", currentIp.String())
					result[currentIp.String()] = "logout failed"
				}
			}()
			msg := todo(currentIp.String(), *token)
			result[currentIp.String()] = msg
		}()
		incrementIP(currentIp)
	}
	return result
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (p *PhoneClient) DownloadPhoneBook(ip string) (map[string]string, string) {
	var result string
	todo := func(ip string, token string) string {
		url := fmt.Sprintf("http://%s:%d/SaveLocalPhonebook", ip, p.port)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		p.log("start phone book download from %s…", ip)
		resp, err := p.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
		err = checkResponse(resp, err)
		if err != nil {
			p.log("could not get phone book for %s: %v", ip, err)
			return "could not get phone book"
		}
		p.log("phone book download from %s successful", ip)
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		result = buf.String()
		return "downloading phone book successful"
	}
	resultMap := p.forEachPhoneIn(ip, 1, todo)
	return resultMap, result
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
