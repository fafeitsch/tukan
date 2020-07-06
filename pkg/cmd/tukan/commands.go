package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan"
	"github.com/urfave/cli"
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
	verbose := context.GlobalBool(verboseFlagName)
	number := context.GlobalInt(numberFlagName)

	channel := connectToPhones(context).Scan()
	results := make([]tukan.SimpleResult, 0, number)
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
