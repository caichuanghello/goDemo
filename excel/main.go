package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

//定义的是excel的列名,不够自己可以按照xls规则加上去
var row  = []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
func main(){
	f := excelize.NewFile()
	//模拟从数据库中获取数据
	result :=getData()
	//利用反射导出数据
	err:=leadOut1(result,f)
	//常规方式
	//err:=leadOut2(result,f)
	if err != nil{
		fmt.Println("excel导出失败:",err)
		return
	}
	fmt.Println("excel数据导出成功")
}




//模拟从数据库中获取数据
func getData() (result []Record){
	result =[]Record{
		{Username: "蔡闯",ApplyDate: "2021-11-23",Mobile: "13588318049",Gender: "男",Company: "杭州睿者",CreateTime: "2021-11-15"},
		{Username: "小李",ApplyDate: "2021-11-25",Mobile: "16756213256",Gender: "女",Company: "杭州知了",CreateTime: "2021-11-24"},
	}
	return
}


type Record struct {
	Username string `row_name:"用户名"`
	ApplyDate string `row_name:"预约时间"`
	Mobile string `row_name:"手机号"`
	Gender string `row_name:"性别"`
	Company string `row_name:"公司名称"`
	CreateTime string  `row_name:"创建时间"`
}

//方法一,使用反射导出数据
func leadOut1(result []Record,f *excelize.File)(err error){
	for i:=0;i<len(result)+1;i++ {
		var value reflect.Value
		if i==0{
			//获取表头设置的参数
			value=reflect.ValueOf(result[i])
		} else{
			value=reflect.ValueOf(result[i-1])
		}
		file_num := value.NumField()
		for j:=0;j<file_num;j++ {
			axis := row[j]+strconv.Itoa(i+1)
			var val string
			if i==0 {
				//设置第一行表头
				val = value.Type().Field(j).Tag.Get("row_name") //第一行第一列的值
				if len(val) == 0 { //如果没有设置,则使用字段名
					val = value.Type().Field(j).Name
				}
			} else{
				//数据插入
				val = value.Field(j).String()
			}
			f.SetCellValue("Sheet1", axis, val)
		}
	}
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("excel/download/Book1.xlsx"); err != nil {
		return err
	}

	return nil
}


//方法二,常规导出数据
func leadOut2(result []Record,f *excelize.File)(err error){

	//设置表头
	f.SetCellValue("Sheet1", "A1", "用户名")
	f.SetCellValue("Sheet1", "B1", "预约时间")
	f.SetCellValue("Sheet1", "C1", "手机号")
	f.SetCellValue("Sheet1", "D1", "性别")
	f.SetCellValue("Sheet1", "E1", "公司名称")
	f.SetCellValue("Sheet1", "F1", "生成时间")

	//循环插入数据
	for i:=0;i<len(result);i++ {
		Username := result[i].Username
		ApplyDate := result[i].ApplyDate
		Mobile := result[i].Mobile
		Gender := result[i].Gender
		Company := result[i].Company
		CreateTime := result[i].CreateTime
		//列名称
		ind :=strconv.Itoa(i+2)
		f.SetCellValue("Sheet1", "A"+ind, Username)
		f.SetCellValue("Sheet1", "B"+ind, ApplyDate)
		f.SetCellValue("Sheet1", "C"+ind, Mobile)
		f.SetCellValue("Sheet1", "D"+ind, Gender)
		f.SetCellValue("Sheet1", "E"+ind, Company)
		f.SetCellValue("Sheet1", "F"+ind, CreateTime)
	}
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("excel/download/Record.xlsx"); err != nil {
		return err
	}
	return nil
}




