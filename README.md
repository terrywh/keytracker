### 基本功能

1. 状态监控
	由各个应用节点、服务连接并上报状态、PING，以确认应用、节点、服务的状态；
2. 配置推送
	向各应用、节点推送设置更新数据；
3. 节点推送
	节点变化通知同命名空间的其他节点；
4. 数据接口
	命名空间、节点数据接口；


### 交互数据示例

节点 => 服务器：

1. 握手登记：

	{"action":"hello","ns":"aaa","key":"yyyyy","ping":5000,"stats":{"key1":1,"key2":100,"key3":"aaaaaaa","key4":true}}
<!--
	{"action":"hello","ns":"aaa","key":"xxxxx","ping":6000,"stats":{"key1":2,"key2":200,"key3":"bbbbbbb","key4":false}}
	{"action":"hello","ns":"bbb","key":"zzzzz","ping":7000,"stats":{"key1":3,"key2":300,"key3":"ccccccc","key4":true}}
-->

	* 命名空间 ping - aaa 标识 key - yyyy 唯一指示一个节点信息；
	* 交互超时 ping - 5000 单位 秒，两次交互中超过该时间服务端将主动断开；
	* 数据初始化 stats 请参见 3.数据上报

2. 探活：

	{"action":"ping"}

3. 数据上报：

	{"action":"data","stats":{"key2":100,"key3":"fffffff"},"delta":{"key1":-1}}

	* 设置 key2 = 100
	* 设置 key3 = "fffffff"
	* 设置 key4 = key4 -1

服务器 => 节点：

1. 探活回复：

	{"action":"pong"}

2. 状态推送：

	{"action":"data","data":{"key3":"ddddddd"}}

	* 设置数据 key3 变更为 "ddddddd"


3. 节点增减、状态：

	{"action":"node","node":{"xxxxx":false,"yyyyy":true}}

	* 节点 xxxxx 下线
	* 节点 yyyyy 上线
