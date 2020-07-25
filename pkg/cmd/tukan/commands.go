package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	channel := make(chan commentedResult)
	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).Run(actionLogin.handler(channel), func(p *tukan.Phone) {}, actionLogout.handler(channel))
	close(channel)
	wg.Wait()
}

func uploadPhoneBook(context *cli.Context) {
	file := context.String(fileFlagName)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not load phone book file: %v", err)
		return
	}

	channel := make(chan commentedResult)

	uploadHandler := actionUploadPhoneBook.handler(channel)
	upload := func(p *tukan.Phone) {
		err := p.UploadPhoneBook(string(content))
		uploadHandler(&tukan.PhoneResult{Address: p.Address, Error: err})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).Run(actionLogin.handler(channel),
		upload,
		actionLogout.handler(channel))
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
	channel := make(chan commentedResult)

	handler := actionDownloadPhoneBook.handler(channel)
	download := func(p *tukan.Phone) {
		book, err := p.DownloadPhoneBook()
		handler(&tukan.PhoneResult{Address: p.Address, Error: err})
		if err == nil && book != nil {
			fileName := phoneBookFileName(p.Address)
			path := filepath.Join(targetDirectory, fileName)
			err := ioutil.WriteFile(path, []byte(*book), os.ModePerm)
			if err != nil {
				comment := fmt.Sprintf("Downloaded content could not be written to file:%v", err)
				channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: err}, comment: comment}
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).Run(actionLogin.handler(channel),
		download,
		actionLogout.handler(channel))
	close(channel)
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
	channel := make(chan commentedResult)

	handler := actionDownloadParameters.handler(channel)
	download := func(p *tukan.Phone) {
		params, err := p.DownloadParameters()
		handler(&tukan.PhoneResult{Address: p.Address, Error: err})
		if err == nil && params != nil {
			fileName := parametersFileName(p.Address)
			file, err := os.OpenFile(filepath.Join(targetDirectory, fileName), os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				comment := fmt.Sprintf("Downloaded content could not be written to file:%v", err)
				channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: err}, comment: comment}
				return
			}
			defer func() { _ = file.Close() }()
			writer := csv.NewWriter(file)
			headerErr := writer.Write([]string{"DisplayName", "PhoneNumber", "CallPickupCode", "Type"})
			if headerErr != nil {
				comment := fmt.Sprintf("Downloaded content could not be written to file:%v", headerErr)
				channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: headerErr}, comment: comment}
				return
			}
			for _, fnKey := range params.FunctionKeys {
				err := writer.Write([]string{fnKey.DisplayName.Value, fnKey.PhoneNumber.Value, fnKey.CallPickupCode.Value, fnKey.Type.Value})
				if err != nil {
					comment := fmt.Sprintf("Downloaded content could not be written to file:%v", err)
					channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: err}, comment: comment}
					return
				}
			}
			writer.Flush()
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).
		Run(actionLogin.handler(channel),
			download,
			actionLogout.handler(channel))
	close(channel)
	wg.Wait()
}

func replaceFunctionKeys(context *cli.Context) {
	original := context.String(originalFlagName)
	replace := context.String(replaceFlagName)

	channel := make(chan commentedResult)

	downloadHandler := actionDownloadParameters.handler(channel)
	uploadHandler := actionUploadParameters.handler(channel)
	replaceOperation := func(p *tukan.Phone) {
		params, err := p.DownloadParameters()
		downloadHandler(&tukan.PhoneResult{Address: p.Address, Error: err})
		if err != nil {
			return
		}
		upload, changed := params.TransformFunctionKeyNames(original, replace)
		comment := fmt.Sprintf("%s (changed keys): %v", actionReplaceFunctionKeys.String(), changed)
		channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: err}, comment: comment}
		err = p.UploadParameters(upload)
		uploadHandler(&tukan.PhoneResult{Address: p.Address, Error: err})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).
		Run(actionLogin.handler(channel),
			replaceOperation,
			actionLogout.handler(channel))
	close(channel)
	wg.Wait()
}

func parametersFileName(address string) string {
	regex := regexp.MustCompile("https?://")
	result := regex.ReplaceAllString(address, "")
	result = strings.ReplaceAll(result, ":", "_")
	return "parameters_" + result + ".csv"
}
