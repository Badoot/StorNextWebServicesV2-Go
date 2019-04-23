package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Setting these variables here since this is just an example that I
// built in my StorNext lab, but you obviously want to provide credentials
// more securely than this in production.
const username = "webuser"
const password = "password"

// I pre-downloaded the self-signed certification from my lab's metadata
// controller with this command:
//    openssl s_client -showcerts -connect mdc01.badoot.local:443 \
//    </dev/null 2>/dev/null|openssl x509 -outform PEM >mycertfile.pem
//
const urlPrefix = "https://mdc01.badoot.local/sws/v2"
const certPath = "./mycertfile.pem"

// I like json
const format = "json"

func main() {
	// Eval length of arguments provided via CLI
	switch len(os.Args) {
	// If just "stornext.go" with no addtional arguments, print usage.
	case 1:
		usage()
	// If an argument is provided, run the fuction for the action requested.
	case 2:
		switch os.Args[1] {
		// System info
		case "info":
			systemInfo()
		// Media info
		case "fsmedinfo":
			fsMedInfo()
		default:
			usage()
		}
	case 3:
		// If 2 arguments are provided, argument 1 is action, 2 is file name
		file := os.Args[2]
		switch arg := os.Args[1]; arg {
		// File info
		case "fsfileinfo":
			fsFileInfo(file)
		// Retrieve File
		case "fsretrieve":
			fsRetrieve(file)
		// Store File
		case "fsstore":
			fsStore(file)
		// Store File
		case "fsrmdiskcopy":
			fsRmDiskCopy(file)
		default:
			usage()
		}
	default:
		usage()
	}
}

func usage() {
	// Print usage
	usage := ("\nUsage: go run stornext.go [action] [filename]" +
		"\n\nActions :\n\n" +
		"info 			: Retrieves the latest status of system components.\n" +
		"fsmedinfo 		: Generate a report on media based on their current status.\n" +
		"fsfileinfo [filename] 	: Generate a report about files known to the Tertiary Storage Manager.\n" +
		"fsstore [filename] 	: Expedite the storage of a file that currently resides on disk to media.\n" +
		"fsretrieve [filname] 	: Retrieve truncated files from media and place on disk.\n" +
		"fsrmdiskcopy [filename] : Remove the copy of a file from disk after the file was stored to a medium.\n")
	fmt.Println(usage)
}

func buildClient() *http.Client {
	// Set up our own certificate pool
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Load our trusted certificate path
	pemData, err := ioutil.ReadFile(certPath)
	if err != nil {
		panic(err)
	}
	// Check if pem can be loaded
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
	if !ok {
		panic("Couldn't load PEM data")
	}
	return client
}

func buildURL(action string) string {
	// Build the URL for your API call
	var url = urlPrefix + action + "&username=" + username + "&password=" + password + "&format=" + format
	return url
}

func printResponse(url string) {
	// Get response from the client
	resp, err := buildClient().Get(url)
	if err != nil {
		panic(err)
	}
	// Read the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// format the bytes returned as a string
	respBody := fmt.Sprintf("%s", body)
	// Print URL at the top
	fmt.Println("URL: " + url)
	// Print response body
	fmt.Println(respBody)
}

func systemInfo() {
	// Retrieves the latest status of system components.
	url := buildURL("/system/info?")
	printResponse(url)
}

func fsMedInfo() {
	// Generate a report on media based on their current status.
	url := buildURL("/fsmedinfo?verbose=true")
	printResponse(url)
}

func fsFileInfo(filename string) {
	// Generate a report about files known to the Tertiary Storage Manager.
	url := buildURL("/file/fsfileinfo?file=" + filename)
	printResponse(url)
}

func fsStore(filename string) {
	// Expedite the storage of a file that currently resides on disk to media.
	url := buildURL("/file/fsstore?file=" + filename)
	printResponse(url)
}

func fsRetrieve(filename string) {
	// Retrieve truncated files from media and place on disk.
	url := buildURL("/file/fsretrieve?file=" + filename)
	printResponse(url)
}

func fsRmDiskCopy(filename string) {
	// Remove the copy of a file from disk after the file was stored to a medium.
	url := buildURL("/file/fsrmdiskcopy?file=" + filename)
	printResponse(url)
}
