package mangadex

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strconv"
)

func (m *Mangadex) Search(query string) ([]*source.Manga, error) {
	if cached, ok := m.cache.mangas.Get(query).Get(); ok {
		for _, manga := range cached {
			manga.Source = m
		}

		return cached, nil
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(100))

	ratings := []string{mangodex.Safe, mangodex.Suggestive}

	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	if viper.GetBool(key.MangadexNSFW) {
		params.Add("contentRating[]", mangodex.Porn)
		params.Add("contentRating[]", mangodex.Erotica)
	}

	params.Set("order[relevance]", "desc")
	params.Set("title", query)

	mangaList, err := m.client.Manga.GetMangaList(params)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var mangas []*source.Manga

	for i, manga := range mangaList.Data {
		m := source.Manga{
			Name:   manga.GetTitle(viper.GetString(key.MangadexLanguage)),
			URL:    fmt.Sprintf("https://mangadex.org/title/%s", manga.ID),
			Index:  uint16(i),
			ID:     manga.ID,
			Source: m,
		}

		mangas = append(mangas, &m)
	}

	_ = m.cache.mangas.Set(query, mangas)
	return mangas, nil
}
