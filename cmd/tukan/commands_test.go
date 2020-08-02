package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/fafeitsch/Tukan/tukan/mock"
	"github.com/fafeitsch/Tukan/tukan/params"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var username = "doe"
var password = "admin123"

func TestScan(t *testing.T) {
	handler1, _ := mock.CreatePhone(username, password)
	server1 := httptest.NewServer(handler1)
	defer server1.Close()

	handler2, _ := mock.CreatePhone(username, "secret")
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	_ = flags.Parse([]string{server1.URL, server2.URL})

	var buff bytes.Buffer
	ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
	scan(ctx)
	got := buff.String()

	assert.Equal(t, 176, len(got), "length of message is wrong")
	assert.Containsf(t, got, server1.URL, "should contain server1 URL %s", server1.URL)
	assert.Containsf(t, got, server2.URL, "should contain server2 URL %s", server1.URL)
}

func TestUploadPhoneBook(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Phonebook = ""
	server1 := httptest.NewServer(handler1)
	defer server1.Close()

	handler2, _ := mock.CreatePhone(username, "secret")
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	_ = flags.Parse([]string{server1.URL, server2.URL})

	t.Run("success", func(t *testing.T) {
		var buff bytes.Buffer
		flags.String(fileFlagName, "./test-resources/phonebook.txt", "")
		ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
		uploadPhoneBook(ctx)
		got := buff.String()

		assert.Equal(t, 209, len(got), "length of message is wrong")
		assert.Containsf(t, got, server1.URL, "should contain server1 URL %s", server1.URL)
		assert.Containsf(t, got, server2.URL, "should contain server2 URL %s", server1.URL)
	})
	t.Run("file not found", func(t *testing.T) {
		var buff bytes.Buffer
		flags := flag.NewFlagSet("", flag.PanicOnError)
		flags.String(fileFlagName, "./test-resources/not-existing.txt", "")
		ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
		uploadPhoneBook(ctx)
		assert.Equal(t, "could not load phone book file: open ./test-resources/not-existing.txt: no such file or directory", buff.String(), "result message in case of error wrong")
	})
}

func TestDownloadPhoneBook(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Phonebook = "phone book of telephone 1"
	server1 := httptest.NewServer(handler1)
	defer server1.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	_ = flags.Parse([]string{server1.URL})

	var buff bytes.Buffer
	number := rand.Int()
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("tukan-test%d", number))
	defer func() { _ = os.RemoveAll(tmpDir) }()
	flags.String(targetDirFlagName, tmpDir, "")
	ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
	downloadPhoneBook(ctx)
	got := strings.Split(buff.String(), "\n")

	assert.Equal(t, 4, len(got)-1, "expected two lines of result")
	assert.Equal(t, server1.URL+":", got[0], "message of first download is wrong")
	assert.Equal(t, "\tLogin successful", got[1], "message of first download is wrong")
	fileContent, err := ioutil.ReadFile(filepath.Join(tmpDir, phoneBookFileName(server1.URL)))
	require.NoError(t, err, "reading the file should not give an error")
	assert.Equal(t, phone1.Phonebook, string(fileContent), "file content is wrong")
}

func TestDownloadParameters(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Parameters = params.Parameters{
		FunctionKeys: []params.FunctionKey{
			{DisplayName: "Linda", PhoneNumber: "89-IN", CallPickupCode: " #0"},
			{},
		},
	}
	server1 := httptest.NewServer(handler1)
	defer server1.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	_ = flags.Parse([]string{server1.URL})

	var buff bytes.Buffer
	number := rand.Int()
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("tukan-test%d", number))
	defer func() { _ = os.RemoveAll(tmpDir) }()
	flags.String(targetDirFlagName, tmpDir, "")
	ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
	backup(ctx)
	got := strings.Split(buff.String(), "\n")

	assert.Equal(t, 4, len(got)-1, "expected two lines of result")
	assert.Equal(t, server1.URL+":", got[0], "message of first download is wrong")
	assert.Equal(t, "\tLogin successful", got[1], "message of first download is wrong")
	fileContent, err := ioutil.ReadFile(filepath.Join(tmpDir, parametersFileName(server1.URL)))
	require.NoError(t, err, "reading the file should not give an error")
	para := params.Parameters{}
	err = yaml.Unmarshal(fileContent, &para)
	require.NoError(t, err, "no error while unmarshalling expected")
	assert.Equal(t, params.FunctionKey{PhoneNumber: "89-IN", DisplayName: "Linda", CallPickupCode: " #0"}, para.FunctionKeys[0], "function key not downloaded correctly")
}

func TestReplaceFunctionKeys(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Parameters = params.Parameters{
		FunctionKeys: []params.FunctionKey{
			{DisplayName: "Linda", PhoneNumber: "89-IN", CallPickupCode: "#0"},
			{},
		},
	}
	handler2, phone2 := mock.CreatePhone(username, password)
	phone2.Parameters = params.Parameters{
		FunctionKeys: []params.FunctionKey{
			{DisplayName: "John", PhoneNumber: "90-DS", CallPickupCode: "#0"},
			{DisplayName: "Linda", PhoneNumber: "89-IN", CallPickupCode: "#0"},
			{},
		},
	}

	server1 := httptest.NewServer(handler1)
	defer server1.Close()
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	flags.String(originalFlagName, "Linda", "")
	flags.String(replaceFlagName, "Eva", "")
	_ = flags.Parse([]string{server1.URL, server2.URL})

	var buff bytes.Buffer
	ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
	replaceFunctionKeys(ctx)

	got1 := phone1.Parameters.FunctionKeys
	assert.Equal(t, 2, len(got1), "length of function keys of first phone not correct")
	assert.Equal(t, "Eva", got1[0].DisplayName.String(), "displayName of first phone not correctly replaced")
	assert.Equal(t, "89-IN", got1[0].PhoneNumber.String(), "phoneNumber of first phone contact should not be changed")
	assert.Equal(t, params.FunctionKey{}, got1[1], "second entry in phone book should still be empty")

	got2 := phone2.Parameters.FunctionKeys
	assert.Equal(t, 3, len(got2), "length of function keys of second phone not correct")
	assert.Equal(t, "John", got2[0].DisplayName.String(), "displayName in second phone should not be changed")
	assert.Equal(t, "90-DS", got2[0].PhoneNumber.String(), "phoneNumber in second phone should not be changed")
	assert.Equal(t, "Eva", got2[1].DisplayName.String(), "displayName in second phone not correctly replaced")
	assert.Equal(t, "89-IN", got2[1].PhoneNumber.String(), "phoneNumber of second phone contact should not be changed")
	assert.Equal(t, params.FunctionKey{}, got2[2], "third entry in phone book should still be empty")

	assert.Equal(t, 348, len(buff.String()), "output is wrong")
}
