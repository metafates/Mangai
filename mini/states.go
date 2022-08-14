package mini

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"regexp"
	"strconv"
	"strings"
)

type state int

const (
	mangasSearchState state = iota + 1
	mangaSelectState
	sourceSelectState
	chapterSelectState
	chapterReadState
	chaptersDownloadState
	historySelectState
	quitState
)

func (m *mini) handleSourceSelectState() error {
	defaultProviders := provider.DefaultProviders()
	customProviders := lo.Must(provider.CustomProviders())

	var providers = make([]*provider.Provider, 0)

	for _, p := range defaultProviders {
		providers = append(providers, p)
	}

	for _, p := range customProviders {
		providers = append(providers, p)
	}

	slices.SortFunc(providers, func(a *provider.Provider, b *provider.Provider) bool {
		return strings.Compare(a.String(), b.String()) < 0
	})

	var err error

	title("Select Source")
	b, p, err := menu(providers)
	if err != nil {
		return err
	}

	if quit.eq(b) {
		m.newState(quitState)
		return nil
	}

	erase := progress("Initializing Source..")
	m.selectedSource, err = p.CreateSource()
	erase()

	m.newState(mangasSearchState)
	return err
}

func (m *mini) handleMangaSearchState() error {
	var searchLoop func() error
	title("Search Manga")

	searchLoop = func() error {
		in, err := getInput(func(s string) bool {
			return s != ""
		})

		if err != nil {
			return err
		}

		query := in.value

		erase := progress("Searching Query..")
		m.cachedMangas[query], err = m.selectedSource.Search(query)
		max := lo.Min([]int{len(m.cachedMangas[query]), viper.GetInt(config.MiniSearchLimit)})
		m.cachedMangas[query] = m.cachedMangas[query][:max]
		erase()

		if len(m.cachedMangas[query]) == 0 {
			fail("No search results found")
			return searchLoop()
		}

		m.query = query
		m.newState(mangaSelectState)
		return err
	}

	return searchLoop()
}

func (m *mini) handleMangaSelectState() error {
	var err error
	title("Query Results >>")
	b, p, err := menu(m.cachedMangas[m.query])
	if err != nil {
		return err
	}

	if quit.eq(b) {
		m.newState(quitState)
		return nil
	}

	m.selectedManga = p
	m.newState(chapterSelectState)
	return err
}

func (m *mini) handleChapterSelectState() error {
	var err error

	erase := progress("Searching Chapters..")
	m.cachedChapters[m.selectedManga.URL], err = m.selectedSource.ChaptersOf(m.selectedManga)
	erase()
	if err != nil {
		return err
	}

	chapters := m.cachedChapters[m.selectedManga.URL]

	if len(chapters) == 0 {
		fail("No chapters found")
		m.selectedManga = nil
		m.newState(mangaSelectState)
		return nil
	}

	title(fmt.Sprintf("To specify a range, use: start_number end_number (Chapters: 1-%d)", len(chapters)))
	oneChapterInput := regexp.MustCompile(`^\d+$`)
	rangeInput := regexp.MustCompile(`^\d+ \d+$`)
	in, err := getInput(func(s string) bool {
		var err error

		switch {
		case rangeInput.MatchString(s):
			var a, b int64
			{
				l := strings.Split(s, " ")
				a, err = strconv.ParseInt(l[0], 10, 16)
				if err != nil {
					return false
				}

				b, err = strconv.ParseInt(l[1], 10, 16)
				if err != nil {
					return false
				}
			}

			return a < b && 0 < a && int(a) < len(chapters) && int(b) <= len(chapters)
		case oneChapterInput.MatchString(s):
			var a int64
			a, err = strconv.ParseInt(s, 10, 16)
			if err != nil {
				return false
			}

			return 0 < a && int(a) <= len(chapters)
		default:
			return s == "q"
		}
	})

	if err != nil {
		return err
	}

	switch {
	case rangeInput.MatchString(in.value):
		nums := strings.Split(in.value, " ")
		from := lo.Must(strconv.ParseInt(nums[0], 10, 16))
		to := lo.Must(strconv.ParseInt(nums[1], 10, 16))

		for i := from - 1; i < to; i++ {
			m.selectedChapters = append(m.selectedChapters, chapters[i])
		}
	case oneChapterInput.MatchString(in.value):
		num := lo.Must(strconv.ParseInt(in.value, 10, 16))
		m.selectedChapters = append(m.selectedChapters, chapters[num-1])
	case in.value == "q":
		m.newState(quitState)
		return nil
	}

	if m.download {
		m.newState(chaptersDownloadState)
	} else {
		m.newState(chapterReadState)
	}

	return nil
}

func (m *mini) handleChapterReadState() error {
	type controls struct {
		next chan struct{}
		prev chan struct{}
		stop chan struct{}
		err  chan error
	}

	var (
		err      error
		readLoop func(*source.Chapter, *controls, bool, bool)
	)

	readLoop = func(chapter *source.Chapter, c *controls, hasPrev, hasNext bool) {
		util.ClearScreen()
		erase := progress("Loading Chapter..")
		m.cachedPages[chapter.URL], err = m.selectedSource.PagesOf(chapter)
		erase()
		if err != nil {
			c.err <- err
			return
		}

		erase = progress("Downloading Pages..")
		err = chapter.DownloadPages()
		erase()

		if err != nil {
			c.err <- err
			return
		}

		erase = progress("Converting..")
		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			c.err <- err
			return
		}

		path, err := conv.SaveTemp(chapter)
		go func(chapter *source.Chapter) {
			if viper.GetBool(config.HistorySaveOnRead) {
				_ = history.Save(chapter)
			}
		}(chapter)
		erase()

		if err != nil {
			c.err <- err
			return
		}

		erase = progress("Opening..")

		if reader := viper.GetString(config.ReaderName); reader != "" {
			err = open.RunWith(path, reader)
			if err != nil {
				c.err <- err
				return
			}
		} else {
			err = open.Run(path)
			if err != nil {
				c.err <- err
				return
			}
		}

		erase()

		title(fmt.Sprintf("Currently reading %s", chapter.Name))

		var options []*bind
		if hasPrev {
			options = append(options, prev)
		}
		if hasNext {
			options = append(options, next)
		}

		options = append(options, reread, back, search)

		b, _, err := menu([]fmt.Stringer{}, options...)
		if err != nil {
			c.err <- err
			return
		}

		switch b {
		case next:
			c.next <- struct{}{}
		case reread:
			readLoop(chapter, c, hasPrev, hasNext)
		case back:
			m.previousState()
			c.stop <- struct{}{}
		case search:
			m.newState(mangasSearchState)
			c.stop <- struct{}{}
		case quit:
			m.newState(quitState)
			c.stop <- struct{}{}
		}
	}

	c := &controls{
		next: make(chan struct{}),
		prev: make(chan struct{}),
		stop: make(chan struct{}),
		err:  make(chan error),
	}

	var i int

	for {
		var (
			hasPrev = i > 0
			hasNext = i+1 < len(m.selectedChapters)
		)

		go readLoop(m.selectedChapters[i], c, hasPrev, hasNext)

		select {
		case <-c.next:
			i++
		case <-c.prev:
			i--
		case <-c.stop:
			return nil
		case err := <-c.err:
			return err
		}
	}
}

func (m *mini) handleChaptersDownloadState() error {
	var (
		err          error
		downloadLoop func(*source.Chapter) error
	)

	downloadLoop = func(chapter *source.Chapter) error {
		util.ClearScreen()

		erase := progress("Fetching pages links..")
		m.cachedPages[chapter.URL], err = m.selectedSource.PagesOf(chapter)
		erase()
		if err != nil {
			return err
		}

		title(fmt.Sprintf("Currently downloading %s %s (%s)", chapter.Manga.Name, chapter.Name, m.selectedSource.Name()))

		erase = progress(fmt.Sprintf("Downloading %d Pages..", len(m.cachedPages[chapter.URL])))
		err = chapter.DownloadPages()
		erase()

		if err != nil {
			return err
		}

		erase = progress("Converting..")
		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			return err
		}

		_, err = conv.Save(chapter)
		go func(chapter *source.Chapter) {
			if viper.GetBool(config.HistorySaveOnDownload) {
				_ = history.Save(chapter)
			}
		}(chapter)
		erase()

		if err != nil {
			return err
		}

		return nil
	}

	for _, chapter := range m.selectedChapters {
		err = downloadLoop(chapter)
		if err != nil {
			return err
		}
	}

	util.ClearScreen()
	title(fmt.Sprintf("%s downloaded.", util.Quantity(len(m.selectedChapters), "chapter")))
	b, _, err := menu([]fmt.Stringer{}, back, search)
	if err != nil {
		return err
	}

	switch b {
	case back:
		m.previousState()
	case search:
		m.newState(mangasSearchState)
	case quit:
		m.newState(quitState)
	}

	return nil
}

func (m *mini) handleHistorySelectState() error {
	h, err := history.Get()
	if err != nil {
		return err
	}

	chapters := lo.Values(h)

	title("History Results >>")
	b, c, err := menu(chapters)
	if err != nil {
		return err
	}

	switch b {
	case quit:
		m.newState(quitState)
		return nil
	}

	defaultProviders := provider.DefaultProviders()
	customProviders, _ := provider.CustomProviders()

	var providers = make([]*provider.Provider, 0)

	for _, p := range defaultProviders {
		providers = append(providers, p)
	}

	for _, p := range customProviders {
		providers = append(providers, p)
	}

	p, _ := lo.Find(providers, func(p *provider.Provider) bool {
		return p.ID == c.SourceID
	})

	erase := progress("Initializing Source..")
	s, err := p.CreateSource()
	if err != nil {
		return err
	}
	m.selectedSource = s
	erase()

	erase = progress("Fetching Chapters..")
	manga := &source.Manga{
		Name:     c.MangaName,
		URL:      c.MangaURL,
		Index:    0,
		SourceID: c.SourceID,
		ID:       c.MangaID,
	}
	chaps, err := m.selectedSource.ChaptersOf(manga)
	erase()

	if err != nil {
		return err
	}

	m.cachedChapters[manga.URL] = chaps
	m.selectedChapters = chaps[c.Index-1:]

	m.newState(chapterReadState)
	return nil
}
