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

func connectToPhones(context *cli.Context) tukan.Connections {
	login := context.GlobalString(loginFlagName)
	password := context.GlobalString(passwordFlagName)
	timeout := context.GlobalInt(timeoutFlagName)
	connector := tukan.Connector{Client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, UserName: login, Password: password}
	addresses := tukan.ExpandAddresses("http", context.Args()...)
	return connector.MultipleConnect(addresses...)
}

func scan(context *cli.Context) {
	channel := connectToPhones(context).Scan()
	handleSimpleResults(channel, context)
}

func handleSimpleResults(channel chan tukan.SimpleResult, context *cli.Context) {
	verbose := context.GlobalBool(verboseFlagName)
	results := make([]tukan.SimpleResult, 0, 0)
	for result := range channel {
		if verbose {
			_, _ = fmt.Fprintf(context.App.Writer, "%s\n", result.String())
		}
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return strings.Compare(results[i].Address, results[j].Address) <= 0
	})
	if verbose {
		_, _ = fmt.Fprint(context.App.Writer, "\n")
	}
	for _, result := range results {
		_, _ = fmt.Fprintf(context.App.Writer, "%s\n", result.String())
	}
}

func uploadPhoneBook(context *cli.Context) {
	file := context.String(fileFlagName)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not load phone book file: %v", err)
		return
	}

	channel := connectToPhones(context).UploadPhoneBook(string(content))
	handleSimpleResults(channel, context)
}

func downloadPhoneBook(context *cli.Context) {
	targetDirectory := context.String(targetDirFlagName)
	err := os.MkdirAll(targetDirectory, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not create target directory: %v", err)
		return
	}
	channel := connectToPhones(context).DownloadPhoneBook()
	simpleResults := make(chan tukan.SimpleResult)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSimpleResults(simpleResults, context)
	}()
	writeErrors := make([]string, 0, 0)
	for phoneBookResult := range channel {
		simpleResults <- phoneBookResult.SimpleResult
		if phoneBookResult.PhoneBook != nil {
			fileName := fileName(phoneBookResult.Address)
			path := filepath.Join(targetDirectory, fileName)
			err := ioutil.WriteFile(path, []byte(*phoneBookResult.PhoneBook), os.ModePerm)
			if err != nil {
				writeErrors = append(writeErrors, err.Error())
			}
		}
	}
	close(simpleResults)
	if len(writeErrors) != 0 {
		_, _ = fmt.Fprintf(context.App.Writer, "There were errors writing the files:\n: %s", strings.Join(writeErrors, "\n"))
	}
	wg.Wait()
}

func fileName(address string) string {
	regex := regexp.MustCompile("https?://")
	result := regex.ReplaceAllString(address, "")
	result = strings.ReplaceAll(result, ":", "_")
	return "phonebook_" + result + ".xml"
}
