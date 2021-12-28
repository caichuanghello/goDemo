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

var syn sync.Mutex
const (
	SystemMessage = iota //系统消息

	BroadcastMessage //广播消息(正常的消息)

	HeartBeatMessage //心跳消息

	ConnectedMessage 	// 上线通知

	DisconnectedMessage //下线通知
)

//
func init(){
	//监听消息推送
	go func() {
		for {
			select {
			case client:=<-SendMsgChannel:

				if client.Msg.MsgType == ConnectedMessage{
<<<<<<< HEAD

					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
						client.Msg.Message ="进入聊天室"
						re,_:=json.Marshal(client)
						v.Conn.WriteMessage(1,re)
						return true
					})

				} else if client.Msg.MsgType == BroadcastMessage{
					re,_:=json.Marshal(client)
					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
						if v !=client {
=======
					onLineMapp.Range(func(key, value interface{}) bool {
						v :=value.(*Client)
						v.Msg.Message = "进入聊天室"
						re,_:=json.Marshal(v)
						v.Conn.WriteMessage(1,re)
						return true
					})
				} else if client.Msg.MsgType == BroadcastMessage{

					onLineMapp.Range(func(key, value interface{}) bool {
						v :=value.(*Client)
						re,_:=json.Marshal(client)
						if v != client {
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
							v.Conn.WriteMessage(1,re)
						}
						return true
					})

<<<<<<< HEAD
				} else if client.Msg.MsgType == DisconnectedMessage{
					re,_:=json.Marshal(client)
					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
						if v !=client {
							v.Conn.WriteMessage(1,re)
						}
						return true
					})
					onLineMapSafe.Delete(client.Addr)

				} else if client.Msg.MsgType == HeartBeatMessage{

					onLineMapSafe.Range(func(k, value interface{}) bool {
						v := value.(*Client)
=======
					//for _,v:=range onLineMap{
					//	re,_:=json.Marshal(client)
					//	if v != client {
					//		v.Conn.WriteMessage(1,re)
					//	}
					//}
				} else if client.Msg.MsgType == DisconnectedMessage{
					re,_:=json.Marshal(client)
					onLineMapp.Range(func(key, value interface{}) bool {
						v :=value.(*Client)
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
						if v != client {
							v.Msg.Message=client.Name+"用户退出!"
							re,_:=json.Marshal(v)
							v.Conn.WriteMessage(1,re)
						} else {
							re,_:=json.Marshal(client)
							v.Conn.WriteMessage(1,re)
						}
						return true
					})
<<<<<<< HEAD
					onLineMapSafe.Delete(client.Addr)
=======
					//关闭连接
					client.Conn.Close()
					onLineMapp.Delete(client.Addr)
					fmt.Println("用户退出",onLineMapp)
					//for _,v:=range onLineMap{
					//	if v != client {
					//		v.Conn.WriteMessage(1,re)
					//	}
					//}
					////关闭连接
					//client.Conn.Close()
					//syn.Lock()
					////删除map中的数据
					//delete(onLineMap,client.Addr)
					//syn.Unlock()
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
				}

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
var onLineMapp sync.Map
var onLineMap = make(map[string]*Client)

//并发安全的map
var onLineMapSafe sync.Map

//定义全局的消息推送channel
var SendMsgChannel = make(chan *Client)

//维持心跳
var isOnline = make(chan bool)

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
<<<<<<< HEAD
=======

	//新用户信息构建
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
	user := &Client{
		Name: name,
		Addr:remoteAddr,
		Conn: ws,
		Header: "https://tse1-mm.cn.bing.net/th/id/R-C.efdb268bb841fb60073dbae826bf2b9f?rik=Ufo6V0eAyp3IkQ&riu=http%3a%2f%2fscimg.jianbihuadq.com%2f202009%2f202009162308095.jpg&ehk=thgEdzkXNa5AqjDy3cJ5aAHwMPSGcbOS7CKvuxvNo3w%3d&risl=&pid=ImgRaw&r=0&sres=1&sresct=1",
		Addtime: time.Now().Format("15:04"),
<<<<<<< HEAD
		Msg:message{MsgType: ConnectedMessage},
	}
	//添加到map中
	onLineMapSafe.Store(remoteAddr,user)
=======
		Msg: message{MsgType: ConnectedMessage},
	}

	//onLineMap[remoteAddr] = user
	//新增
	onLineMapp.Store(remoteAddr,user)
	fmt.Println("所有用户",onLineMapp)
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
	//通知其他用户上线消息
	SendMsgChannel<-user

	for {
<<<<<<< HEAD
		_,p,err:=ws.ReadMessage()
		value,_:=onLineMapSafe.Load(remoteAddr)
		v :=value.(*Client)
		//用户退出
		if err != nil {
			fmt.Println("用户退出!")
			v.Msg = message{MsgType: DisconnectedMessage}
			SendMsgChannel<-v
			return
		}
		//群发用户发送的信息
		v.Msg =message{MsgType: BroadcastMessage,Message:string(p)}
		SendMsgChannel<-v
=======

		_,p,err:=ws.ReadMessage() //这里是阻塞的

		val,_:=onLineMapp.Load(remoteAddr)
		//用户退出
		if err != nil {
			val.(*Client).Msg = message{MsgType: DisconnectedMessage}
			SendMsgChannel<-val.(*Client)
			return
		}
		//群发用户发送的信息
		val.(*Client).Msg = message{MsgType: BroadcastMessage,Message:string(p)}
		SendMsgChannel<-val.(*Client)
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
		if err != nil {
			fmt.Println("socket写入失败:",err)
			break
		}

	}


	//for{
	//	select {
	//		case <-isOnline:
	//	case <-time.After(10*time.Second):
	//		//用户离线,
	//		user.Msg = message{MsgType: DisconnectedMessage,Message: "已断开连接.."}
	//		SendMsgChannel<-user
	//		return
	//	}
	//}

}


//获取当前已在线用户
func  (UserController)GetOnlineUser(c *gin.Context){
<<<<<<< HEAD
	c.Header("Access-Control-Allow-Origin", "*")
	var onLineSlice []*Client
	onLineMapSafe.Range(func(k, value interface{}) bool {
		v := value.(*Client)
		onLineSlice = append(onLineSlice,v)
		return true
	})
=======
	var onLineSlice []*Client
	onLineMapp.Range(func(key, value interface{}) bool {
		onLineSlice = append(onLineSlice,value.(*Client))
		return true
	})

>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
	if len(onLineSlice) ==0 {
		c.JSON(200, []struct {
		}{})

	} else {
		c.JSON(200,onLineSlice)

	}
	fmt.Println("在线用户",onLineMapp)
	return
}



