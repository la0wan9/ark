package spider

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/la0wan9/ark/internal/adoc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// NewAdocCmd creates a new adoc command
func NewAdocCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "adoc",
		Short: "administrative division of China",
		Run:   adocCmd,
	}
}

func adocCmd(command *cobra.Command, args []string) {
	URL := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"
	responseCallback := func(r *colly.Response) {
		// os.WriteFile("response.html", r.Body, 0666)
		if strings.Contains(string(r.Body), "请开启JavaScript并刷新该页") {
			time.Sleep(1 * time.Minute)
			r.Request.Retry()
		}
	}
	errorCallback := func(r *colly.Response, err error) {
		log.Info(err)
		r.Request.Retry()
	}
	entrance := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.DetectCharset(),
	)
	entrance.Limit(&colly.LimitRule{
		DomainGlob: "*.stats.gov.cn",
		Delay:      100 * time.Millisecond,
	})
	entrance.OnResponse(responseCallback)
	entrance.OnError(errorCallback)
	extensions.RandomUserAgent(entrance)
	extensions.Referer(entrance)
	target := entrance.Clone()
	target.OnResponse(responseCallback)
	target.OnError(errorCallback)
	extensions.RandomUserAgent(target)
	extensions.Referer(target)
	selector := "ul.center_list_contlist > li:first-child > a"
	entrance.OnHTML(selector, func(e *colly.HTMLElement) {
		entrance.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})
	selector = "tr.provincetr a"
	entrance.OnHTML(selector, func(e *colly.HTMLElement) {
		code := strings.TrimSuffix(filepath.Base(e.Attr("href")), ".html")
		if count := 12 - len(code); count > 0 {
			code += strings.Repeat("0", count)
		}
		adoc := &adoc.Adoc{}
		adoc.Name = e.Text
		adoc.Code = cast.ToInt64(code)
		fmt.Println(adoc)
		e.Request.Ctx.Put("parent", adoc.Code)
		target.Request("GET", e.Request.AbsoluteURL(e.Attr("href")), nil, e.Request.Ctx, nil)
	})
	target.OnHTML("tr", func(e *colly.HTMLElement) {
		var href string
		adoc := &adoc.Adoc{}
		if parent, ok := e.Request.Ctx.GetAny("parent").(int64); ok {
			adoc.Parent = parent
		}
		e.ForEach("td", func(i int, e *colly.HTMLElement) {
			if i == 0 {
				href = e.ChildAttr("a", "href")
				adoc.Code = cast.ToInt64(e.Text)
			} else {
				adoc.Name = e.Text
			}
		})
		if adoc.Code == 0 {
			return
		}
		fmt.Println(adoc)
		if href != "" {
			URL := e.Request.AbsoluteURL(href)
			e.Request.Ctx.Put("parent", adoc.Code)
			target.Request("GET", URL, nil, e.Request.Ctx, nil)
		}
	})
	entrance.Visit(URL)
}
