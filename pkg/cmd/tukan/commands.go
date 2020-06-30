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
	firstIp := context.GlobalString(ipFlagName)
	number := context.GlobalInt(numberFlagName)
	login := context.GlobalString(loginFlagName)
	password := context.GlobalString(passwordFlagName)
	port := context.GlobalInt(portFlagName)
	timeout := context.GlobalInt(timeoutFlagName)
	connector := tukan.Connector{Client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, UserName: login, Password: password}
	addresses := tukan.CreateAddresses("http", firstIp, port, number)
	return connector.MultipleConnect(addresses...)
}

func scan(context *cli.Context) {
	verbose := context.GlobalBool(verboseFlagName)
	number := context.GlobalInt(numberFlagName)

	channel := connectToPhones(context).Scan()
	results := make([]tukan.SimpleResult, 0, number)
	for result := range channel {
		if verbose {
			fmt.Printf("%s\n", result.String())
		}
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return strings.Compare(results[i].Address, results[j].Address) <= 0
	})
	if verbose {
		fmt.Printf("\n")
	}
	for _, result := range results {
		fmt.Printf("%s\n", result.String())
	}
}
