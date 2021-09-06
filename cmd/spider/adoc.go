package spider

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/spf13/cobra"

	"github.com/la0wan9/ark/internal/adoc"
)

const url = "http://www.mca.gov.cn/article/sj/xzqh/1980/"

// NewAdocCmd creates a new adoc command
func NewAdocCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "adoc",
		Short: "administrative division of China",
		Run:   adocCmd,
	}
}

func adocCmd(command *cobra.Command, args []string) {
	entrance := colly.NewCollector()
	_ = entrance.Limit(&colly.LimitRule{
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	extensions.RandomUserAgent(entrance)
	extensions.Referer(entrance)
	target := entrance.Clone()
	selector := "div#list_content tr:nth-child(1) > td.arlisttd > a"
	entrance.OnHTML(selector, func(e *colly.HTMLElement) {
		err := entrance.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if err != nil {
			panic(err)
		}
	})
	selector = "div#zoom p:nth-child(2) > a"
	entrance.OnHTML(selector, func(e *colly.HTMLElement) {
		err := target.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if err != nil {
			panic(err)
		}
	})
	var adocs []*adoc.Adoc
	parents := make(map[int]int64)
	target.OnHTML("tr", func(e *colly.HTMLElement) {
		adoc := &adoc.Adoc{}
		e.ForEachWithBreak("td", func(_ int, e *colly.HTMLElement) bool {
			if adoc.Code == 0 {
				code, _ := strconv.ParseInt(e.Text, 10, 64)
				if code == 0 {
					return true
				}
				level := 0
				switch {
				case code%10000 == 0:
				case code%100 == 0:
					level = 1
				default:
					level = 2
				}
				parents[level] = code
				parent := parents[level-1]
				if parent == 0 && level == 2 {
					parent = parents[0]
				}
				adoc.Code = code
				adoc.Parent = parent
				return true
			}
			if adoc.Name == "" {
				adoc.Name = strings.TrimSpace(e.Text)
				return true
			}
			return false
		})
		if adoc.Code > 0 {
			adocs = append(adocs, adoc)
		}
	})
	if err := entrance.Visit(url); err != nil {
		panic(err)
	}
	if isJSON, _ := command.Flags().GetBool("json"); isJSON {
		out, _ := json.MarshalIndent(adocs, "", "  ")
		fmt.Println(string(out))
	} else if isXML, _ := command.Flags().GetBool("xml"); isXML {
		out, _ := xml.MarshalIndent(adocs, "", "  ")
		fmt.Println(xml.Header + string(out))
	} else {
		for _, adoc := range adocs {
			fmt.Println(adoc)
		}
	}
}
