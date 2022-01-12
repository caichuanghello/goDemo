package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"sync"
	"time"
)

type UserController struct {

}
const (
	// SystemMessage 系统消息
	SystemMessage = iota
	// BroadcastMessage 广播消息(正常的消息)
	BroadcastMessage
	// HeartBeatMessage 心跳消息
	HeartBeatMessage

	ConnectedMessage 	// 上线通知
	// DisconnectedMessage 下线通知
	DisconnectedMessage
)

//
func init(){
	//监听消息推送
	go func() {
		for {
			select {
			case client:=<-SendMsgChannel:
				re,_:=json.Marshal(client)
				if client.Msg.MsgType == ConnectedMessage{

					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
						client.Msg.Message ="进入聊天室"
						//re,_:=json.Marshal(client)
						v.Conn.WriteMessage(1,re)
						return true
					})
					//for _,v:=range onLineMap{
					//	client.Msg.Message ="进入聊天室"
					//	re,_:=json.Marshal(client)
					//	v.Conn.WriteMessage(1,re)
					//}
				} else if client.Msg.MsgType == BroadcastMessage{

					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
						if v !=client {
							v.Conn.WriteMessage(1,re)
						}
						return true
					})

					//for _,v:=range onLineMap{
					//	re,_:=json.Marshal(client)
					//	if v != client {
					//		v.Conn.WriteMessage(1,re)
					//	}
					//}
				} else if client.Msg.MsgType == DisconnectedMessage{

					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
						if v !=client {
							v.Conn.WriteMessage(1,re)
						}
						return true
					})
					onLineMapSafe.Delete(client.Addr)
					//for _,v:=range onLineMap{
					//	if v != client {
					//		v.Conn.WriteMessage(1,re)
					//	}
					//}
					//关闭连接
					//client.Conn.Close()
					//var wg sync.WaitGroup
					//wg.Add(1)
					////删除map中的数据
					//
					//wg.Done()
				}
				//case <-time.After(10*time.Second):
				//	//主动断开
				//	fmt.Println("主动断开")
			}
		}
	}()
}




var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//定义全局在线用户
var onLineMap = make(map[string]*Client)

//并发安全的map

var onLineMapSafe sync.Map

//定义全局的消息推送channel
var SendMsgChannel = make(chan *Client)

//定义用户结构体
type Client struct {
	Name string
	Conn *websocket.Conn
	Header string
	Msg message
	Addr string
	Addtime string

}

func (UserController)Login(c *gin.Context){
	c.JSON(200,ApiSuccess(struct {}{}))
}

type message struct {
	MsgType int
	Message string
}
func (UserController)Ws(c *gin.Context){
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer ws.Close()
	if err != nil {
		c.Writer.Write([]byte(err.Error()))
		return
	}
	defer  ws.Close()

	//第一步,首先把该用户加入到在线列表中
	remoteAddr := ws.RemoteAddr().String()
	//name是 "用户-端口"
	strsli:=strings.Split(remoteAddr,":")
	name :="用户"+strsli[1]
	user := &Client{
		Name: name,
		Addr:remoteAddr,
		Conn: ws,
		Header: "https://tse1-mm.cn.bing.net/th/id/R-C.efdb268bb841fb60073dbae826bf2b9f?rik=Ufo6V0eAyp3IkQ&riu=http%3a%2f%2fscimg.jianbihuadq.com%2f202009%2f202009162308095.jpg&ehk=thgEdzkXNa5AqjDy3cJ5aAHwMPSGcbOS7CKvuxvNo3w%3d&risl=&pid=ImgRaw&r=0&sres=1&sresct=1",
		Addtime: time.Now().Format("15:04"),
		Msg:message{MsgType: ConnectedMessage},
	}
	//添加到map中
	onLineMapSafe.Store(remoteAddr,user)
	SendMsgChannel<-user
	//通知其他用户上线消息
	//onLineMap[remoteAddr].Msg = message{MsgType: ConnectedMessage}
	//SendMsgChannel<-onLineMap[remoteAddr]
	for {
		_,p,err:=ws.ReadMessage()

		value,_:=onLineMapSafe.Load(remoteAddr)
		v :=value.(*Client)


		//用户退出
		if err != nil {
			fmt.Println("用户退出!")

			//value,_:=onLineMapSafe.Load(remoteAddr)
			//v :=value.(*Client)
			v.Msg = message{MsgType: DisconnectedMessage}

			SendMsgChannel<-v

			//onLineMap[remoteAddr].Msg = message{MsgType: DisconnectedMessage}
			//SendMsgChannel<-onLineMap[remoteAddr]
			return
		}

		//群发用户发送的信息
		v.Msg =message{MsgType: BroadcastMessage,Message:string(p)}

		//onLineMap[remoteAddr].Msg = message{MsgType: BroadcastMessage,Message:string(p)}
		SendMsgChannel<-v
		if err != nil {
			fmt.Println("socket写入失败:",err)
			break
		}


	}

}


//获取当前已在线用户
func  (UserController)GetOnlineUser(c *gin.Context){
	c.Header("Access-Control-Allow-Origin", "*")
	var onLineSlice []*Client

	onLineMapSafe.Range(func(k, value interface{}) bool {
		v := value.(*Client)
		onLineSlice = append(onLineSlice,v)
		return true
	})

	//for _,v := range onLineMap{
	//	onLineSlice = append(onLineSlice,v)
	//}

	if len(onLineSlice) ==0 {
		c.JSON(200, []struct {
		}{})

	} else {
		c.JSON(200,onLineSlice)

	}
	return
}



