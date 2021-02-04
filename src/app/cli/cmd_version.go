package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"klog/app"
	"net/http"
	"time"
)

type Version struct{}

func (args *Version) Run(ctx *app.Context) error {
	fmt.Printf("Version: %s\n(Build %s)\n\n", ctx.MetaInfo().Version, ctx.MetaInfo().BuildHash)

	fmt.Printf("Checking for newer version...\n")
	v := fetchVersionInfo("https://klog.jotaen.net/latest-stable-version.json")
	if v == nil {
		return errors.New("Failed to check for new version, please try again later")
	}
	if v.Latest.Version == ctx.MetaInfo().Version && v.Latest.BuildHash == ctx.MetaInfo().BuildHash {
		fmt.Printf("You already have the latest version!\n")
	} else {
		fmt.Printf("New version available: %s (%s)\n", v.Latest.Version, v.Latest.BuildHash)
		fmt.Printf("See: https://github.com/jotaen/klog\n")
	}
	return nil
}

func fetchVersionInfo(url string) *versionInfo {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := (&http.Client{
		Timeout: time.Second * 5,
	}).Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	v := &versionInfo{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil
	}
	return v
}

type versionInfo struct {
	Latest struct {
		Version   string `json:"version"`
		BuildHash string `json:"buildHash"`
	} `json:"latest"`
}
