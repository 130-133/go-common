package errorx

// https://mini1.feishu.cn/docx/doxcn6SZIPJhNqovp2GAtrFktkg

type CategoryCode int

const (
	SystemError  CategoryCode = iota + 10 //系统错误
	ParamError                            //参数错误
	GetDataError                          //获取数据错误
	CacheError                            //缓存错误
	DbError                               //数据库错误
	MqError                               //RMQ错误
	HttpError                             //HTTP请求错误
	RpcError                              //Rpc请求错误
)

type AppCode int

//应用标识(2位数字)
const (
	Common AppCode = 10
	Admin  AppCode = 11
)
