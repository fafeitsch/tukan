package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

const loginFlagName = "login"
const passwordFlagName = "password"
const portFlagName = "port"
const timeoutFlagName = "timeout"
const verboseFlagName = "verbose"
const targetDirFlagName = "targetDir"
const sourceDirFlagName = "sourceDir"
const originalFlagName = "original"
const replaceFlagName = "replace"

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Tukan REST Client for IP-Phones"
	app.Usage = "This application connects to the REST endpoints of some VoIP telephones and offers to upload/download data."

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
		Usage:  "Scans an IP range for IP phones.",
		Action: scan,
	}

	phoneBookUploadCommand := cli.Command{
		Name:  "pb-up",
		Usage: "Uploads a phone book to a set of VoIP phones.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: sourceDirFlagName, Required: true, Usage: "The directory where the phone books to upload can be found.", TakesFile: true},
		},
		Action: uploadPhoneBook,
	}

	phonebookDownloadCommand := cli.Command{
		Name:  "pb-down",
		Usage: "Downloads the phone books from a set of VoIP phones and stores them in files.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: targetDirFlagName, Required: true, Usage: "The directory where the downloaded phonebooks are saved."},
		},
		Action: downloadPhoneBook,
	}

	downloadCommand := cli.Command{
		Name:  "downloadConfig",
		Usage: "Downloads all parameters from the phone and stores them in a json file. Though possible, the downloaded params are only meant for analyzing the settings, not for a complete restore on the phone.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: targetDirFlagName, Required: true, Usage: "The directory where the downloaded parameters are saved."},
		},
		Action: saveConfig,
	}

	backup := cli.Command{
		Name:  "backup",
		Usage: "Downloads a binary backup from the phones which can be restored.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: targetDirFlagName, Required: true, Usage: "The directory where the downloaded parameters are saved."},
		},
		Action: backup,
	}

	restoreCommand := cli.Command{
		Name:  "restore",
		Usage: "Reads a json parameters file and restores it onto the telephone.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: sourceDirFlagName, Required: true, Usage: "The directory where to find the file used for restoring."},
		},
		Action: restore,
	}

	functionKeysReplaceCommand := cli.Command{
		Name:  "fnkeys-replace",
		Usage: "Replaces display names of function keys from VoIP phones.",
		Flags: []cli.Flag{
			replaceFlag,
			originalFlag,
		},
		Action: replaceFunctionKeys,
	}

	sipOverrideDisplayNamesCommand := cli.Command{
		Name:  "sip-override",
		Usage: "Overrides SIP display names if they are not empty",
		Flags: []cli.Flag{
			replaceFlag,
		},
		Action: SipOverrideDisplayNames,
	}

	resetCommand := cli.Command{
		Name:   "reset",
		Usage:  "Resets the whole telephone.",
		Action: reset,
	}

	app.Commands = []cli.Command{scanCommand, phoneBookUploadCommand, phonebookDownloadCommand, downloadCommand, restoreCommand, functionKeysReplaceCommand, resetCommand, backup, sipOverrideDisplayNamesCommand}

	app.Flags = []cli.Flag{loginFlag, passwordFlag, portFlag, timeoutFlag, verboseFlag}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
