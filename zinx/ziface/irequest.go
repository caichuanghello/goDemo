package ziface

//把客户端请求的连接信息,包装到了这里
type IRequest interface {
	GetConnection()  IConnection //得到当前连接
	GetData() []byte
}
