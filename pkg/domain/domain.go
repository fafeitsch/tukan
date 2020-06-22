package domain

import (
	"fmt"
	"sort"
)

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
