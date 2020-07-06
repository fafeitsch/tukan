package main

import (
	"fmt"
	http2 "github.com/fafeitsch/Tukan/pkg/http"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

const loginFlagName = "login"
const passwordFlagName = "password"
const ipFlagName = "ip"
const numberFlagName = "number"
const portFlagName = "port"
const timeoutFlagName = "timeout"
const verboseFlagName = "verbose"

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Elmeg ip620/630 HTTP Configurator"
	app.Usage = "This application configures some parts of Elmeg ip620/630 telephones"

	var login, password, original, replace string
	var port, timeout int
	var noLogging bool

	loginFlag := cli.StringFlag{Name: loginFlagName, Value: "Admin", Usage: "The login to be used", Destination: &login}
	passwordFlag := cli.StringFlag{Name: passwordFlagName, Value: "admin", Usage: "The password to be used", Destination: &password}
	portFlag := cli.IntFlag{Name: portFlagName, Value: 80, Usage: "The port to be used to connect to the telephones", Destination: &port}
	ipFlag := cli.StringFlag{Name: ipFlagName, Required: false, Usage: "The IP of the first phone to interact with"}
	numberFlag := cli.IntFlag{Name: numberFlagName, Value: 1, Usage: "The number of phones to contact, including IP"}
	noLogFlag := cli.BoolFlag{Name: verboseFlagName, Usage: "Disables the logging and only prints the final results", Destination: &noLogging}
	timeoutFlag := cli.IntFlag{Name: timeoutFlagName, Value: 20, Usage: "Number of seconds to wait for remote connection", Destination: &timeout}
	originalFlag := cli.StringFlag{Name: "original", Value: "", Usage: "The display name to be replaced", Destination: &original, Required: true}
	replaceFlag := cli.StringFlag{Name: "replace", Value: "", Usage: "The new display name", Destination: &replace, Required: true}

	var phoneClient http2.PhoneClient

	scanCommand := cli.Command{
		Name:   "scan",
		Usage:  "Scans an IP range for elmeg ip620/630 and tries to log into them",
		Action: scan,
	}

	phonebookUploadCommand := cli.Command{
		Name:  "pb-up",
		Usage: "Uploads a phone book to a set of elmeg ip 620/630 phones",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file", Required: true, Usage: "The phone book file to up", TakesFile: true},
			ipFlag,
			numberFlag,
		},
		Action: func(c *cli.Context) error {
			content, err := ioutil.ReadFile(c.String("file"))
			if err != nil {
				return fmt.Errorf("could not read phone book file %s: %v", c.String("file"), err)
			}
			phonebook := string(content)
			result := phoneClient.UploadPhoneBook(c.String("ip"), c.Int("number"), phonebook)
			fmt.Printf("%v", result)
			return nil
		},
	}

	phonebookDownloadCommand := cli.Command{
		Name:  "pb-down",
		Usage: "Downloads a phone book from a elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "ip", Required: true, Usage: "The IP of the phone to download to phonebook from."},
		},
		Action: func(c *cli.Context) error {
			_, down := phoneClient.DownloadPhoneBook(c.String("ip"))
			fmt.Printf(down)
			return nil
		},
	}

	functionKeysDownloadCommand := cli.Command{
		Name:  "fnKeys-down",
		Usage: "Downloads the function keys from an elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			ipFlag,
			numberFlag,
		},
		Action: func(c *cli.Context) error {
			result := phoneClient.DownloadFunctionKeys(c.String("ip"), c.Int("number"))
			fmt.Printf("%v", result)
			return nil
		},
	}

	functionKeysReplaceCommand := cli.Command{
		Name:  "fnKeys-replace",
		Usage: "Replaces display names of function keys from an elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			ipFlag,
			numberFlag,
			replaceFlag,
			originalFlag,
		},
		Action: func(c *cli.Context) error {
			result := phoneClient.ReplaceFunctionKeyName(c.String("ip"), c.Int("number"), original, replace)
			fmt.Printf("%v", result)
			return nil
		},
	}

	app.Commands = []cli.Command{scanCommand, phonebookUploadCommand, phonebookDownloadCommand, functionKeysDownloadCommand, functionKeysReplaceCommand}

	app.Flags = []cli.Flag{loginFlag, passwordFlag, portFlag, timeoutFlag, noLogFlag, ipFlag, numberFlag}

	app.Before = func(context *cli.Context) error {
		phoneClient = http2.BuildPhoneClient(port, login, password, timeout)
		if noLogging {
			phoneClient.Logger.SetFlags(0)
			phoneClient.Logger.SetOutput(ioutil.Discard)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
