package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

func connectToPhones(context *cli.Context, onError tukan.ResultCallback) tukan.Connections {
	login := context.GlobalString(loginFlagName)
	password := context.GlobalString(passwordFlagName)
	timeout := context.GlobalInt(timeoutFlagName)
	connector := tukan.Connector{Client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, UserName: login, Password: password}
	addresses := tukan.ExpandAddresses("http", context.Args()...)
	return connector.MultipleConnect(onError, addresses...)
}

func scan(context *cli.Context) {
	channel := make(chan *tukan.PhoneResult)
	collectResults := func(result *tukan.PhoneResult) {
		channel <- result
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSimpleResults(channel, context)
	}()
	connectToPhones(context, collectResults).Logout(collectResults)
	close(channel)
	wg.Wait()
}

func handleSimpleResults(channel chan *tukan.PhoneResult, context *cli.Context) {
	verbose := context.GlobalBool(verboseFlagName)
	results := make(map[string][]string)
	keys := make([]string, 0, 0)
	for result := range channel {
		if verbose {
			_, _ = fmt.Fprintf(context.App.Writer, "%s\n", result.String())
		}
		if _, present := results[result.Address]; !present {
			results[result.Address] = make([]string, 0, 0)
			keys = append(keys, result.Address)
		}
		results[result.Address] = append(results[result.Address], result.Comment)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) <= 0
	})
	if verbose {
		_, _ = fmt.Fprint(context.App.Writer, "\n")
	}
	for _, key := range keys {
		_, _ = fmt.Fprintf(context.App.Writer, "%s: %s\n", key, strings.Join(results[key], "\n\t"))
	}
}

func uploadPhoneBook(context *cli.Context) {
	file := context.String(fileFlagName)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not load phone book file: %v", err)
		return
	}

	channel := make(chan *tukan.PhoneResult)
	collectResults := func(result *tukan.PhoneResult) {
		channel <- result
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSimpleResults(channel, context)
	}()
	connectToPhones(context, collectResults).
		UploadPhoneBook(collectResults, string(content)).
		Logout(collectResults)
	close(channel)
	wg.Wait()
}

func downloadPhoneBook(context *cli.Context) {
	targetDirectory := context.String(targetDirFlagName)
	err := os.MkdirAll(targetDirectory, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not create target directory: %v", err)
		return
	}
	channel := make(chan *tukan.PhoneResult)
	collectResults := func(result *tukan.PhoneResult) {
		channel <- result
	}

	bookChannel := make(chan *tukan.PhoneBookResult)
	collectBooks := func(result *tukan.PhoneBookResult) {
		bookChannel <- result
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		handleSimpleResults(channel, context)
	}()
	go func() {
		defer wg.Done()
		writeErrors := make([]string, 0, 0)
		for phoneBookResult := range bookChannel {
			fileName := fileName(phoneBookResult.Address)
			path := filepath.Join(targetDirectory, fileName)
			err := ioutil.WriteFile(path, []byte(phoneBookResult.PhoneBook), os.ModePerm)
			if err != nil {
				writeErrors = append(writeErrors, err.Error())
			}
		}
		if len(writeErrors) != 0 {
			_, _ = fmt.Fprintf(context.App.Writer, "There were errors writing the files:\n: %s", strings.Join(writeErrors, "\n"))
		}
	}()
	connectToPhones(context, collectResults).
		DownloadPhoneBook(collectResults, collectBooks).
		Logout(collectResults)
	close(channel)
	close(bookChannel)
	wg.Wait()
}

func fileName(address string) string {
	regex := regexp.MustCompile("https?://")
	result := regex.ReplaceAllString(address, "")
	result = strings.ReplaceAll(result, ":", "_")
	return "phonebook_" + result + ".xml"
}
