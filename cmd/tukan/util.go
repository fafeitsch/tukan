package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/tukan"
	"github.com/urfave/cli"
	"sort"
	"strings"
	"sync"
)

type commentedResult struct {
	*tukan.PhoneResult
	comment string
}

type action int

const (
	actionLogin action = iota
	actionLogout
	actionUploadPhoneBook
	actionDownloadPhoneBook
	actionReplaceFunctionKeys
	actionDownloadParameters
	actionUploadParameters
)

func (a action) String() string {
	names := []string{"Login", "Logout", "Uploading Phone Book", "Downloading Phone Book", "Replacing Function Keys", "Downloading Parameters", "Uploading Parameters"}
	return names[a]
}

func (a action) handler(consumer chan<- commentedResult) func(*tukan.PhoneResult) {
	return func(result *tukan.PhoneResult) {
		var comment string
		if result.Error == nil {
			comment = fmt.Sprintf("%s successful", a.String())
		} else {
			comment = fmt.Sprintf("%s returned error: %v", a.String(), result.Error)
		}
		commentedResult := commentedResult{PhoneResult: result, comment: comment}
		consumer <- commentedResult
	}
}

func handleResults(wg *sync.WaitGroup, channel chan commentedResult, context *cli.Context) {
	defer wg.Done()
	verbose := context.GlobalBool(verboseFlagName)
	results := make(map[string][]string)
	keys := make([]string, 0, 0)
	for result := range channel {
		if verbose {
			_, _ = fmt.Fprintf(context.App.Writer, "%s: %s\n", result.Address, result.comment)
		}
		if _, present := results[result.Address]; !present {
			results[result.Address] = make([]string, 0, 0)
			keys = append(keys, result.Address)
		}
		results[result.Address] = append(results[result.Address], result.comment)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) <= 0
	})
	if verbose {
		_, _ = fmt.Fprint(context.App.Writer, "\n")
	}
	for _, key := range keys {
		_, _ = fmt.Fprintf(context.App.Writer, "%s:\n\t%s\n", key, strings.Join(results[key], "\n\t"))
	}
}
