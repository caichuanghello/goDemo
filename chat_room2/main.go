package main

import (
	"chat_room2/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main(){

	r:=gin.Default()
	r.LoadHTMLGlob("view/*")
	r.StaticFS("/static", http.Dir("./static"))
	router.ApiRoutersInit(r)
	
	r.Run(":8080")
}
