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

// Downloads the phone book from the telephone. In case of an error
// the returned string is nil.
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

// Uploads the given payload to the phone book endpoint of every telephone that
// comes into the connection channel. This function returns immediately and reports
// the results of the uploads by means of the onProcess callback.
//
// The function only uploads the string and does not check the string for valid XML format.
// The behaviour of the telephone is undefined when unproper content is uploaded. This method does
// not check the format of the payload.
func (c Connections) UploadPhoneBook(onProcess ResultCallback, payload string) Connections {
	result := make(Connections)
	singleAction := func(phone *Phone) {
		err := phone.UploadPhoneBook(payload)
		if err != nil {
			phoneResult := PhoneResult{Address: phone.address, Comment: err.Error()}
			onProcess(&phoneResult)
		} else {
			onProcess(&PhoneResult{Address: phone.address, Comment: "Upload successful"})
		}
		result <- phone
	}
	end := func() {
		close(result)
	}
	go c.loop(singleAction, end)
	return result
}

type PhoneBookResult struct {
	Address   string
	PhoneBook string
}

// Uses the connection to download the phone books from the telephone. The method
// returns immediately. The onProcess callback is called for every tried phone, independently of
// the download success. The onSuccess callback is called for those phones of which the phone book
// could be downloaded successfully.
func (c Connections) DownloadPhoneBook(onProcess ResultCallback, onSuccess func(result *PhoneBookResult)) Connections {
	result := make(Connections)
	singleAction := func(phone *Phone) {
		book, err := phone.DownloadPhoneBook()
		if err != nil {
			phoneResult := PhoneResult{Address: phone.address, Comment: err.Error(), Error: err}
			onProcess(&phoneResult)
		} else {
			onProcess(&PhoneResult{Address: phone.address, Comment: "Download successful"})
			onSuccess(&PhoneBookResult{Address: phone.address, PhoneBook: *book})
		}
		result <- phone
	}
	end := func() {
		close(result)
	}
	go c.loop(singleAction, end)
	return result
}
