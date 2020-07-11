package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
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
	got := strings.Split(buff.String(), "\n")

	assert.Equal(t, 2, len(got)-1, "expected two lines of result")
	want := map[string]bool{
		server2.URL + ": false (could not connect to address 1 \"" + server2.URL + "\": authentication error, status code: 403 with message \"403 Forbidden\")": true,
		server1.URL + ": true (connection established and login successful)":                                                                                    true,
	}
	_, srv1 := want[got[0]]
	_, srv2 := want[got[1]]
	assert.True(t, srv1 && srv2, "both wanted strings should be contained")
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
		got := strings.Split(buff.String(), "\n")

		assert.Equal(t, 2, len(got)-1, "expected two lines of result")
		want := map[string]bool{
			server2.URL + ": false (could not connect to address 1 \"" + server2.URL + "\": authentication error, status code: 403 with message \"403 Forbidden\")": true,
			server1.URL + ": true (Upload successful)": true,
		}
		_, srv1 := want[got[0]]
		_, srv2 := want[got[1]]
		assert.True(t, srv1 && srv2, "both wanted strings should be contained")
		assert.Equal(t, "John Doe: 50\n", phone1.Phonebook, "phone book should be saved on the phone")
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

	assert.Equal(t, 1, len(got)-1, "expected two lines of result")
	assert.Equal(t, server1.URL+": true (Download successful)", got[0], "message of first download is wrong")
	fileContent, err := ioutil.ReadFile(filepath.Join(tmpDir, fileName(server1.URL)))
	require.NoError(t, err, "reading the file should not give an error")
	assert.Equal(t, phone1.Phonebook, string(fileContent), "file content is wrong")
}