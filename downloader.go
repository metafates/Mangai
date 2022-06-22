package main

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// RemoveIfExists removes file if it exists
func RemoveIfExists(path string) error {
	exists, err := Afero.Exists(path)

	if err != nil {
		return err
	}

	if exists {
		err = Afero.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveTemp saves file to OS temp dir and returns its path
// It's a caller responsibility to remove created file
func SaveTemp(contents *[]byte) (string, error) {
	out, err := Afero.TempFile("", TempPrefix+"*")

	if err != nil {
		return "", err
	}

	defer func(out afero.File) {
		err := out.Close()
		if err != nil {
			log.Fatal("Unexpected error while closing file")
		}
	}(out)

	_, err = out.Write(*contents)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

type DownloaderStage int

const (
	Scraping DownloaderStage = iota
	Downloading
	Converting
	Cleanup
	Done
)

type ChaptersDownloadProgress struct {
	Current   *URL
	Done      bool
	Failed    []*URL
	Succeeded []string
	Total     int
	Proceeded int
}

type ChapterDownloadProgress struct {
	Stage   DownloaderStage
	Message string
}

// DownloadChapter downloads chapter from the given url and returns its path
func DownloadChapter(chapter *URL, progress chan ChapterDownloadProgress, temp bool) (string, error) {
	mangaTitle := chapter.Relation.Info
	var (
		mangaPath string
		err       error
	)

	// Get future path to manga
	if temp {
		mangaPath = os.TempDir()
	} else {
		absPath, err := filepath.Abs(UserConfig.Path)

		if err != nil {
			return "", err
		}

		mangaPath = filepath.Join(absPath, mangaTitle)
	}

	showProgress := progress != nil

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Scraping,
			Message: "Getting pages",
		}
	}

	var chapterPath string
	// Get future path to chapter
	if temp {
		chapterPath = filepath.Join(mangaPath, fmt.Sprintf(TempPrefix+" [%d] %s", chapter.Index, chapter.Info))
	} else {
		chapterPath = filepath.Join(mangaPath, fmt.Sprintf("[%d] %s", chapter.Index, chapter.Info))
	}

	// Replace whitespaces with underscore and colon with unicode colon
	// Windows is very bad at escaping whitespaces, so have to use this workaround
	if runtime.GOOS == "windows" {
		chapterPath = strings.ReplaceAll(chapterPath, " ", "-")

		volumeName := filepath.VolumeName(chapterPath)
		if strings.Contains(volumeName, ":") {
			chapterPath = strings.ReplaceAll(chapterPath, ":", "꞉") // Unicode U+A789
			chapterPath = strings.Replace(chapterPath /* unicode colon */, "꞉" /* ascii colon */, ":", 1)
		} else {
			chapterPath = strings.ReplaceAll(chapterPath, ":", "꞉") // Unicode U+A789
		}
	}

	// Get chapter contents
	pages, err := chapter.Scraper.GetPages(chapter)
	pagesCount := len(pages)

	if err != nil {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Downloading,
			Message: fmt.Sprintf("Downloading %d pages", pagesCount),
		}
	}

	var (
		tempPaths        = make([]string, pagesCount)
		wg               sync.WaitGroup
		errorEncountered bool
	)

	wg.Add(pagesCount)

	// Download pages in parallel
	for _, page := range pages {
		go func(p *URL) {
			defer wg.Done()
			var (
				data     *[]byte
				tempPath string
			)

			data, err = chapter.Scraper.GetFile(p)

			if err != nil {
				// TODO: use channel
				errorEncountered = true
				return
			}

			tempPath, err = SaveTemp(data)
			fixedTempPath := tempPath + ".jpg"
			err = Afero.Rename(tempPath, fixedTempPath)
			tempPaths[p.Index] = fixedTempPath
		}(page)
	}

	wg.Wait()

	defer chapter.Scraper.ResetFiles()

	if errorEncountered {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Converting,
			Message: fmt.Sprintf("Converting %d pages to %s", pagesCount, UserConfig.Format),
		}
	}

	if len(tempPaths) == 0 {
		return "", errors.New("pages was not downloaded")
	}

	// Convert pages to desired format
	chapterPath, err = Packers[UserConfig.Format](tempPaths, chapterPath)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Cleanup,
			Message: "Removing temp files",
		}
	}

	if err != nil {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Done,
			Message: fmt.Sprintf("Chapter %s downloaded", chapter.Info),
		}
	}

	return chapterPath, nil
}
