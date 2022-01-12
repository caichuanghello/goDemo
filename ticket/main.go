package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"ticket/utils"
	"time"
)

const (
	depCity = "杭州"
	desCity = "松阳"
	date = "2022-01-10"
	cook ="_uab_collina=163307956321853337672101; JSESSIONID=13E577830AD7E6F75EC133419958E543; tk=Af3YUSkvAAwmQGRE8h7hDwz1pGR1bsNfmkCXfed_b7E92c1c0; _jc_save_wfdc_flag=dc; _jc_save_fromStation=%u676D%u5DDE%2CHZH; _jc_save_toStation=%u677E%u9633%2CSUU; BIGipServerotn=351273482.38945.0000; guidesStatus=off; highContrastMode=defaltMode; cursorStatus=off; BIGipServerpool_passport=199492106.50215.0000; RAIL_EXPIRATION=1642029649260; RAIL_DEVICEID=gIDYb5wfqKssDjTv-ZSGjazEyWz13hPR-42oh8Dt7OTO1QJ38guaUazFRwhRkFa2kM4kgrU3NJslvQhgJoI27mS6KjAw7Bo_o4MuGJ8j65ZuZM4T0I_NwWW1ZQs7sJSx3ZeWGn7tOArr1wWKwuGK214N4EVv8DzU; route=6f50b51faa11b987e576cdb301e545c4; _jc_save_toDate=2022-01-09; current_captcha_type=Z; _jc_save_showIns=true; _jc_save_fromDate=2022-01-10; uKey=0312b8c6d7e2c15cf1944632fdaab905015c1346dfea2de3e697520bc485429d"

)

func main01(){

	cityMap := utils.GetCityMap("en")

	che,err:=getTripsInfo()
	if err != nil{
		fmt.Println("getTripsInfo error:",err.Error())
		return
	}
	res:=che.Data.Result
	for _,v:=range res{
		detailInfo:=strings.Split (v,"|")
		fmt.Println("车次:",detailInfo[3])
		fmt.Println("始发站:",cityMap[detailInfo[4]])
		fmt.Println("终点站:",cityMap[detailInfo[5]])
		fmt.Println("本次上车点:",cityMap[detailInfo[6]])
		fmt.Println("本次下车点:",cityMap[detailInfo[7]])
		fmt.Println("出发时间:",detailInfo[8])
		fmt.Println("到达时间:",detailInfo[9])
		fmt.Println("历时:",detailInfo[10])
		fmt.Println("一等座:",detailInfo[26])
		fmt.Println("二等座:",detailInfo[30])
		fmt.Println("无座:",detailInfo[31])
		fmt.Println()
	}


	//下单


}


type User struct {
	Name string `json:"name"`
	CreateTime MyTime `json:"create_time"`
}
type MyTime int64

//返回json数据的时候重写该方法,进行格式化转化
func (t MyTime) MarshalJSON() ([]byte, error) {
	fmt.Println("进来了!")
	tTime := int64(t)
	return []byte(fmt.Sprintf("\"%v\"", time.Unix(tTime,0).Format("2006-01-02 15:04:05"))), nil
}

func main(){

	var user = []User{}
	user = append(user,User{"小李",123123123})
	re,_:=json.MarshalIndent(user,""," ")
	fmt.Println(string(re))

	var user1 = User{"小王",123123123123}

	re1,_:=json.MarshalIndent(user1,""," ")
	fmt.Println(string(re1))
}



func getTripsInfo() (ceh checi,err error){



	var che checi

	cityMap := utils.GetCityMap("cn")
	//cityMap["杭州"] = "HZH"
	//cityMap["松阳"] = "SUU"

	apiUrl:="https://kyfw.12306.cn/otn/leftTicket/queryT?leftTicketDTO.train_date="+date+"&leftTicketDTO.from_station="+cityMap[depCity]+"&leftTicketDTO.to_station="+cityMap[desCity]+"&purpose_codes=ADULT"
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Set("Cookie", cook)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error :",err)
		return che ,nil
	}
	defer resp.Body.Close()

	res,err:=ioutil.ReadAll(resp.Body)
	if err !=nil {
		fmt.Println("错误",err.Error())
		return che,err
	}
	json.Unmarshal(res,&che)
	return che,nil
}




type checi struct {
	Data DataStruct
	HttpStatus int
	Message string
	Status bool
}

type DataStruct struct {
	Result []string
	Flag string
	Map map[string]string
}
