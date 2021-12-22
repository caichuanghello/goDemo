package middleware

import (
	"chat_room2/controller/api"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var Claims struct {
	authorities interface{}
	created int
	exp int
	sub string
}

/*
验证token,判断是否正确
*/
func CheckLogin(c *gin.Context) {
	val:=c.GetHeader("token")
	if val =="" {
		c.Abort()
		c.JSON(200,api.ApiSuccess("token不存在!"))
	}
	mySignKey := "abcdefgh"
	mySignKeys,_ := base64.URLEncoding.DecodeString(mySignKey)
	parseAtth, _ := jwt.Parse(val, func(*jwt.Token) (interface{}, error) {
		return mySignKeys, nil
	})
	claims :=parseAtth.Claims.(jwt.MapClaims)
	tim :=time.Now().Unix()
	if val,ok := claims["exp"]; ok && int64(val.(float64)) < tim { //如果验证不通过,则直接返回
		c.Abort()
		c.JSON(200,api.ApiSuccess("token过期"))
	}
	if val,ok := claims["exp"]; ok && int64(val.(float64)) > tim {
		c.Set("uid",claims["uid"])
		c.Set("role",claims["role"])
	}
	c.Next()

}