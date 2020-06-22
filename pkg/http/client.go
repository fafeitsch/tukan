package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api/down"
	"github.com/fafeitsch/Tukan/pkg/api/up"
	"github.com/fafeitsch/Tukan/pkg/domain"
	"github.com/fafeitsch/Tukan/pkg/tukan"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type PhoneClient struct {
	client   *http.Client
	port     int
	login    string
	password string
	Logger   *log.Logger
}

func BuildPhoneClient(port int, login string, password string, timeoutSeconds int) PhoneClient {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return PhoneClient{port: port, client: client, login: login, password: password, Logger: logger}
}

func (p *PhoneClient) Scan(ip string, number int) domain.TukanResult {
	forEach := func(phone tukan.Phone) string {
		p.log("%v is reachable and login is possible", phone)
		return "phone is reachable, login worked"
	}
	result := p.forEachPhoneIn(ip, number, forEach)
	return result
}

func (p *PhoneClient) UploadPhoneBook(ip string, number int, payload string) domain.TukanResult {
	todo := func(phone tukan.Phone) string {
		err := phone.UploadPhoneBook(payload)
		if err != nil {
			p.log("could not up phone book to %s: %s", ip, err)
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

func (p *PhoneClient) forEachPhoneIn(ip string, number int, todo func(phone tukan.Phone) string) domain.TukanResult {
	currentIp := net.ParseIP(ip)
	result := make(map[string]string)
	for i := 0; i < number; i++ {
		func() {
			p.log("fetching token for %s…", currentIp.String())
			url := fmt.Sprintf("http://%s:%d", ip, p.port)
			phone, err := tukan.Connect(p.client, url, p.login, p.password)
			if err != nil {
				p.log("fetching token for %s failed: %v", currentIp.String(), err)
				result[currentIp.String()] = "login failed"
				return
			}
			defer func() {
				p.log("logging out of %s…", currentIp.String())
				err := phone.Logout()
				if err != nil {
					p.log("could not logout from %s", currentIp.String())
					result[currentIp.String()] = "logout failed"
				}
			}()
			msg := todo(*phone)
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

func (p *PhoneClient) DownloadPhoneBook(ip string) (domain.TukanResult, string) {
	var result string
	todo := func(phone tukan.Phone) string {
		ptr, err := phone.DownloadPhoneBook()
		if err != nil {
			return "could not get phonebook"
		}
		result = *ptr
		return "Success"
	}
	resultMap := p.forEachPhoneIn(ip, 1, todo)
	return resultMap, result
}

func (p *PhoneClient) DownloadFunctionKeys(ip string, number int) domain.TukanResult {
	todo := func(phone tukan.Phone) string {
		url := fmt.Sprintf("%s/Parameters", phone.Host())
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+phone.Token())
		p.log("start function key download from %s…", ip)
		resp, err := p.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
		err = checkResponse(resp, err)
		if err != nil {
			p.log("could not get function keys from %s: %v", ip, err)
			return "could not get function keys"
		}
		p.log("function keys successfully downloaded from %s", ip)
		params := down.Parameters{}
		err = json.NewDecoder(resp.Body).Decode(&params)
		if err != nil {
			p.log("error deserializing the function keys from %s: %v", ip, err)
			return "could not deserialize function keys"
		}
		params.PurgeTrailingFunctionKeys()
		return params.FunctionKeys.String()
	}
	return p.forEachPhoneIn(ip, number, todo)
}

func (p *PhoneClient) ReplaceFunctionKeyName(ip string, number int, original string, replace string) domain.TukanResult {
	todo := func(phone tukan.Phone) string {
		url := fmt.Sprintf("%s/Parameters", phone.Host())
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+phone.Token())
		p.log("downloading function keys from %s …", phone.Host())
		resp, err := p.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
		err = checkResponse(resp, err)
		if err != nil {
			p.log("could not get function keys from %s: %v", phone.Host(), err)
			return "could not get function keys"
		}
		p.log("function keys successfully downloaded from %s", phone.Host())
		params := down.Parameters{}
		err = json.NewDecoder(resp.Body).Decode(&params)
		if err != nil {
			p.log("error deserializing the function keys from %s: %v", phone.Host(), err)
			return "could not deserialize function keys"
		}
		params.PurgeTrailingFunctionKeys()
		newKeys := p.buildNewFunctionKeys(params, original, replace)
		payload, _ := json.Marshal(&newKeys)
		reader := bytes.NewBuffer(payload)
		req, _ = http.NewRequest("POST", url, reader)
		req.Header.Add("Authorization", "Bearer "+phone.Token())
		p.log("uploading new function keys from %s …", phone.Host())
		resp, err = p.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
		err = checkResponse(resp, err)
		if err != nil {
			p.log("could not upload function keys to %s: %v", ip, err)
			return "could not upload function keys"
		}
		return "function keys updated"
	}
	return p.forEachPhoneIn(ip, number, todo)
}

func (p *PhoneClient) buildNewFunctionKeys(params down.Parameters, original string, replace string) up.Parameters {
	keys := make([]map[string]string, 0, len(params.FunctionKeys))
	for index, fnKey := range params.FunctionKeys {
		var key = map[string]string{}
		if fnKey.DisplayName.Value == original {
			key = map[string]string{"DisplayName": replace}
			p.log("replacing display name \"%s\" of %dth function key with display name \"%s\"", fnKey.DisplayName.Value, index, replace)
		}
		keys = append(keys, key)
	}
	return up.Parameters{FunctionKeys: keys}
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
