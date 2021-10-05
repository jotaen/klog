package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"io"
	"net/http"
	"time"
)

type Version struct {
	NoCheck bool `name:"no-check" help:"Don’t check online for updates"` // used for the smoke test
	lib.QuietArgs
}

func (opt *Version) Run(ctx app.Context) error {
	if opt.Quiet {
		ctx.Print(ctx.MetaInfo().Version + "\n")
		return nil
	}
	ctx.Print("Command line tool: " + ctx.MetaInfo().Version)
	ctx.Print("  [" + ctx.MetaInfo().BuildHash + "]\n")
	ctx.Print("File format: version " + SPEC_VERSION + "\n")

	if opt.NoCheck {
		return nil
	}
	ctx.Print(fmt.Sprintf("\nChecking for newer version...\n"))
	v := fetchVersionInfo("https://api.github.com/repos/jotaen/klog/releases/latest")
	if v == nil {
		return errors.New("Failed to check for new version, please try again later")
	}
	if v.Version() == ctx.MetaInfo().Version && ctx.MetaInfo().BuildHash == v.BuildHash() {
		ctx.Print(fmt.Sprintf("You already have the latest version!\n"))
	} else {
		ctx.Print(fmt.Sprintf("New version available: %s  [%s]\n", v.Version(), v.BuildHash()))
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
	body, err := io.ReadAll(res.Body)
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
	Tag    string `json:"tag_name"`
	Commit string `json:"target_commitish"`
}

func (v *versionInfo) Version() string { return v.Tag }
func (v *versionInfo) BuildHash() string {
	if len(v.Commit) < 7 {
		return ""
	}
	return v.Commit[:7]
}
