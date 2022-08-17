package downloader

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
)

func Download(src source.Source, chapter *source.Chapter, progress func(string)) (string, error) {
	log.Info("downloading " + chapter.Name)
	progress("Getting pages")
	pages, err := src.PagesOf(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Info("found " + fmt.Sprintf("%d", len(pages)) + " pages")

	log.Info("downloading " + fmt.Sprintf("%d", len(pages)) + " pages")
	progress(fmt.Sprintf("Downloading %d pages", len(pages)))
	err = chapter.DownloadPages()
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Info("gettings " + viper.GetString(config.FormatsUse) + " converter")
	progress(fmt.Sprintf("Converting %d pages to %s", len(pages), viper.GetString(config.FormatsUse)))
	conv, err := converter.Get(viper.GetString(config.FormatsUse))
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Info("converting " + viper.GetString(config.FormatsUse))
	path, err := conv.Save(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if viper.GetBool(config.HistorySaveOnDownload) {
		go func() {
			log.Info("saving history")
			err = history.Save(chapter)
			if err != nil {
				log.Warn(err)
			} else {
				log.Info("history saved")
			}
		}()
	}

	log.Info("downloaded without errors")
	progress("Downloaded")
	return path, nil
}
