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

func createConnector(context *cli.Context) *tukan.Connector {
	login := context.GlobalString(loginFlagName)
	password := context.GlobalString(passwordFlagName)
	timeout := context.GlobalInt(timeoutFlagName)
	connector := tukan.Connector{Client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, UserName: login, Password: password}
	addresses := tukan.ExpandAddresses("http", context.Args()...)
	connector.Addresses = addresses
	return &connector
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
	createConnector(context).Run(collectResults, []tukan.PhoneAction{}, collectResults)
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
		// TODO: different result collectors for the different actions with different messages
		channel <- result
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSimpleResults(channel, context)
	}()
	createConnector(context).Run(collectResults,
		[]tukan.PhoneAction{tukan.PreparePhoneBookUpload(collectResults, string(content))},
		collectResults)
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
			fileName := phoneBookFileName(phoneBookResult.Address)
			path := filepath.Join(targetDirectory, fileName)
			err := ioutil.WriteFile(path, []byte(*phoneBookResult.PhoneBook), os.ModePerm)
			if err != nil {
				writeErrors = append(writeErrors, err.Error())
			}
		}
		if len(writeErrors) != 0 {
			_, _ = fmt.Fprintf(context.App.Writer, "There were errors writing the files:\n: %s", strings.Join(writeErrors, "\n"))
		}
	}()
	createConnector(context).Run(collectResults,
		[]tukan.PhoneAction{tukan.PreparePhoneBookDownload(collectBooks)},
		collectResults)
	close(channel)
	close(bookChannel)
	wg.Wait()
}

func phoneBookFileName(address string) string {
	regex := regexp.MustCompile("https?://")
	result := regex.ReplaceAllString(address, "")
	result = strings.ReplaceAll(result, ":", "_")
	return "phonebook_" + result + ".xml"
}

func downloadParameters(context *cli.Context) {
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

	parametersChannel := make(chan *tukan.ParametersResult)
	collectParameters := func(result *tukan.ParametersResult) {
		parametersChannel <- result
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
		for parametersResult := range parametersChannel {
			fileName := parametersFileName(parametersResult.Address)
			path := filepath.Join(targetDirectory, fileName)
			err := ioutil.WriteFile(path, []byte(parametersResult.Parameters.FunctionKeys.String()), os.ModePerm)
			if err != nil {
				writeErrors = append(writeErrors, err.Error())
			}
		}
		if len(writeErrors) != 0 {
			_, _ = fmt.Fprintf(context.App.Writer, "There were errors writing the files:\n: %s", strings.Join(writeErrors, "\n"))
		}
	}()
	createConnector(context).
		Run(collectResults,
			[]tukan.PhoneAction{tukan.PrepareParameterDownload(collectParameters)},
			collectResults)
	close(channel)
	close(parametersChannel)
	wg.Wait()
}

func parametersFileName(address string) string {
	regex := regexp.MustCompile("https?://")
	result := regex.ReplaceAllString(address, "")
	result = strings.ReplaceAll(result, ":", "_")
	return "parameters_" + result + ".txt"
}
