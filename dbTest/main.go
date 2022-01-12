package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
var DB *sql.DB

func init(){
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	DB,_=sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/gva?charset=utf8mb4")
	//if err !=nil {
	//	fmt.Println("open error:",err.Error())
	//	return
	//}
	if err:=DB.Ping(); err !=nil{
		fmt.Println("数据库Ping连接失败:",err.Error())
		return
	}
}

func main()  {

	var nickname []byte
	var gender int

	err:=DB.QueryRow(`select nickname,gender from fk_users where uid = ?`,1).Scan(&nickname,&gender)
	fmt.Println(nickname)
	fmt.Println(gender)
	if err !=nil {
		fmt.Println("query error:",err.Error())
		return
	}




}
