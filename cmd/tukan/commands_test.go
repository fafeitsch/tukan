package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fafeitsch/Tukan/tukan/mock"
	"github.com/fafeitsch/Tukan/tukan/params"
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
	"time"
)

var username = "doe"
var password = "admin123"

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

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

	assert.Equal(t, 221, len(got), "length of message is wrong")
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

	number := rand.Int()
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("tukan-test%d", number))
	err := os.Mkdir(tmpDir, os.ModePerm)
	require.NoError(t, err, "no error expected")
	defer func() { _ = os.RemoveAll(tmpDir) }()
	err = ioutil.WriteFile(filepath.Join(tmpDir, phoneBookFileName(server1.URL)), []byte("phone book 1"), os.ModePerm)
	require.NoError(t, err, "no error expected")
	err = ioutil.WriteFile(filepath.Join(tmpDir, phoneBookFileName(server2.URL)), []byte("phone book 2"), os.ModePerm)
	require.NoError(t, err, "no error expected")

	t.Run("success", func(t *testing.T) {
		var buff bytes.Buffer
		flags.String(sourceDirFlagName, tmpDir, "")
		ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
		uploadPhoneBook(ctx)
		got := buff.String()

		assert.Equal(t, 254, len(got), "length of message is wrong")
		assert.Containsf(t, got, server1.URL, "should contain server1 URL %s", server1.URL)
		assert.Containsf(t, got, server2.URL, "should contain server2 URL %s", server2.URL)
		assert.Equal(t, "phone book 1\n", phone1.Phonebook, "phone book 1 not uploaded correctly")
	})
	t.Run("file not found", func(t *testing.T) {
		var buff bytes.Buffer
		flags := flag.NewFlagSet("", flag.PanicOnError)
		flags.String(sourceDirFlagName, filepath.Join(tmpDir, "not_existing"), "")
		flags.String(loginFlagName, username, "")
		flags.String(passwordFlagName, password, "")
		_ = flags.Parse([]string{server1.URL})
		ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
		uploadPhoneBook(ctx)
		assert.Contains(t, buff.String(), "Uploading Phone Book returned error: open "+filepath.Join(tmpDir, "not_existing", phoneBookFileName(server1.URL))+": no such file or directory", "result message in case of error wrong")
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
	saveConfig(ctx)
	got := strings.Split(buff.String(), "\n")

	assert.Equal(t, 4, len(got)-1, "expected two lines of result")
	assert.Equal(t, server1.URL+":", got[0], "message of first download is wrong")
	assert.Equal(t, "\tLogin successful", got[1], "message of first download is wrong")
	fileContent, err := ioutil.ReadFile(filepath.Join(tmpDir, parametersFileName(server1.URL)))
	require.NoError(t, err, "reading the file should not give an error")
	para := params.Parameters{}
	err = json.Unmarshal(fileContent, &para)
	require.NoError(t, err, "no error while unmarshalling expected")
	assert.Equal(t, params.FunctionKey{PhoneNumber: "89-IN", DisplayName: "Linda", CallPickupCode: " #0"}, para.FunctionKeys[0], "function key not downloaded correctly")
}

func TestReplaceFunctionKeys(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Parameters = params.Parameters{
		PhoneModel: "Phone ABC",
		FunctionKeys: []params.FunctionKey{
			{DisplayName: "Linda", PhoneNumber: "89-IN", CallPickupCode: "#0"},
			{DisplayName: "John", PhoneNumber: "90-DS", CallPickupCode: "#0"},
		},
	}
	handler2, phone2 := mock.CreatePhone(username, password)
	phone2.Parameters = params.Parameters{
		PhoneModel: "Phone ABC",
		FunctionKeys: []params.FunctionKey{
			{DisplayName: "John", PhoneNumber: "90-DS", CallPickupCode: "#0"},
			{DisplayName: "Linda", PhoneNumber: "89-IN", CallPickupCode: "#0"},
			{DisplayName: "Hugh", PhoneNumber: "65-ID", CallPickupCode: "***"},
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
	assert.Equal(t, "Eva", got1[0].DisplayName, "displayName of first phone not correctly replaced")
	assert.Equal(t, "89-IN", got1[0].PhoneNumber, "phoneNumber of first phone contact should not be changed")
	assert.Equal(t, params.FunctionKey{DisplayName: "John", PhoneNumber: "90-DS", CallPickupCode: "#0"}, got1[1], "second entry in phone book should still be empty")
	assert.Empty(t, phone1.Parameters.PhoneModel, "Phone Model should be empty now because only the function keys should be sent to the phone")

	got2 := phone2.Parameters.FunctionKeys
	assert.Equal(t, 3, len(got2), "length of function keys of second phone not correct")
	assert.Equal(t, "John", got2[0].DisplayName, "displayName in second phone should not be changed")
	assert.Equal(t, "90-DS", got2[0].PhoneNumber, "phoneNumber in second phone should not be changed")
	assert.Equal(t, "Eva", got2[1].DisplayName, "displayName in second phone not correctly replaced")
	assert.Equal(t, "89-IN", got2[1].PhoneNumber, "phoneNumber of second phone contact should not be changed")
	assert.Equal(t, params.FunctionKey{DisplayName: "Hugh", PhoneNumber: "65-ID", CallPickupCode: "***"}, got2[2], "third entry in phone book should still be empty")
	assert.Empty(t, phone2.Parameters.PhoneModel, "Phone Model should be empty now because only the function keys should be sent to the phone")

	assert.Equal(t, 348, len(buff.String()), "output is wrong")
}

func TestSipOverrideDisplayNames(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Parameters = params.Parameters{
		PhoneModel: "Phone ABC",
		Sip: params.Sips{
			{DisplayName: "John"},
			{DisplayName: ""},
			{DisplayName: "Fabian"},
		},
	}
	handler2, phone2 := mock.CreatePhone(username, password)
	phone2.Parameters = params.Parameters{
		PhoneModel: "Phone ABC",
		Sip: params.Sips{
			{DisplayName: "Fabian"},
		},
	}

	server1 := httptest.NewServer(handler1)
	defer server1.Close()
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	flags.String(replaceFlagName, "999 Eva", "")
	_ = flags.Parse([]string{server1.URL, server2.URL})

	var buff bytes.Buffer
	ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
	SipOverrideDisplayNames(ctx)

	got1 := phone1.Parameters.Sip
	assert.Equal(t, 3, len(got1), "length of sips of first phone not correct")
	assert.Equal(t, "999 Eva", got1[0].DisplayName, "displayName of sip correctly replaced")
	assert.Equal(t, "", got1[1].DisplayName, "second entry in Sips should still be empty")
	assert.Equal(t, "999 Eva", got1[2].DisplayName, "displayName of sip correctly replaced")
	assert.Empty(t, phone1.Parameters.PhoneModel, "Phone Model should be empty now because only the function keys should be sent to the phone")

	got2 := phone2.Parameters.Sip
	assert.Equal(t, 1, len(got2), "length of sips of second phone not correct")
	assert.Equal(t, "999 Eva", got2[0].DisplayName, "displayName in second phone not correctly replaced")
	assert.Empty(t, phone2.Parameters.PhoneModel, "Phone Model should be empty now because only the function keys should be sent to the phone")

	assert.Equal(t, 374, len(buff.String()), "output is wrong")
}

func TestRestore(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Backup = []byte("")
	server1 := httptest.NewServer(handler1)
	defer server1.Close()

	handler2, _ := mock.CreatePhone(username, "secret")
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String(loginFlagName, username, "")
	flags.String(passwordFlagName, password, "")
	_ = flags.Parse([]string{server1.URL, server2.URL})

	number := rand.Int()
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("tukan-test%d", number))
	err := os.Mkdir(tmpDir, os.ModePerm)
	require.NoError(t, err, "no error expected")
	defer func() { _ = os.RemoveAll(tmpDir) }()
	err = ioutil.WriteFile(filepath.Join(tmpDir, backupFileName(server1.URL)), []byte("Model XYZ"), os.ModePerm)
	require.NoError(t, err, "no error expected")
	err = ioutil.WriteFile(filepath.Join(tmpDir, backupFileName(server2.URL)), []byte("Model XYZ"), os.ModePerm)
	require.NoError(t, err, "no error expected")

	t.Run("success", func(t *testing.T) {
		var buff bytes.Buffer
		flags.String(sourceDirFlagName, tmpDir, "")
		ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
		restore(ctx)
		got := buff.String()

		assert.Equal(t, 254, len(got), "length of message is wrong")
		assert.Containsf(t, got, server1.URL, "should contain server1 URL %s", server1.URL)
		assert.Containsf(t, got, server2.URL, "should contain server2 URL %s", server2.URL)
		assert.Equal(t, "Model XYZ", string(phone1.Backup), "parameters not uploaded correctly")
	})
	t.Run("file not found", func(t *testing.T) {
		var buff bytes.Buffer
		flags := flag.NewFlagSet("", flag.PanicOnError)
		flags.String(sourceDirFlagName, filepath.Join(tmpDir, "not_existing"), "")
		flags.String(loginFlagName, username, "")
		flags.String(passwordFlagName, password, "")
		_ = flags.Parse([]string{server1.URL})
		ctx := cli.NewContext(&cli.App{Writer: &buff}, flags, nil)
		restore(ctx)
		assert.Contains(t, buff.String(), "Uploading Parameters returned error: open "+filepath.Join(tmpDir, "not_existing", backupFileName(server1.URL))+": no such file or directory", "result message in case of error wrong")
	})
}

func TestBackup(t *testing.T) {
	handler1, phone1 := mock.CreatePhone(username, password)
	phone1.Backup = []byte("this is my backup of phone1")
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
	fileContent, err := ioutil.ReadFile(filepath.Join(tmpDir, backupFileName(server1.URL)))
	require.NoError(t, err, "reading the file should not give an error")
	assert.Equal(t, phone1.Backup, fileContent, "downloaded cfg not correct")
}
