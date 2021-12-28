package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strings"
	"time"
)

type Hash [32]byte

//爬取豆瓣电影top250
func main(){

	getTop250()
}


type Record struct {
	Username string `row_name:"用户名"`
	Content string `row_name:"评论内容"`
}


func getTop250(){
	fName := "douban_movie_top250.csv"
	file, _ := os.Create(fName)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"排行", "影名", "评分","简介","链接"})

	c:=colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
		//colly.MaxDepth(1),
		//colly.Debugger(&debug.LogDebugger{}),
	)

	var url = "https://movie.douban.com/top250"

	c.OnHTML("ol[class='grid_view']", func(element *colly.HTMLElement) {
		element.ForEach("li", func(i int, item *colly.HTMLElement) {
			//获取电影链接
			detailUrl := item.ChildAttr("div > div[class='info'] > div[class='hd'] > a","href")
			getDetail(c,detailUrl,writer)
		})
	})

	c.OnHTML("div.paginator > span.next", func(element *colly.HTMLElement) {
		href, found := element.DOM.Find("a").Attr("href")
		// 如果有下一页，则继续访问
		if found {
			err:=element.Request.Visit(element.Request.AbsoluteURL(href)) //这里会拼接得到的重写去调用上面的OnHTML.然后解析页面
			if err !=nil {
				fmt.Println("下一页爬取失败:",err.Error())
			}
		}
	})
	c.Visit(url)

}

func getDetail(c *colly.Collector, url string,w *csv.Writer){
	collector := c.Clone()

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 2 * time.Second,
	})

	collector.OnHTML("body", func(element *colly.HTMLElement) {
		selection := element.DOM.Find("div#content")
		//排行
		idx:=selection.Find("div.top250 > span").First().Text()
		//标题
		title:=selection.Find("h1 > span").First().Text()
		//简介
		content:= selection.Find("div#link-report  span[property='v:summary']").Text()
		content = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(content, " ", ""),"\n",""))
		//评分
		mark :=selection.Find("div[class='grid-16-8 clearfix'] div#interest_sectl div[class='rating_self clearfix'] > strong").Text()
		fmt.Println("排行:",idx)
		fmt.Println("电影:",title)
		fmt.Println("评分:",mark)
		fmt.Println("简介:", content)
		fmt.Println("链接:",url)
		fmt.Println()

		w.Write([]string{idx,title,mark,content,url})
	})

	collector.Visit(url)

}
