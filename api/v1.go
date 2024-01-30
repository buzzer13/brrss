package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/buzzer13/brrss"
	"gitlab.com/buzzer13/brrss/util"
	"net/http"
	"net/url"
	"strings"
)

// V1GetFeed godoc
//
//	@Description	Generates RSS/Atom feed
//	@Tags			feeds fetch
//	@Accept			json
//	@Produce		application/atom+xml
//	@Produce		application/feed+json
//	@Produce		application/rss+xml
//	@Param			format		path	brrss.FeedFormat	true	"Output feed format"
//	@Param			url			query	string				true	"Source URL"					format(string)
//	@Param			item		query	string				true	"Article selector"				format(string)
//	@Param			feed-title	query	string				false	"Feed title selector"			format(string)
//	@Param			feed-desc	query	string				false	"Feed description selector"		format(string)
//	@Param			item-time	query	string				false	"Article time selector"			format(string)
//	@Param			item-desc	query	string				false	"Article description selector"	format(string)
//	@Param			item-link	query	string				false	"Article link selector"			format(string)
//	@Param			item-title	query	string				false	"Article title selector"		format(string)
//	@Param			req-headers	query	[]string			false	"Outgoing request headers"		collectionFormat(multi)
//	@Router			/v1/feed/{format} [get]
//	@Security		ApiKeyAuth || BasicAuth
func V1GetFeed(ctx echo.Context) error {
	feedFormat := brrss.FeedFormat(ctx.Param("format"))
	feedURL, err := url.Parse(ctx.QueryParam("url"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid feed url - "+err.Error())
	}

	headers := http.Header{}

	for _, reqHeader := range ctx.Request().URL.Query()["req-headers"] {
		hdr := strings.SplitN(reqHeader, ":", 2)

		if len(hdr) != 2 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid header - "+reqHeader)
		}

		headers.Set(strings.TrimSpace(hdr[0]), strings.TrimSpace(hdr[1]))
	}

	res, err := util.Fetch("GET", feedURL.String(), &util.FetchOptions{
		Header: headers,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "http request failure - "+err.Error())
	}

	if res.StatusCode >= 400 {
		return echo.NewHTTPError(res.StatusCode, "remote server error - "+err.Error())
	}

	feed, err := brrss.HTMLToFeed(res.Body, feedFormat, brrss.HTMLToFeedOptions{
		BaseURL:      *feedURL,
		SelItem:      ctx.QueryParam("item"),
		SelFeedTitle: ctx.QueryParam("feed-title"),
		SelFeedDesc:  ctx.QueryParam("feed-desc"),
		SelItemTime:  ctx.QueryParam("item-time"),
		SelItemDesc:  ctx.QueryParam("item-desc"),
		SelItemLink:  ctx.QueryParam("item-link"),
		SelItemTitle: ctx.QueryParam("item-title"),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "feed generation failure - "+err.Error())
	}

	return ctx.String(http.StatusOK, feed)
}
