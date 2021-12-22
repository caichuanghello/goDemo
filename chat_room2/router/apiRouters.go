package router

import (
	"chat_room2/controller/api"
	"chat_room2/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func ApiRoutersInit(r *gin.Engine) {

	//需要验证token的
	apiRouters :=r.Group("/api",middleware.CheckLogin)
	{
		apiRouters.GET("/getUserInfo",func(c *gin.Context){
			fmt.Println("success")
		})

	}

	//不需要验证token的
	apiRouters =r.Group("/")
	{
		//登入页面
		apiRouters.GET("/login", func(c *gin.Context) {
			c.HTML(200,"login.html", struct {}{})
			return
		})
		apiRouters.POST("/login",api.UserController{}.Login)


		//首页
		apiRouters.GET("/index", func(c *gin.Context) {
			c.HTML(200,"index.html", struct {}{})
			return
		})

		//websocket
		apiRouters.GET("/ws", api.UserController{}.Ws)

		//获取在线用户数据
		apiRouters.GET("/getOnlineUser", api.UserController{}.GetOnlineUser)

	}




}



