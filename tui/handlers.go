package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/installer"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"strings"
)

func (b *statefulBubble) loadScrapers() tea.Cmd {
	return func() tea.Msg {
		b.progressStatus = "Loading scrapers"
		scrapers, err := installer.Scrapers()
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}
		b.progressStatus = "Scrapers Loaded"

		slices.SortFunc(scrapers, func(a, b *installer.Scraper) bool {
			return strings.Compare(a.Name, b.Name) < 0
		})

		var items []list.Item
		for _, s := range scrapers {
			items = append(items, &listItem{
				title:       s.Name,
				description: s.GithubURL(),
				internal:    s,
			})
		}

		cmd := b.scrapersInstallC.SetItems(items)
		b.scrapersLoadedChannel <- scrapers
		return cmd
	}
}

func (b *statefulBubble) waitForScrapersLoaded() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.scrapersLoadedChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) installScraper(s *installer.Scraper) tea.Cmd {
	return func() tea.Msg {
		b.progressStatus = fmt.Sprintf("Installing %s", s.Name)
		err := s.Install()
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("scraper " + s.Name + " installed")
			b.scraperInstalledChannel <- s
		}

		return nil
	}
}

func (b *statefulBubble) waitForScraperInstallation() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.scraperInstalledChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) loadSource(p *provider.Provider) tea.Cmd {
	return func() tea.Msg {
		log.Info("loading source " + p.ID)
		b.progressStatus = "Initializing source"
		s, err := p.CreateSource()

		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("source " + p.ID + " loaded")
			b.sourceLoadedChannel <- s
		}

		return nil
	}
}

func (b *statefulBubble) waitForSourceLoaded() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.sourceLoadedChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) searchManga(query string) tea.Cmd {
	return func() tea.Msg {
		log.Info("searching for " + query)
		b.progressStatus = "Searching"
		mangas, err := b.selectedSource.Search(query)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("found " + fmt.Sprintf("%d", len(mangas)) + " mangas")
			b.foundMangasChannel <- mangas
		}

		return nil
	}
}

func (b *statefulBubble) waitForMangas() tea.Cmd {
	return func() tea.Msg {
		select {
		case found := <-b.foundMangasChannel:
			return found
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) getChapters(manga *source.Manga) tea.Cmd {
	return func() tea.Msg {
		log.Info("getting chapters of " + manga.Name)
		chapters, err := b.selectedSource.ChaptersOf(manga)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("found " + fmt.Sprintf("%d", len(chapters)) + " chapters")
			b.foundChaptersChannel <- chapters
		}

		return nil
	}
}

func (b *statefulBubble) waitForChapters() tea.Cmd {
	return func() tea.Msg {
		select {
		case found := <-b.foundChaptersChannel:
			return found
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) readChapter(chapter *source.Chapter) tea.Cmd {
	return func() tea.Msg {
		defer func() {
			if viper.GetBool(config.HistorySaveOnRead) {
				go func() {
					log.Info("saving history")
					err := history.Save(chapter)
					if err != nil {
						log.Warn(err)
					} else {
						log.Info("history saved")
					}
				}()
			}
		}()

		if viper.GetBool(config.ReaderReadInBrowser) {
			err := open.Start(chapter.URL)
			if err != nil {
				b.errorChannel <- err
			} else {
				b.chapterReadChannel <- struct{}{}
			}

			return nil
		}

		log.Info("downloading " + chapter.Name + " from " + chapter.Manga.Name + " for reading. Provider is " + b.selectedSource.ID())
		b.progressStatus = "Gettings pages"
		b.currentDownloadingChapter = chapter
		log.Info("getting pages of " + chapter.Name)
		pages, err := b.selectedSource.PagesOf(chapter)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}

		log.Info("downloading " + fmt.Sprintf("%d", len(pages)) + " pages")
		b.progressStatus = fmt.Sprintf("Downloading %d pages", len(pages))
		err = chapter.DownloadPages()
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}

		log.Info("getting " + viper.GetString(config.FormatsUse) + " converter")
		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}

		log.Info("converting " + viper.GetString(config.FormatsUse))
		b.progressStatus = fmt.Sprintf("Converting %d pages to %s", len(pages), viper.GetString(config.FormatsUse))
		path, err := conv.SaveTemp(chapter)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}

		log.Info("downloaded without errors. Opening " + path)
		if reader := viper.GetString(config.ReaderName); reader != "" {
			log.Info("opening with " + reader)
			b.progressStatus = fmt.Sprintf("Opening %s", reader)
			err = open.RunWith(reader, path)
			if err != nil {
				log.Error(err)
				b.errorChannel <- err
				return nil
			}
			log.Info("opened without errors")
		} else {
			log.Info("no reader specified. opening with default")
			b.progressStatus = "Opening"
			err = open.Run(path)
			if err != nil {
				log.Error(err)
				b.errorChannel <- err
				return nil
			}
			log.Info("opened without errors")
		}

		b.progressStatus = "Done"
		b.chapterReadChannel <- struct{}{}

		return nil
	}
}

func (b *statefulBubble) waitForChapterRead() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.chapterReadChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) downloadChapter(chapter *source.Chapter) tea.Cmd {
	return func() tea.Msg {
		log.Info("downloading " + chapter.Name)
		b.currentDownloadingChapter = chapter
		b.progressStatus = "Getting pages"
		pages, err := b.selectedSource.PagesOf(chapter)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}
		log.Info("found " + fmt.Sprintf("%d", len(pages)) + " pages")

		log.Info("downloading " + fmt.Sprintf("%d", len(pages)) + " pages")
		b.progressStatus = fmt.Sprintf("Downloading %d pages", len(pages))
		err = chapter.DownloadPages()
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		log.Info("gettings " + viper.GetString(config.FormatsUse) + " converter")
		b.progressStatus = fmt.Sprintf("Converting %d pages to %s", len(pages), viper.GetString(config.FormatsUse))
		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}

		log.Info("converting " + viper.GetString(config.FormatsUse))
		path, err := conv.Save(chapter)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
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
		b.progressStatus = "Downloaded"
		b.chapterDownloadChannel <- struct{}{}
		b.lastDownloadedChapterPath = path

		return nil
	}
}

func (b *statefulBubble) waitForChapterDownload() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.chapterDownloadChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}
