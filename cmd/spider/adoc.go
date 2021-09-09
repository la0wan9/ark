package spider

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/la0wan9/ark/internal/adoc"
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
	tmp := filepath.Join(
		os.TempDir(), fmt.Sprintf("%x.html", md5.Sum([]byte(URL))),
	)
	defer os.Remove(tmp)
	responseCallback := func(r *colly.Response) {
		os.WriteFile(tmp, r.Body, 0600)
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
		a := &adoc.Adoc{
			Parent: -1,
		}
		code := strings.TrimSuffix(filepath.Base(e.Attr("href")), ".html")
		if count := 12 - len(code); count > 0 {
			code += strings.Repeat("0", count)
		}
		a.Code = cast.ToInt64(code)
		a.Name = e.Text
		fmt.Println(a)
		if href := e.Attr("href"); href != "" {
			e.Request.Ctx.Put("parent", a.Code)
			target.Request(
				"GET", e.Request.AbsoluteURL(href), nil, e.Request.Ctx, nil,
			)
		}
	})
	target.OnHTML("tr", func(e *colly.HTMLElement) {
		var href string
		a := &adoc.Adoc{}
		if parent, ok := e.Request.Ctx.GetAny("parent").(int64); ok {
			a.Parent = parent
		}
		e.ForEach("td", func(i int, e *colly.HTMLElement) {
			if i == 0 {
				href = e.ChildAttr("a", "href")
				a.Code = cast.ToInt64(e.Text)
			} else {
				a.Name = e.Text
			}
		})
		if a.Code == 0 {
			return
		}
		fmt.Println(a)
		if href != "" {
			e.Request.Ctx.Put("parent", a.Code)
			target.Request(
				"GET", e.Request.AbsoluteURL(href), nil, e.Request.Ctx, nil,
			)
		}
	})
	entrance.Visit(URL)
}
