package installer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Scraper struct {
	Name        string
	URL         string
	Description string
	Contents    string
}

func (s *Scraper) Path() string {
	filename := fmt.Sprintf("%s.lua", s.Name)
	return filepath.Join(where.Sources(), filename)
}

func (s *Scraper) GithubURL() string {
	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/scrapers/%s.lua", collector.user, collector.repo, collector.branch, s.Name)
}

func (s *Scraper) download() error {
	if s.Contents != "" {
		return nil
	}

	if s.URL == "" {
		return fmt.Errorf("url must be set")
	}

	res, err := http.Get(s.URL)
	if err != nil {
		return err
	}

	defer util.Ignore(res.Body.Close)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get %s: %s", s.URL, res.Status)
	}

	var b []byte
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var info = struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}{}

	err = json.Unmarshal(b, &info)
	if err != nil {
		return err
	}

	switch info.Encoding {
	case "base64":
		text, err := base64.StdEncoding.DecodeString(info.Content)
		if err != nil {
			return err
		}

		s.Contents = string(text)
	default:
		return fmt.Errorf("unsupported encoding: %s", info.Encoding)
	}

	return nil
}

func (s *Scraper) Install() error {
	err := s.download()

	if err != nil {
		return err
	}

	return filesystem.Api().WriteFile(s.Path(), []byte(s.Contents), os.ModePerm)
}
