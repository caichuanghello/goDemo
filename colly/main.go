package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"os"
	"strings"
	"time"
)



func main(){

	getNovel()
	//getTop250() //爬取豆瓣电影top250
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

//笔趣阁网站小说下载
func getNovel(){

	baseurl := "https://www.xbiquge.la" //笔趣阁官网
	url :=baseurl+"/48/48900/" //需要下载的书的首页连接

	//获取书名
	title:=getTitle(url)
	if len(title) <1 {
		fmt.Println("文件名获取失败...")
		return
	}
	file,err:=os.OpenFile(title+".txt",os.O_RDONLY|os.O_CREATE|os.O_APPEND,0666)
	defer file.Close()
	if err != nil {
		fmt.Println("os.openfile error:",err)
	}
	c:=colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
		colly.MaxDepth(3),//这里表示采集10章
	)

	//获取第一章的url
	c.OnHTML("div[id='list'] dl", func(element *colly.HTMLElement) {
		href,exit:=element.DOM.Find("dd").Find("a").Attr("href")
		if exit {
			getContext(baseurl+href,file,c)
		} else {
			fmt.Println("连接不存在")
		}
	})

	c.Visit(url)

}

//获取小说的标题
func getTitle(url string)(title string){
	c:=colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
		colly.MaxDepth(1),
		colly.Debugger(&debug.LogDebugger{}),
	)
	
	c.OnHTML("div[id='info']", func(element *colly.HTMLElement) {
		//title =element.DOM.ChildrenFiltered("h1").Text()
		//title =element.DOM.Children().Filter("h1").Text() //chidren是后代.find是子后代
		title = element.DOM.Find("h1").Text()

	})

	c.Visit(url)
	return
}

//根据连接获取小数主题内容
func getContext(url string,file *os.File,c *colly.Collector){

	c.OnHTML("div[class='box_con']", func(element *colly.HTMLElement) {
		//获取章节名
		title:=element.DOM.Find("div[class='bookname']").Find("h1").Text()
		n := strings.Index(title,"第")
		if n != -1 {
			title = "    "+string([]byte(title)[n:])+"\r\n\r\n"
			//获取主体内容
			content :=element.DOM.Find("div[id='content']").Text()+"\r\n\r\n\r\n"
			file.Write([]byte(title))
			file.Write([]byte(content))
		}
	})

	//获取下一章内容,然后重复执行
	c.OnHTML("div[class='bottem2']", func(element *colly.HTMLElement) {
		href, found:=element.DOM.Find("a").Eq(3).Attr("href")
		if found {
			err:=element.Request.Visit(element.Request.AbsoluteURL(href)) //这里会拼接得到的重写去调用上面的OnHTML.然后解析页面
			if err !=nil {
				fmt.Println("下一页爬取失败:",err.Error())
				return
			}
		}
	})

	c.Visit(url)
}