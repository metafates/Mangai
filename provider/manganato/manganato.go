package manganato

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/metafates/mangal/provider/generic"
)

var Config = &generic.Configuration{
	Name:            "Manganato",
	Delay:           50 * time.Millisecond,
	Parallelism:     50,
	ReverseChapters: true,
	BaseURL:         "https://manganato.com/",
	GenerateSearchURL: func(query string) string {
		query = strings.ToLower(query)

		replacements := map[string]string{
			`[àáạảãâầấậẩẫăằắặẳẵ]`: "a",
			`[èéẹẻẽêềếệểễ]`:       "e",
			`[ìíịỉĩ]`:             "i",
			`[òóọỏõôồốộổỗơờớợởỡ]`: "o",
			`[ùúụủũưừứựửữ]`:       "u",
			`[ỳýỵỷỹ]`:             "y",
			`[đ]`:                 "d",
		}

		for pattern, replacement := range replacements {
			query = regexp.MustCompile(pattern).ReplaceAllString(query, replacement)
		}

		query = strings.ReplaceAll(query, " ", "_")
		query = regexp.MustCompile(`[^0-9a-z_]`).ReplaceAllString(query, "_")
		query = regexp.MustCompile(`_+`).ReplaceAllString(query, "_")
		query = strings.Trim(query, "_")
		query = url.QueryEscape(query)

		return fmt.Sprintf("https://manganato.com/search/story/%s", query)
	},
	MangaExtractor: &generic.Extractor{
		Selector: "div.search-story-item",
		Name: func(selection *goquery.Selection) string {
			return strings.TrimSpace(selection.Find("a.item-title").Text())
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find("a.item-title").AttrOr("href", "")
		},
		Cover: func(selection *goquery.Selection) string {
			return selection.Find("img").AttrOr("src", "")
		},
	},
	ChapterExtractor: &generic.Extractor{
		Selector: "li.a-h",
		Name: func(selection *goquery.Selection) string {
			name := selection.Find("a").Text()
			if strings.HasPrefix(name, "Vol.") {
				splitted := strings.Split(name, " ")
				name = strings.Join(splitted[1:], " ")
			}
			return name
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find("a").AttrOr("href", "")
		},
		Volume: func(selection *goquery.Selection) string {
			name := selection.Find("a").Text()
			if strings.HasPrefix(name, "Vol.") {
				splitted := strings.Split(name, " ")
				return splitted[0]
			}
			return ""
		},
	},
	PageExtractor: &generic.Extractor{
		Selector: ".container-chapter-reader img",
		URL: func(selection *goquery.Selection) string {
			return selection.AttrOr("src", "")
		},
	},
}
