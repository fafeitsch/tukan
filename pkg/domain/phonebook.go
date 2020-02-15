package domain

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const payloadTemplate = `--%s
Content-Disposition: form-data; name="file"; filename="LocalPhonebook.xml"
Content-Type: text/xml

%s

--%s--`

func LoadAndEmbedPhonebook(file string, delimiter string) (*string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read phone book file %s: %v", file, err)
	}
	text := string(content)
	if strings.Contains(text, delimiter) {
		return nil, fmt.Errorf("the phone book contains the delimiter string \"%s\", which is not allowed", delimiter)
	}
	result := fmt.Sprintf(payloadTemplate, delimiter, text, delimiter)
	return &result, nil
}
