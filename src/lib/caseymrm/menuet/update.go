// +build darwin

package menuet

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func (a *Application) checkForUpdates() {
	checkForRestart()
	ticker := time.NewTicker(24 * time.Hour)
	for ; true; <-ticker.C {
		release := checkForNewRelease(a.AutoUpdate.Repo, a.AutoUpdate.Version)
		if release == nil {
			continue
		}
		button := a.Alert(Alert{
			MessageText:     fmt.Sprintf("New version of %s available", a.Name),
			InformativeText: fmt.Sprintf("Looks like %s of %s is now available- you're running %s", release.TagName, a.Name, a.AutoUpdate.Version),
			Buttons:         []string{"Update now", "Remind me later"},
		})
		if button.Button == 0 {
			err := updateApp(release)
			if err != nil {
				log.Printf("Unable to update app: %v", err)
			}
		}
	}
}

func checkForRestart() {
	restarting := false
	for _, arg := range os.Args {
		if arg == "-restarting" {
			restarting = true
			break
		}
	}
	if !restarting {
		return
	}
	ppid := syscall.Getppid()
	log.Printf("%d: Detected restart, killing ppid %d", os.Getpid(), ppid)
	syscall.Kill(ppid, syscall.SIGTERM)
}

func checkForNewRelease(githubProject, currentVersion string) *release {
	if currentVersion == "" {
		log.Printf("Not checking updates for dev version")
		return nil
	}
	releases, err := getReleasesFromGitHub(githubProject)
	if err != nil {
		log.Printf("Error fetching github releases: %v", err)
		return nil
	}
	return getReleaseToUpdateTo(releases, currentVersion)
}

func updateApp(release *release) error {
	name, url := downloadURL(release)
	dir, err := ioutil.TempDir("", "menuetupdater")
	if err != nil {
		return fmt.Errorf("Not updating, couldn't get tempdir: %v", err)
	}
	defer os.RemoveAll(dir)
	log.Printf("Downloading archive...")
	archivefile, err := downloadArchive(dir, name, url)
	if err != nil {
		return err
	}
	log.Printf("Unzipping bundle...")
	newAppPath, err := unzipBundle(archivefile)
	if err != nil {
		return err
	}
	return replaceExecutableAndRestart(newAppPath)
}

func replaceExecutableAndRestart(newAppPath string) error {
	currentExecutable, currentAppPath := appPath()
	if currentExecutable == "" {
		log.Fatalf("Cannot update app, can't find executable")
	}
	if currentAppPath == "" {
		log.Fatalf("Cannot update app, not running in Mac app bundle (%s doesn't have /Contents/MacOS)", currentExecutable)
	}
	backupAppPath := currentAppPath + ".updating"
	log.Printf("Updating app (%s to %s)", currentAppPath, newAppPath)
	log.Printf("Move %s to %s", currentAppPath, backupAppPath)
	err := os.Rename(currentAppPath, backupAppPath)
	if err != nil {
		return err
	}
	log.Printf("Move %s to %s", newAppPath, currentAppPath)
	err = os.Rename(newAppPath, currentAppPath)
	if err != nil {
		err := os.Rename(backupAppPath, currentAppPath)
		if err != nil {
			return fmt.Errorf("os.Rename roll back: %v", err)
		}
		return fmt.Errorf("os.Rename move (rollled back): %v", err)
	}
	log.Printf("Removing")
	err = os.RemoveAll(backupAppPath)
	if err != nil {
		return err
	}
	log.Printf("Executing")
	cmd := exec.Command(currentExecutable, "-restarting")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func appPath() (string, string) {
	currentPath, err := os.Executable()
	if err != nil {
		log.Fatalf("os.Executable: %v", err)
	}
	d := strings.Split(currentPath, string(os.PathSeparator))
	if len(d) < 5 || d[len(d)-2] != "MacOS" || d[len(d)-3] != "Contents" {
		return currentPath, ""
	}
	return currentPath, strings.Join(d[0:len(d)-3], string(os.PathSeparator))
}

func getReleasesFromGitHub(project string) ([]release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", project)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	releases := make([]release, 0)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&releases)
	if err != nil {
		return nil, err
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("Could not check for updates: no releases found")
	}
	return releases, nil
}

func downloadURL(release *release) (string, string) {
	name := ""
	url := ""
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			name = asset.Name
			url = asset.DownloadURL
			break
		}
	}
	return name, url
}

func downloadArchive(tempdir, name, url string) (string, error) {
	filename := filepath.Join(tempdir, name)
	out, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("Not updating, couldn't create file in tempdir: %v", err)
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Not updating, couldn't open url: %v", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Not updating, couldn't copy data: %v", err)
	}
	return filename, nil
}

func unzipBundle(filename string) (string, error) {
	destination := filepath.Dir(filename)
	bundle := ""
	r, err := zip.OpenReader(filename)
	if err != nil {
		return "", err
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		fpath := filepath.Join(destination, f.Name)
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				return "", err
			}
			if strings.HasSuffix(f.Name, ".app/") && !strings.Contains(filepath.Dir(f.Name), "/") {
				bundle = fpath
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return "", err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
			if err != nil {
				return "", err
			}
		}
	}
	return bundle, nil
}

func getReleaseToUpdateTo(releases []release, currentVersion string) *release {
	if len(releases) == 0 {
		log.Printf("No github releases found")
		return nil
	}
	found := false
	for ind, release := range releases {
		if release.TagName == currentVersion {
			if ind == 0 {
				log.Printf("Not updating, latest version already running")
				return nil
			}
			found = true
			break
		}
	}
	if !found {
		log.Printf("Our version isn't on the page, not updating")
		return nil
	}
	return &releases[0]
}
