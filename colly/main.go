package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"os"
	"strings"
	"time"
)



func main(){

	//getNovel() //笔趣阁小说下载
	//getTop250() //爬取豆瓣电影top250

	//getVideoList()

	fmt.Println("没运行吗")

	url:="https://upos-sz-mirrorcoso1.bilivideo.com/upgcxcode/52/92/260129252/260129252-1-30033.m4s?e=ig8euxZM2rNcNbdlhoNvNC8BqJIzNbfqXBvEqxTEto8BTrNvN0GvT90W5JZMkX_YN0MvXg8gNEV4NC8xNEV4N03eN0B5tZlqNxTEto8BTrNvNeZVuJ10Kj_g2UB02J0mN0B5tZlqNCNEto8BTrNvNC7MTX502C8f2jmMQJ6mqF2fka1mqx6gqj0eN0B599M=&uipk=5&nbs=1&deadline=1641283387&gen=playurlv2&os=coso1bv&oi=1941969994&trid=beaa4b83939b4c1eb757e728f3cd4d9au&platform=pc&upsig=6dfa7c1e5fb928dbfadc576f88615206&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,platform&mid=0&bvc=vod&nettype=0&orderid=0,3&agrr=1&bw=20560&logo=80000000"
	c :=colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
	)
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("绝了")
		//r.Headers.Set("content-Type","application/json; charset=utf-8")
		file,err:=os.OpenFile("1.mp4",os.O_CREATE|os.O_WRONLY,0666)
		if err != nil {
			fmt.Println("OpenFile error :",err.Error())
			return
		}
		n,err:=file.Write(r.Body)
		if err != nil {
			fmt.Println("视频数据采集错误!")
			return
		}
		fmt.Println("文件写入:",n)

	})

	err:=c.Visit(url)

	if err != nil {
		fmt.Println("visti error :",err.Error())
	}

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
		colly.MaxDepth(3),//这里表示采集3章
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
		//title =element.DOM.Children().Filter("h1").Text() //Children是后代.find是子后代
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
			content :=element.DOM.Find("div[id='content']").BeforeSelection(element.DOM.Children().Find("p")).Text()+"\r\n\r\n\r\n" //过滤掉内容中的p标签包裹的内容
			file.Write([]byte(title))
			file.Write([]byte(content))
			//获取下一个url
			href, found:=element.DOM.Find("div[class='bottem2']").Find("a").Eq(3).Attr("href")
			if found{
				element.Request.Visit(element.Request.AbsoluteURL(href)) //然后请求下一个url,会自动触发c.OnHTML()
			}
		}
	})
	c.Visit(url)
}



//下载bilibili视频
func getVideoList(){
	//baseurl := "https://www.bilibili.com/video/"
	videoId := "BV1ME411Y71o"
	//url :=baseurl + videoId
	//获取视频列表
	mapList :=make(map[string]interface{})
	c :=colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
	)

	//因为是动态渲染的,所以需要去获取真正的数据返回
	apiUrl :="https://api.bilibili.com/x/player/pagelist?bvid="+videoId+"&jsonp=jsonp"

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("refer","https://www.bilibili.com/video/BV1ME411Y71o")
		r.Headers.Set("accept-encoding","identity")
		r.Headers.Set("origin","https://www.bilibili.com")
		r.Headers.Set("range","bytes=1040-4563")
	})

	c.OnResponse(func(r *colly.Response) {
		//r.Headers.Set("content-Type","application/json; charset=utf-8")
		var item  PageResult
		err := json.Unmarshal(r.Body, &item)
		if err !=nil {
			fmt.Println("构建结构体失败:",err.Error())
			return
		}
		for _,v := range item.Data {
			mapList[v.Part] = v.Cid
		}
		fmt.Println("构建结构体成功:",mapList)
	})




	//c.OnHTML("div[id='multi_page']", func(element *colly.HTMLElement) {
	//	fmt.Println(element.DOM.Html())
	//
	//	//element.ForEach("li", func(i int, item *colly.HTMLElement) {
	//	//	fmt.Println(item.DOM.Html())
	//	//	//获取名称
	//	//	name,_ := item.DOM.Find("a").Attr("href")
	//	//	//获取连接
	//	//	url := item.DOM.Find("a").Text()
	//	//
	//	//	fmt.Printf("%s:%s",name,url)
	//	//
	//	//})
	//})



	c.Visit(apiUrl)

}

func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result

}

type PageResult struct {
	Code int
	Message string
	ttl int
	Data []VideoInfo
}

type VideoInfo struct {
	Cid int
	Page int
	From string
	Part string
	Duration int
	Vid string
	Weblink string
	Dimension interface{}
}
