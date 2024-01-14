package brrss

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"io"
	"net/url"
	"strings"
	"time"
)

type FeedFormat string
type SelKind int64

type SelectTextOptions struct {
	BaseURL url.URL
}

type HTMLToFeedOptions struct {
	BaseURL      url.URL
	SelItem      string
	SelFeedTitle string
	SelFeedDesc  string
	SelItemTime  string
	SelItemDesc  string
	SelItemLink  string
	SelItemTitle string
}

const (
	AtomFeedFormat FeedFormat = "atom"
	JSONFeedFormat FeedFormat = "json"
	RSSFeedFormat  FeedFormat = "rss"
)

const (
	FeedTitleSelKind SelKind = iota
	FeedDescSelKind
	ItemTitleSelKind
	ItemLinkSelKind
	ItemDescSelKind
	ItemTimeSelKind
)

var selDefaults = map[SelKind][]string{
	FeedTitleSelKind: {
		"meta[property='og:title']!attr:content",
	},
	FeedDescSelKind: {
		"meta[property='og:description']!attr:content",
	},
	ItemTitleSelKind: {
		"h1!text", "h2!text", "h3!text", "h4!text", "h5!text", "h6!text",
		"header!text",
	},
	ItemLinkSelKind: {
		"&!attr:href!link", "a!attr:href!link",
	},
	ItemDescSelKind: {
		".description!text", ".desc!text",
		"p!text",
	},
	ItemTimeSelKind: {},
}

func HTMLToFeed(body io.Reader, feedFormat FeedFormat, options HTMLToFeedOptions) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	stOptions := &SelectTextOptions{
		BaseURL: options.BaseURL,
	}

	if err != nil {
		return "", errors.New("page parse failed")
	}

	now := time.Now()
	items := make([]*feeds.Item, 0)

	doc.Find(options.SelItem).Each(func(i int, item *goquery.Selection) {
		items = append(items, &feeds.Item{
			Title:       TrySelectText(item, ItemTitleSelKind, options.SelItemTitle, stOptions),
			Link:        &feeds.Link{Href: TrySelectText(item, ItemLinkSelKind, options.SelItemLink, stOptions)},
			Description: TrySelectText(item, ItemDescSelKind, options.SelItemDesc, stOptions),
			Created:     TryParseTime(time.RFC3339, TrySelectText(item, ItemTimeSelKind, options.SelItemTime, stOptions)),
		})
	})

	if len(items) <= 0 {
		return "", errors.New("no articles found")
	}

	feed := &feeds.Feed{
		Title:       TrySelectText(doc.Selection, FeedTitleSelKind, options.SelFeedTitle, stOptions),
		Link:        &feeds.Link{Href: options.BaseURL.String()},
		Description: TrySelectText(doc.Selection, FeedDescSelKind, options.SelFeedDesc, stOptions),
		Created:     now,
		Items:       items,
	}

	switch feedFormat {
	case AtomFeedFormat:
		return feed.ToAtom()
	case JSONFeedFormat:
		return feed.ToJSON()
	case RSSFeedFormat:
		return feed.ToRss()
	default:
		return "", errors.New("unsupported feed format")
	}
}

func SelectText(item *goquery.Selection, selKind SelKind, selCustom string, options *SelectTextOptions) (string, error) {
	selectors := append([]string{selCustom}, selDefaults[selKind]...)

	if item == nil {
		return "", errors.New("item is nil")
	}

	if item.Length() <= 0 {
		return "", errors.New("item is empty")
	}

	for _, sel := range selectors {
		actions := strings.Split(sel, "!")
		selItem := item
		result := ""

		if actions[0] != "&" {
			selItem = item.Find(actions[0])
		}

		if selItem.Length() <= 0 {
			continue
		}

		for _, act := range actions[1:] {
			args := strings.Split(act, ":")

			switch fmt.Sprintf("%s:%d", args[0], len(args)-1) {
			// Data extraction actions
			case "attr:1":
				result = selItem.AttrOr(args[1], "")
			case "attr:2":
				result = selItem.AttrOr(args[1], args[2])
			case "text:0":
				result = selItem.Text()

			// Text processing actions
			case "link:0":
				articleURL, err := url.Parse(result)

				if err == nil {
					newURL := options.BaseURL
					newURL.Path = articleURL.Path
					newURL.RawQuery = articleURL.RawQuery

					result = newURL.String()
				}
			case "time:1":
				dt, err := time.Parse(args[1], result)

				if err == nil {
					result = dt.UTC().Format(time.RFC3339)
				}

			default:
				return "", errors.New("unsupported action or invalid arguments")
			}

			if result == "" {
				break
			}
		}

		if result != "" {
			return result, nil
		}
	}

	return "", errors.New("no selector match")
}

func TrySelectText(item *goquery.Selection, selKind SelKind, selCustom string, options *SelectTextOptions) string {
	value, _ := SelectText(item, selKind, selCustom, options)

	return value
}

func TryParseTime(layout string, value string) time.Time {
	dt, err := time.Parse(layout, value)

	if err != nil {
		return time.Now()
	}

	return dt
}
