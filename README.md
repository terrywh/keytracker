### 基本功能

1. 状态监控
	由各个应用节点、服务连接并上报状态、PING，以确认应用、节点、服务的状态；
2. 配置推送
	向各应用、节点推送设置更新的的状态数据；


### 交互数据示例

节点 => 服务器：

1. 握手登记：

	{"action":"hello","ns":"abc","key":"yyyyy","ping":6000,"stats":{"key1":50,"key2":200,"key3":"abcdefg","key4":true}}
	{"action":"hello","ns":"abc","key":"xxxxx","ping":6000,"stats":{"key1":50,"key2":200,"key3":"abcdefg","key4":true}}
	{"action":"hello","ns":"abc","key":"zzzzz","ping":6000,"stats":{"key1":50,"key2":200,"key3":"abcdefg","key4":true}}


2. 探活：

	{"action":"ping"}

3. 数据上报：

	{"action":"data","stats":{"key1":100,"key2":200,"key3":"xxxxxx"},"delta":{"key1":-2, "key2":1}}

服务器 => 节点：

1. 探活回复：

	{"action":"pong"}

2. 状态推送：

	{"action":"data","data":{"key3":"yyyyyy"}}

3. 节点增减、状态：

	{"action":"node","node":{"xxxxx":false, "yyyyy":true}}
