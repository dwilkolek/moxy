package app

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/dwilkolek/moxy/internal/logger"
	"github.com/inconshreveable/go-update"
)

var version string

func Update() error {
	var url string
	logger := logger.New("Moxy")
	var client = &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		url = req.URL.String()
		return nil
	}

	_, err := client.Get("https://github.com/dwilkolek/moxy/releases/latest")
	if err != nil {
		return err
	}

	if runtime.GOARCH == "amd64" {
		url := "/moxy-" + runtime.GOOS + "-" + runtime.GOARCH
		if runtime.GOOS == "windows" {
			url = url + ".exe"
		}
	}

	url = strings.Replace(url, "/tag/", "/download/", -1)
	logger.Printf("Latest version available at %s \n", url)
	if strings.Contains(url, version) {
		logger.Printf("Up to date. \n")
		return nil
	}

	logger.Printf("Downloading update... %s \n", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})

	if err != nil {
		logger.Println("Downloading update failed")
	}
	return err

}
