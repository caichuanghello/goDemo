package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func GetCityMap(t string) (cityMap map[string]string){

	cityMap = make(map[string]string)

	filename := "utils/city_"+t+".json"

	file,err:=os.OpenFile(filename,os.O_RDONLY,0666)
	if err !=nil{
		fmt.Println("openfile err:",err.Error())
		return
	}
	res,err:=ioutil.ReadAll(file)

	err=json.Unmarshal(res,&cityMap)
	if err !=nil{
		fmt.Println("Unmarshal error :",err.Error())
		return
	}

	return

}


func NewCityMap(){
	url:="https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9226"
	resp,err:= http.Get(url)
	if err !=nil{
		fmt.Println("获取cityMap失败:",err.Error())
		return
	}
	defer resp.Body.Close()

	res,_:=ioutil.ReadAll(resp.Body)
	ress := strings.TrimPrefix(string(res),"var station_names ='")
	ress = strings.TrimSuffix(ress,"';")

	reg,err:=regexp.Compile(`@[a-z]{3}\|`)
	if err != nil{
		fmt.Println("regexp.Compile error ",err.Error())
		return
	}

	rest:=reg.Split(ress, -1)

	cityMap :=make(map[string]string)
	cityMap1 :=make(map[string]string)
	for _,v:=range rest{
		resst:=strings.Split(v,"|")
		if len(v) >0 {
			cityMap[resst[0]] = resst[1]
			cityMap1[resst[1]] = resst[0]
		}
	}
	//写入数据
	jsonData,err :=json.MarshalIndent(&cityMap,"","\t")
	jsonData1,err :=json.MarshalIndent(&cityMap1,"","\t")
	if err != nil {
		fmt.Println("json 数据失败:",err.Error())
		return
	}

	file,_:=os.OpenFile("utils/city_cn.json",os.O_WRONLY|os.O_CREATE,0666)
	file1,_:=os.OpenFile("utils/city_en.json",os.O_WRONLY|os.O_CREATE,0666)
	defer file.Close()
	defer file1.Close()
	file.Write(jsonData)
	file1.Write(jsonData1)
}