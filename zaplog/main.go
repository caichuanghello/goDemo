package main

import (
	"encoding/json"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger (logpath string,loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logpath,  //日志文件路径
		MaxSize:    128, //最大字节
		MaxAge:     30,
		MaxBackups: 7,
		Compress:   true,
	}
	w := zapcore.AddSync(&hook)
	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level= zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)

	logger := zap.New(core)
	logger.Info("DefaultLogger init success")
	return logger
}


type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func main(){
	t:= &Person{
		Name: "zhangsan",
		Age:  21,
	}
	data,err := json.Marshal(t)
	if err!=nil {
		fmt.Println("marshal is failed,err:",err)
	}
	// 历史记录日志名字为：all-2018-11-15T07-45-51.763.log，服务重新启动，日志会追加，不会删除
	logger := initLogger("./log/all.log","debug")
	logger.Info(fmt.Sprint("test log"),zap.Int("line",47),zap.String("k1","v1"))
	logger.Debug(fmt.Sprint("debug log"),zap.ByteString("level",data))
	logger.Info(fmt.Sprint("Info log"),zap.String("level",`{'a':4,'b':5}`))
	logger.Warn(fmt.Sprint("Info log"), zap.String("level", `{'a':7,'b':8}`))

}