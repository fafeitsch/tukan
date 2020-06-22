package tukan

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

const payloadTemplate = `--%s
Content-Disposition: form-data; name="file"; filename="LocalPhonebook.xml"
Content-Type: text/xml

%s

--%s--`

// Uploads the payload as phone book to the telephone. Note, that the IP620/630 treat local
// phone books in XML format. However, this method does not check the XML format of the payload, it
// rather uploads it and leaves the parsing to the telephone.
// If an error occurs, or the response does not carry a successful status, an non-nil error is returned.
func (p *Phone) UploadPhoneBook(payload string) error {
	url := fmt.Sprintf("%s/LocalPhonebook", p.address)
	var delimiter string
	for ok := true; ok; ok = len(delimiter) == 0 || strings.Contains(payload, delimiter) {
		randomBytes := make([]byte, 16)
		rand.Read(randomBytes)
		delimiter = hex.EncodeToString(randomBytes)
	}
	multipartFormData := fmt.Sprintf(payloadTemplate, delimiter, payload, delimiter)
	req, _ := http.NewRequest("POST", url, strings.NewReader(multipartFormData))
	req.Header.Add("Authorization", "Bearer "+p.token)
	multipartHeader := fmt.Sprintf("multipart/form-data; boundary=%s", delimiter)
	req.Header.Add("Content-Type", multipartHeader)
	resp, err := p.client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	err = checkResponse(resp, err)
	return err
}

func (p *Phone) DownloadPhoneBook() (*string, error) {
	url := fmt.Sprintf("%s/SaveLocalPhonebook", p.address)
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
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	result := buf.String()
	return &result, nil
}
