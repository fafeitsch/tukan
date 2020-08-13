package tukan

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func (p *Phone) Backup() ([]byte, error) {
	url := fmt.Sprintf("%s/SaveAllSettings", p.Address)
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
	return ioutil.ReadAll(resp.Body)
}

func (p *Phone) Restore(data []byte) error {
	url := fmt.Sprintf("%s/RestoreSettings", p.Address)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := injectableCreateFormFile(writer)
	if err != nil {
		return nil
	}
	_, err = part.Write(data)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+p.token)
	resp, err := p.client.Do(request)
	return checkResponse(resp, err)
}

func injectableCreateFormFile(w *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="AllSettings.cfg"`))
	h.Set("Content-Type", "application/octet-stream")
	return w.CreatePart(h)
}
