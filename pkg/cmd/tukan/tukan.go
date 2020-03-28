package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/domain"
	http2 "github.com/fafeitsch/Tukan/pkg/http"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Elmeg ip620/630 HTTP Configurator"
	app.Usage = "This application configures some parts of Elmeg ip620/630 telephones"

	var login, password string
	var port, timeout int
	var noLogging bool
	loginFlag := cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login to be used", Destination: &login}
	passwordFlag := cli.StringFlag{Name: "password", Value: "admin", Usage: "The password to be used", Destination: &password}
	portFlag := cli.IntFlag{Name: "port", Value: 80, Usage: "The port to be used to connect to the telephones", Destination: &port}
	ipFlag := cli.StringFlag{Name: "ip", Required: true, Usage: "The IP of the first phone to interact with"}
	numberFlag := cli.IntFlag{Name: "number", Value: 1, Usage: "The number of phones to contact, including IP"}
	noLogFlag := cli.BoolFlag{Name: "nolog", Usage: "Disables the logging and only prints the final results", Destination: &noLogging}
	timeoutFlag := cli.IntFlag{Name: "timeout", Value: 20, Usage: "Number of seconds to wait for remote connection", Destination: &timeout}

	var phoneClient http2.PhoneClient

	scanCommand := cli.Command{
		Name:  "scan",
		Usage: "Scans an IP range for elmeg ip620/630 and tries to log into them",
		Flags: []cli.Flag{
			numberFlag,
			ipFlag,
		},
		Action: func(c *cli.Context) error {
			result := phoneClient.Scan(c.String("ip"), c.Int(numberFlag.Name))
			fmt.Printf("%v", result)
			return nil
		},
	}

	phonebookUploadCommand := cli.Command{
		Name:  "pb-up",
		Usage: "Uploads a phone book to a set of elmeg ip 620/630 phones",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file", Required: true, Usage: "The phone book file to upload", TakesFile: true},
			cli.StringFlag{Name: "delimiter", Usage: "A string that is not contained in the phone book. Needed for the upload. Must consist of at most 70 bytes of ASCII printable-characters", Value: "XXXXX"},
			ipFlag,
			numberFlag,
		},
		Action: func(c *cli.Context) error {
			payload, err := domain.LoadAndEmbedPhonebook(c.String("file"), c.String("delimiter"))
			if err != nil {
				return fmt.Errorf("could not prepare payload for sending to phones: %v", err)
			}
			result := phoneClient.UploadPhoneBook(c.String("ip"), c.Int("number"), *payload, c.String("delimiter"))
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

	app.Commands = []cli.Command{scanCommand, phonebookUploadCommand, phonebookDownloadCommand, functionKeysDownloadCommand}

	app.Flags = []cli.Flag{loginFlag, passwordFlag, portFlag, timeoutFlag, noLogFlag}

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
