package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
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
