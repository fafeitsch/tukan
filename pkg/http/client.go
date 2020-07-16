package http

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/domain"
	"github.com/fafeitsch/Tukan/pkg/tukan"
	"github.com/fafeitsch/Tukan/pkg/tukan/down"
	"github.com/fafeitsch/Tukan/pkg/tukan/up"
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
			connector := tukan.Connector{Client: p.client, UserName: p.login, Password: p.password}
			phone, err := connector.SingleConnect(url)
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

func (p *PhoneClient) ReplaceFunctionKeyName(ip string, number int, original string, replace string) domain.TukanResult {
	todo := func(phone tukan.Phone) string {
		parameters, err := phone.DownloadParameters()
		if err != nil {
			p.log("could not get function keys from %s: %v", phone.Host(), err)
			return "could not get function keys"
		}
		newKeys := p.buildNewFunctionKeys(*parameters, original, replace)
		err = phone.UploadParameters(newKeys)
		if err != nil {
			p.log("could not upload function keys to %s: %v", ip, err)
			return "could not upload function keys"
		}
		return "function keys replaced successfully"
	}
	return p.forEachPhoneIn(ip, number, todo)
}

func (p *PhoneClient) buildNewFunctionKeys(params down.Parameters, original string, replace string) up.Parameters {
	keys := make([]up.FunctionKey, 0, len(params.FunctionKeys))
	for index, fnKey := range params.FunctionKeys {
		var key = up.FunctionKey{}
		if fnKey.DisplayName.Value == original {
			key = up.FunctionKey{DisplayName: replace}
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
