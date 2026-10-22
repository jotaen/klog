package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/jotaen/klog/klog/app/cli/prettify"
)

type Version struct {
	NoCheck bool `name:"no-check" help:"Don’t check online for updates."` // used e.g. for the smoke test or other automated contexts
	args.QuietArgs
}

func (opt *Version) Help() string {
	return `
If you don’t use a package manager for managing your klog installation, you can subscribe to the release notifications on the Github repository (https://github.com/jotaen/klog).
`
}

const KLOG_WEBSITE_URL = "https://klog.jotaen.net"

type versionFetcher interface {
	url() string
	decode() app.BuildInfo
}

var versionFetchers = []versionFetcher{
	&versionFromJotaen{},
	&versionFromGithub{},
}

func (opt *Version) Run(ctx app.Context) app.Error {
	buildInfo := ctx.Meta().BuildInfo

	if opt.Quiet {
		if !buildInfo.IsDetermined() {
			return app.NewError(
				"Version unknown",
				"There is no version information embedded into your binary.",
				nil,
			)
		}
		ctx.Print(buildInfo.Version + "\n")
		return nil
	}

	if !buildInfo.IsDetermined() {
		styler, _ := ctx.Serialise()
		ctx.Print(prettify.PrettifyWarning("There is no version information embedded in your klog binary\n", styler))
	}

	ctx.Print("Command line tool: " + buildInfo.Version)
	ctx.Print(fmt.Sprintf("  [%s, %s]\n", buildInfo.SrcHash, buildInfo.Date))
	ctx.Print("File format: version " + klog.SPEC_VERSION + "\n")

	if opt.NoCheck {
		return nil
	}

	ctx.Print("\nChecking for newer version...")
	stopTicker := make(chan bool)
	go progressTicker(func() {
		ctx.Print(".")
	}, stopTicker)

	fetchUrl, latestRelease := fetchVersionInfo(versionFetchers)
	close(stopTicker)
	ctx.Print("\n")
	if fetchUrl == "" {
		return app.NewError(
			"Failed to retrieve version information.",
			"Please try again later, or check manually at "+KLOG_WEBSITE_URL,
			nil,
		)
	}

	ctx.Debug(func() {
		ctx.Print("Fetched from: " + fetchUrl + "\n")
	})

	if latestRelease.Version == buildInfo.Version && latestRelease.SrcHash == buildInfo.SrcHash {
		ctx.Print("You already have the latest version!\n")
	} else {
		ctx.Print(fmt.Sprintf("New version available: %s  [%s, %s]\n", latestRelease.Version, latestRelease.SrcHash, latestRelease.Date))
		ctx.Print("See: " + KLOG_WEBSITE_URL + "\n")
	}
	return nil
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
func fetchVersionInfo(fetcher []versionFetcher) (string, app.BuildInfo) {
	for _, origin := range fetcher {
		req, err := http.NewRequest(http.MethodGet, origin.url(), nil)
		if err != nil {
			continue
		}
		res, err := (&http.Client{
			Timeout: time.Second * 5,
		}).Do(req)
		if err != nil {
			continue
		}
		if res.StatusCode != http.StatusOK {
			_ = res.Body.Close()
			continue
		}
		body, err := io.ReadAll(io.LimitReader(res.Body, 64*1024))
		_ = res.Body.Close()
		if err != nil {
			continue
		}
		err = json.Unmarshal(body, &origin)
		if err != nil {
			// Skip if JSON cannot be decoded.
			continue
		}
		buildInfo := origin.decode()
		if !buildInfo.IsDetermined() {
			// Skip if build info data isn’t legit.
			continue
		}
		return origin.url(), buildInfo
	}
	return "", app.NewBuildInfo("", "", "")
}

// versionFromJotaen checks the version via klog.jotaen.net
type versionFromJotaen struct {
	Version_     string `json:"version"`
	BuildHash_   string `json:"src_hash"`
	ReleaseDate_ string `json:"release_date"`
}

func (v *versionFromJotaen) url() string {
	return "https://klog.jotaen.net/versions/latest.json"
}
func (v *versionFromJotaen) decode() app.BuildInfo {
	return app.NewBuildInfo(v.Version_, v.BuildHash_, v.ReleaseDate_)
}

// versionFromGithub checks the version via github.com
type versionFromGithub struct {
	Tag        string `json:"tag_name"`
	CommitHash string `json:"target_commitish"`
	CreatedAt  string `json:"created_at"` // Refers to the build date, not the publishing date.
}

func (v *versionFromGithub) url() string {
	return "https://api.github.com/repos/jotaen/klog/releases/latest"
}
func (v *versionFromGithub) decode() app.BuildInfo {
	return app.NewBuildInfo(v.Tag, v.CommitHash, v.CreatedAt)
}
