package domain

import (
	"fmt"
	"io/ioutil"
	"sort"
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

func InsertIntoTemplate(payload string, delimiter string) string {
	return fmt.Sprintf(payloadTemplate, delimiter, payload, delimiter)
}

type TukanResult map[string]string

func (t TukanResult) String() string {
	result := "\n==========\n Results\n==========\n"
	keys := make([]string, 0, len(t))
	for ip, _ := range t {
		keys = append(keys, ip)
	}
	sort.Strings(keys)
	for _, ip := range keys {
		result = result + fmt.Sprintf("%s: %s\n", ip, t[ip])
	}
	return result
}
