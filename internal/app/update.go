package app

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/dwilkolek/moxy/internal/logger"
	"github.com/inconshreveable/go-update"
)

var version string

func CheckForUpdateUrl() (string, error) {
	var url string
	var client = &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		url = req.URL.String()
		return nil
	}

	_, err := client.Get("https://github.com/dwilkolek/moxy/releases/latest")
	if err != nil {
		return "", err
	}

	if runtime.GOARCH == "amd64" {
		url := "/moxy-" + runtime.GOOS + "-" + runtime.GOARCH
		if runtime.GOOS == "windows" {
			url = url + ".exe"
		}
	}

	url = strings.Replace(url, "/tag/", "/download/", -1)
	if strings.Contains(url, version) {
		return "", nil
	}
	return url, nil
}

func Update() error {
	logger := logger.New("Moxy")
	updateUrl, err := CheckForUpdateUrl()
	if err != nil {
		return err
	}

	logger.Printf("Downloading update... %s \n", updateUrl)

	resp, err := http.Get(updateUrl)
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
