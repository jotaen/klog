package cli

import (
	"encoding/json"
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"io"
	"net/http"
	"strings"
	"time"
)

type Version struct {
	NoCheck bool `name:"no-check" help:"Don’t check online for updates."` // used for the smoke test
	util.QuietArgs
}

func (opt *Version) Help() string {
	return `
If you don’t use a package manager for managing your klog installation, you can subscribe to the release notifications on the Github repository (https://github.com/jotaen/klog).
`
}

const KLOG_WEBSITE_URL = "https://klog.jotaen.net"

var versionCheckers = []versionChecker{
	{"https://klog.jotaen.net/versions/latest.json", &versionFromJotaen{}},
	{"https://api.github.com/repos/jotaen/klog/releases/latest", &versionFromGithub{}},
}

func (opt *Version) Run(ctx app.Context) app.Error {
	if ctx.Meta().Version == "" {
		return app.NewError(
			"Cannot check version",
			"There is no version information embedded in your binary. Please check manually for updates on "+KLOG_WEBSITE_URL,
			nil,
		)
	}

	if opt.Quiet {
		ctx.Print(ctx.Meta().Version + "\n")
		return nil
	}
	ctx.Print("Command line tool: " + ctx.Meta().Version)
	ctx.Print("  [" + ctx.Meta().SrcHash + "]\n")
	ctx.Print("File format: version " + klog.SPEC_VERSION + "\n")

	if opt.NoCheck {
		return nil
	}

	ctx.Print("\nChecking for newer version...")
	stopTicker := make(chan bool)
	go progressTicker(func() {
		ctx.Print(".")
	}, stopTicker)

	origin, v := fetchVersionInfo(versionCheckers)
	close(stopTicker)
	ctx.Print("\n")
	if v == nil {
		return app.NewError(
			"Failed to retrieve version information.",
			"Please try again later, or check manually at "+KLOG_WEBSITE_URL,
			nil,
		)
	}

	ctx.Debug(func() {
		ctx.Print("Fetched from: " + origin.url + "\n")
	})
	if v.Version() == ctx.Meta().Version && ctx.Meta().SrcHash == v.SrcHash() {
		ctx.Print("You already have the latest version!\n")
	} else {
		ctx.Print(fmt.Sprintf("New version available: %s  [%s]\n", v.Version(), v.SrcHash()))
		downloadLinkPath := ""
		if v.DownloadLinkPath() != "" {
			downloadLinkPath = "/" + v.DownloadLinkPath()
		}
		ctx.Print("See: " + KLOG_WEBSITE_URL + downloadLinkPath + "\n")
	}
	return nil
}

type versionInfo interface {
	Version() string
	SrcHash() string
	DownloadLinkPath() string
	IsValid() bool
}

type versionChecker struct {
	url       string
	structure versionInfo
}

func progressTicker(onTick func(), stop chan bool) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for i := 1; i <= 20; i++ {
		select {
		case <-ticker.C:
			onTick()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

// fetchVersionInfo requests version info from the origins by trying them
// one after the other. It returns the first response that succeeds.
func fetchVersionInfo(origins []versionChecker) (*versionChecker, versionInfo) {
	for _, origin := range origins {
		req, err := http.NewRequest(http.MethodGet, origin.url, nil)
		if err != nil {
			continue
		}
		res, err := (&http.Client{
			Timeout: time.Second * 5,
		}).Do(req)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			continue
		}
		v := origin.structure
		err = json.Unmarshal(body, &v)
		if err != nil || !v.IsValid() {
			continue
		}
		return &origin, v
	}
	return nil, nil
}

// versionFromJotaen checks the version from klog.jotaen.net
type versionFromJotaen struct {
	Version_          string `json:"version"`
	BuildHash_        string `json:"src_hash"`
	DownloadLinkPath_ string `json:"download_link_path"`
}

func (v *versionFromJotaen) Version() string { return v.Version_ }
func (v *versionFromJotaen) SrcHash() string {
	if len(v.BuildHash_) < 7 {
		return ""
	}
	return v.BuildHash_[:7]
}
func (v *versionFromJotaen) DownloadLinkPath() string { return v.DownloadLinkPath_ }
func (v *versionFromJotaen) IsValid() bool {
	return strings.HasPrefix(v.Version_, "v") && len(v.BuildHash_) >= 7
}

// versionFromGithub checks the version from github.com
type versionFromGithub struct {
	Tag        string `json:"tag_name"`
	CommitHash string `json:"target_commitish"`
}

func (v *versionFromGithub) Version() string { return v.Tag }
func (v *versionFromGithub) SrcHash() string {
	if len(v.CommitHash) < 7 {
		return v.CommitHash
	}
	return v.CommitHash[:7]
}
func (v *versionFromGithub) DownloadLinkPath() string {
	return ""
}
func (v *versionFromGithub) IsValid() bool {
	return strings.HasPrefix(v.Tag, "v") && len(v.CommitHash) >= 7
}
