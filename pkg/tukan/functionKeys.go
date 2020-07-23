package tukan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/down"
	"github.com/fafeitsch/Tukan/pkg/tukan/up"
	"net/http"
)

// Downloads the phone's parameters, for example the function key definitions from the
// telephone or returns an error if the download is not successful.
func (p *Phone) DownloadParameters() (*down.Parameters, error) {
	url := fmt.Sprintf("%s/Parameters", p.address)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+p.token)
	resp, err := p.client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	params := down.Parameters{}
	err = json.NewDecoder(resp.Body).Decode(&params)
	if err == nil {
		params.FunctionKeys = purgeTrailingFunctionKeys(params.FunctionKeys)
		return &params, nil
	}
	return nil, err
}

func PrepareDownloadParameters(callback func(result *ParametersResult)) func(p *Phone) {
	return func(p *Phone) {
		params, err := p.DownloadParameters()
		callback(&ParametersResult{Address: p.address, Parameters: params, PhoneResult: PhoneResult{Error: err}})
	}
}

func purgeTrailingFunctionKeys(keys down.FunctionKeys) down.FunctionKeys {
	index := len(keys) - 1
	for index >= 0 && keys[index].IsEmpty() {
		index = index - 1
	}
	current := 0
	result := make([]down.FunctionKey, 0, 8)
	for current <= index {
		result = append(result, keys[current])
		current = current + 1
	}
	return result
}

// Uploads the parameters to the telephone. Returns an error if
// an error occurred during the request or if the response code was not successful.
func (p *Phone) UploadParameters(params up.Parameters) error {
	url := fmt.Sprintf("%s/Parameters", p.address)
	payload, _ := json.Marshal(params)
	reader := bytes.NewBuffer(payload)
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+p.token)
	resp, err := p.client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return checkResponse(resp, err)
}

type ParametersResult struct {
	PhoneResult
	Address    string
	Parameters *down.Parameters
}

func PrepareUploadParameters(params up.Parameters, callback ResultCallback) func(p *Phone) {
	return func(p *Phone) {
		err := p.UploadParameters(params)
		callback(&PhoneResult{Address: p.address, Error: err})
	}
}
