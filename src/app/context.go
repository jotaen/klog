package app

import (
	"io/ioutil"
	"klog"
	"klog/parser"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Context struct {
	bookmarkedFile string
	config         Config
	homeDir        string
}

func NewContext(config Config, homeDir string) (*Context, error) {
	return &Context{
		config:  config,
		homeDir: homeDir,
	}, nil
}

func NewContextFromEnv() (*Context, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, err
	}
	//config, err := func() (Config, error) {
	//	configToml, err := readFile(homeDir.HomeDir + "/.klog.toml")
	//	if err != nil {
	//		return NewDefaultConfig(), nil
	//	}
	//	return NewConfigFromToml(configToml)
	//}()
	return NewContext(NewDefaultConfig(), homeDir.HomeDir)
}

func (c *Context) HomeDir() string {
	return c.homeDir
}

func (c *Context) RetrieveRecords(paths ...string) ([]klog.Record, error) {
	var records []klog.Record
	for _, p := range paths {
		content, err := readFile(p)
		if err != nil {
			return nil, err
		}
		rs, errs := parser.Parse(content)
		if errs != nil {
			return nil, errs
		}
		records = append(records, rs...)
	}
	return records, nil
}

func (c *Context) WriteRecords(rs []klog.Record, path string) error {
	s := parser.SerialiseRecords(rs, nil)
	return ioutil.WriteFile(path, []byte(s), 0644)
}

func (c *Context) Config() Config {
	return Config{} // TODO
}

type File struct {
	Name     string
	Location string
	Path     string
}

func (c *Context) Bookmark() ([]klog.Record, File, error) {
	bookmarkPath := c.HomeDir() + "/.klog/bookmark.klg"
	dest, _ := os.Readlink(bookmarkPath)
	file := File{
		Name:     filepath.Base(dest),
		Location: filepath.Dir(dest),
		Path:     dest,
	}
	rs, err := c.RetrieveRecords(bookmarkPath)
	return rs, file, err
}

func (c *Context) SetBookmark(path string) error {
	bookmark, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	klogFolder := c.HomeDir() + "/.klog"
	err = os.MkdirAll(klogFolder, 0700)
	if err != nil {
		return err
	}

	symlink := klogFolder + "/bookmark.klg"
	_ = os.Remove(symlink)
	return os.Symlink(bookmark, symlink)
}

func (c *Context) OpenInFileBrowser(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
