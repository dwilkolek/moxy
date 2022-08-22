package app

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/dwilkolek/moxy/internal/logger"
	"github.com/inconshreveable/go-update"
)

var version string

func Version() string {
	return version
}

func CheckForUpdateUrl() (string, error) {
	var url string
	var client = &http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		url = req.URL.String()
		if runtime.GOARCH == "amd64" {
			url = url + "/moxy-" + runtime.GOOS + "-" + runtime.GOARCH
			if runtime.GOOS == "windows" {
				url = url + ".exe"
			}
		}
		return nil
	}

	_, err := client.Get("https://github.com/dwilkolek/moxy/releases/latest")
	if err != nil {
		return "", err
	}

	url = strings.Replace(url, "/tag/", "/download/", -1)

	if strings.Contains(url, version) {
		return "", nil
	}
	return url, nil
}

func Update() {
	logger := logger.New("Moxy")
	updateUrl, err := CheckForUpdateUrl()

	if err != nil {
		panic("Downloading update failed " + err.Error())
	}
	if updateUrl == "" {
		logger.Printf("Up to date. %s\n", version)
		return
	}
	logger.Printf("Downloading update... %s \n", updateUrl)
	resp, err := http.Get(updateUrl)
	if err != nil {
		panic("Downloading update failed " + err.Error())
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})

	if err != nil {
		panic("Downloading update failed")
	}

	logger.Println("Successful update.")
}
