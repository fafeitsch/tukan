package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

const loginFlagName = "actionLogin"
const passwordFlagName = "password"
const portFlagName = "port"
const timeoutFlagName = "timeout"
const verboseFlagName = "verbose"
const fileFlagName = "file"
const targetDirFlagName = "targetDir"
const originalFlagName = "original"
const replaceFlagName = "replace"

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Elmeg ip620/630 HTTP Configurator"
	app.Usage = "This application configures some parts of Elmeg ip620/630 telephones"

	var login, password, original, replace string
	var port, timeout int
	var noLogging bool

	loginFlag := cli.StringFlag{Name: loginFlagName, Value: "Admin", Usage: "The actionLogin to be used", Destination: &login}
	passwordFlag := cli.StringFlag{Name: passwordFlagName, Value: "admin", Usage: "The password to be used", Destination: &password}
	portFlag := cli.IntFlag{Name: portFlagName, Value: 80, Usage: "The port to be used to connect to the telephones", Destination: &port}
	verboseFlag := cli.BoolFlag{Name: verboseFlagName, Usage: "Disables the logging and only prints the final results", Destination: &noLogging}
	timeoutFlag := cli.IntFlag{Name: timeoutFlagName, Value: 20, Usage: "Number of seconds to wait for remote connection", Destination: &timeout}
	originalFlag := cli.StringFlag{Name: originalFlagName, Value: "", Usage: "The display name to be replaced", Destination: &original, Required: true}
	replaceFlag := cli.StringFlag{Name: replaceFlagName, Value: "", Usage: "The new display name", Destination: &replace, Required: true}

	scanCommand := cli.Command{
		Name:   "scan",
		Usage:  "Scans an IP range for elmeg ip620/630 and tries to log into them",
		Action: scan,
	}

	phoneBookUploadCommand := cli.Command{
		Name:  "pb-up",
		Usage: "Uploads a phone book to a set of elmeg ip 620/630 phones",
		Flags: []cli.Flag{
			cli.StringFlag{Name: fileFlagName, Required: true, Usage: "The phone book file to be loaded up.", TakesFile: true},
		},
		Action: uploadPhoneBook,
	}

	phonebookDownloadCommand := cli.Command{
		Name:  "pb-down",
		Usage: "Downloads a phone book from a elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			cli.StringFlag{Name: targetDirFlagName, Required: true, Usage: "The directory where the downloaded phonebooks are saved."},
		},
		Action: downloadPhoneBook,
	}

	functionKeysDownloadCommand := cli.Command{
		Name:  "fnKeys-down",
		Usage: "Downloads the function keys from an elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			cli.StringFlag{Name: targetDirFlagName, Required: true, Usage: "The directory where the downloaded phonebooks are saved."},
		},
		Action: downloadParameters,
	}

	functionKeysReplaceCommand := cli.Command{
		Name:  "fnKeys-replace",
		Usage: "Replaces display names of function keys from an elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			replaceFlag,
			originalFlag,
		},
		Action: actionReplaceFunctionKeys,
	}

	app.Commands = []cli.Command{scanCommand, phoneBookUploadCommand, phonebookDownloadCommand, functionKeysDownloadCommand, functionKeysReplaceCommand}

	app.Flags = []cli.Flag{loginFlag, passwordFlag, portFlag, timeoutFlag, verboseFlag}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
