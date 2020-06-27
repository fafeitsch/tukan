package tukan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api/down"
	"github.com/fafeitsch/Tukan/pkg/api/up"
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
		return nil, fmt.Errorf("response error: %v", err)
	}
	params := down.Parameters{}
	err = json.NewDecoder(resp.Body).Decode(&params)
	if err == nil {
		params.FunctionKeys = purgeTrailingFunctionKeys(params.FunctionKeys)
		return &params, nil
	}
	return nil, err
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
