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

type Version struct {
	NoCheck bool `name:"no-check" help:"Don’t check online for updates"` // used for the smoke test
}

func (args *Version) Run(ctx app.Context) error {
	ctx.Print("Command line tool: " + ctx.MetaInfo().Version)
	ctx.Print("  [" + ctx.MetaInfo().BuildHash + "]\n")
	ctx.Print("File format: version 1 (RFC)\n")

	if args.NoCheck {
		return nil
	}
	ctx.Print(fmt.Sprintf("\nChecking for newer version...\n"))
	v := fetchVersionInfo("https://klog.jotaen.net/latest-stable-version.json")
	if v == nil {
		return errors.New("Failed to check for new version, please try again later")
	}
	if v.Latest.Version == ctx.MetaInfo().Version && v.Latest.BuildHash == ctx.MetaInfo().BuildHash {
		ctx.Print(fmt.Sprintf("You already have the latest version!\n"))
	} else {
		ctx.Print(fmt.Sprintf("New version available: %s  [%s]\n", v.Latest.Version, v.Latest.BuildHash))
		ctx.Print(fmt.Sprintf("See: https://github.com/jotaen/klog\n"))
	}
	return nil
}

func fetchVersionInfo(url string) *versionInfo {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := (&http.Client{
		Timeout: time.Second * 7,
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
