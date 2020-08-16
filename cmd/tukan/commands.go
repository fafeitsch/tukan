package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/tukan"
	"github.com/fafeitsch/Tukan/tukan/params"
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

func reset(context *cli.Context) {
	channel := make(chan commentedResult)
	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	handler := actionReset.handler(channel)
	resetPhone := func(p *tukan.Phone) {
		err := p.Reset()
		handler(&tukan.PhoneResult{Address: p.Address, Error: err})
	}
	// Do nothing with logout because it fails nonetheless (the phone immediately resets itself)
	logoutCallback := func(p *tukan.PhoneResult) {}
	connector := createConnector(context)
	addresses := connector.Addresses
	_, _ = fmt.Fprintf(context.App.Writer, "Do you really want to reset %d phones? Type YES: ", len(addresses))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if input != "YES\n" {
		close(channel)
		wg.Wait()
		return
	}
	connector.Run(actionLogin.handler(channel), resetPhone, logoutCallback)
	close(channel)
	wg.Wait()
}

func uploadPhoneBook(context *cli.Context) {
	sourceDirectory := context.String(sourceDirFlagName)
	channel := make(chan commentedResult)

	uploadHandler := actionUploadPhoneBook.handler(channel)
	upload := func(p *tukan.Phone) {
		fileName := phoneBookFileName(p.Address)
		path := filepath.Join(sourceDirectory, fileName)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			uploadHandler(&tukan.PhoneResult{Address: p.Address, Error: err})
			return
		}
		err = p.UploadPhoneBook(string(content))
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

func saveConfig(context *cli.Context) {
	targetDirectory := context.String(targetDirFlagName)
	err := os.MkdirAll(targetDirectory, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not create target directory: %v", err)
		return
	}
	channel := make(chan commentedResult)

	handler := actionDownloadParameters.handler(channel)
	download := func(p *tukan.Phone) {
		parameters, err := p.DownloadParameters()
		if err == nil && parameters != nil {
			fileName := parametersFileName(p.Address)
			bytes, _ := json.MarshalIndent(&parameters, "", "  ")
			err = ioutil.WriteFile(filepath.Join(targetDirectory, fileName), bytes, os.ModePerm)
		}
		handler(&tukan.PhoneResult{Address: p.Address, Error: err})
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

func backup(context *cli.Context) {
	targetDirectory := context.String(targetDirFlagName)
	err := os.MkdirAll(targetDirectory, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(context.App.Writer, "could not create target directory: %v", err)
		return
	}
	channel := make(chan commentedResult)

	handler := actionBackup.handler(channel)
	backup := func(p *tukan.Phone) {
		data, err := p.Backup()
		if err == nil && data != nil {
			fileName := backupFileName(p.Address)
			fileName = filepath.Join(targetDirectory, fileName)
			err = ioutil.WriteFile(fileName, data, os.ModePerm)
		}
		handler(&tukan.PhoneResult{Address: p.Address, Error: err})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).
		Run(actionLogin.handler(channel),
			backup,
			actionLogout.handler(channel))
	close(channel)
	wg.Wait()
}

func restore(context *cli.Context) {
	sourceDirectory := context.String(sourceDirFlagName)
	channel := make(chan commentedResult)

	handler := actionUploadParameters.handler(channel)
	upload := func(p *tukan.Phone) {
		fileName := backupFileName(p.Address)
		data, err := ioutil.ReadFile(filepath.Join(sourceDirectory, fileName))
		if err != nil {
			handler(&tukan.PhoneResult{Address: p.Address, Error: err})
			return
		}
		err = p.Restore(data)
		handler(&tukan.PhoneResult{Address: p.Address, Error: err})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleResults(&wg, channel, context)
	createConnector(context).
		Run(actionLogin.handler(channel),
			upload,
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
		parameters, err := p.DownloadParameters()
		downloadHandler(&tukan.PhoneResult{Address: p.Address, Error: err})
		if err != nil {
			return
		}
		upload, changed := parameters.FunctionKeys.Transform(params.ReplaceDisplayName(original, replace))
		comment := fmt.Sprintf("%s (changed keys): %v", actionReplaceFunctionKeys.String(), changed)
		channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: err}, comment: comment}
		err = p.UploadParameters(params.Parameters{FunctionKeys: upload})
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

func SipOverrideDisplayNames(context *cli.Context) {
	replace := context.String(replaceFlagName)

	channel := make(chan commentedResult)

	downloadHandler := actionDownloadParameters.handler(channel)
	uploadHandler := actionSipOverrideDisplayName.handler(channel)
	replaceOperation := func(p *tukan.Phone) {
		parameters, err := p.DownloadParameters()
		downloadHandler(&tukan.PhoneResult{Address: p.Address, Error: err})
		if err != nil {
			return
		}
		upload, changed := parameters.Sip.Transform(params.SipOverrideDisplayName(replace))
		comment := fmt.Sprintf("%s (changed sip): %v", actionSipOverrideDisplayName.String(), changed)
		channel <- commentedResult{PhoneResult: &tukan.PhoneResult{Address: p.Address, Error: err}, comment: comment}
		err = p.UploadParameters(params.Parameters{Sip: upload})
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
	return "parameters_" + result + ".json"
}

func backupFileName(address string) string {
	regex := regexp.MustCompile("https?://")
	result := regex.ReplaceAllString(address, "")
	result = strings.ReplaceAll(result, ":", "_")
	return "parameters_" + result + ".cfg"
}
