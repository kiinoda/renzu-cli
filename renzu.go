package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const iniFile = "~/.renzu.rc"

func main() {
	baseURL := getBaseURL()

	if len(os.Args[1:]) == 0 {
		getRenzu(baseURL + "renzu/introspect")
	} else {
		params := strings.Join(os.Args[1:], "/")
		getRenzu(baseURL + params)
	}
}

func expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(user.HomeDir, path[1:]), nil
}

func getBaseURL() string {
	filePath, err := expandPath(iniFile)
	if err != nil {
		panic("Unable to find configuration file")
	}

	f, err := os.Open(filePath)
	if err != nil {
		panic("Could not open configuration file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()

	// baseURL := strings.TrimSuffix(string(firstLine), "\n")
	return firstLine
}

func getRenzu(url string) {
	res, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to Renzu API")
	}

	data, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))
}
