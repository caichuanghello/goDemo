package api

import "github.com/gin-gonic/gin"

type baseController struct {

}


type ApiReturn struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func (baseController)Success(c *gin.Context) {
	c.String(200,"success")
}

func (baseController)Error(c *gin.Context) {
	c.String(200,"error")
}

/*
正确返回数据
*/
func ApiSuccess(data interface{}) ApiReturn{
	var res ApiReturn
	res.Code = "000000"
	res.Message = "SUCCESS"
	res.Data = data
	return res
}

/*
错误返回数据
*/
func ApiError(err interface{}) ApiReturn{
	var res ApiReturn
	if err,ok :=err.(error);ok{
		res.Message = err.Error()
	}
	if msg,ok :=err.(string);ok {
		res.Message = msg
	}
	res.Code = "000001"
	res.Data = struct{}{}
	return res
}



