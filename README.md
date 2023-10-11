# go-common

go公共代码仓库

[toc]

# utils
## context
基于上下文提取header，meta，auth信息
- header 提取UserAgent、IP等
- meta 提取RPC meta信息
- auth 提取鉴权信息

## crontab
定时任务基础包，用于继承

## dao
Dao层基础包，用于继承

## encrypt
加密算法工具包，encrypt入口，集成多个
- aes CBC
- jwt
- md5
- sha256
- param_sign 支付宝参数签名方式
- mini_authorized 迷你世界auth验签

## errorx
生成带错误码错误类，需定义系统应用码

## grpc
grpc客户端基础包，用于连接时使用
- Connect

## help
常用方法工具
- context 生成上下文
- datetime 时间格式化或解析
- helper 常用方法
- goroutine 协程方法
- ip 获取本地网卡ip
- mongo mongo相关方法
- price 金额类字符整形转化
- string 字符串转化

## locker
加锁类
- redis实现分布式锁

## middleware
HTTP中间件
- checkauth 鉴权

## middleware_grpc
grpc中间件
- grpcLog 日志 

## mongo
mongo连接基础包
- Connect

## mysql
mysql连接基础包
- Connect

## queue
常驻任务实现的基础包

## rabbitmq
rabbitmq基础包，实现生产消费方法
- 生产
- 延迟生产
- 带阶梯时间延迟生产
- 消费

## redis
redis连接基础包，追加方法
- LoadSetEx 查询string，不存在即写入
- LoadHSetEx 查询hash单字段，不存在即写入
- LoadHMSetEx 查询hash所有字段，不存在即写入
- PushList 基于List实现队列发布
- PullList 基于List实现队列消费
- PushQueue 基于stream实现队列发布(需redis 5.0支持)
- PullQueue 基于stream实现队列消费(需redis 5.0支持)

## request
resty实现的http请求客户端
- get
- post

## response
项目HTTP响应统一方法

## thirdpay
第三方支付
- appleiap 苹果内购

## tracer
链路追踪相关方法