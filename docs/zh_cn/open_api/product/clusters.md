# 集群
## 1 创建集群
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
|端点 |	/products/{product_name}/clusters | |
|动作 |	POST  | |
|含义 |	创建产品的集群 | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |


#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| name| string |  集群名 | Y | 集群名必须全局唯一 | 
| description| string |  集群描述信息| N |  | 
| basic| object |  基本参数| Y |  | 
| basic.connection| object |  连接管理| Y | 内容见 [表：连接设置](#connection) | 
| basic.retries| object |  重试次数| Y | 内容见 [表：重试设置](#retries) | 
| basic.buffers| object |  缓冲设置| Y |  | 
| basic.buffers.req_write_buffer_size| string |  接受请求的缓冲字节数| Y |  | 
| basic.timeouts| object |  超时设置| Y |  内容见 [表：超时设置](#timeouts) | 
| sticky_sessions| object |  会话保持| Y | 内容见 [表：会话保持](#sticky_sessions)| 
| sub_clusters| []string |  集群中挂载的子集群| Y |  | 
| scheduler| object |  内网流量配置| Y | 具体说明见 [调度说明](traffic.md#scheduler_explain)  | 
| passive_health_check| object |  被动健康检查| Y | 具体字段见 [表：被动健康检查](#passive_health_check) | 

<a id="connection">表：连接设置</a>

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| max_idle_conn_per_rs| int | 连接池| Y | 每个BFE实例，为集群中每个RS维持的空闲长连接数。一般情况下，无需特别维持，设置为0 。<br/>设置为非0时，可以提升转发性能 | 
| cancel_on_client_close| string |  连接是否级联关闭 | Y | 设置为true时，当客户端关闭连接后，BFE同时关闭对应RS的连接 <br/>设置为false时，当客户端关闭连接后，BFE按默认策略关闭对应RS的连接 | 

<a id="retries">表： 重试设置</a>

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| max_retry_in_subcluster| string |  同一个子集群内重试次数| Y |  | 
| max_retry_cross_subcluster| string |  跨子集群重试次数| Y | - | 

<a id="sticky_sessions">表：会话保持</a>

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| session_sticky_type| string |  会话保持的粒度 | Y | INSTANCE，实例级会话保持 <br/>	SUB_CLUSTER，子集群级别会话保持| 
| hash_strategy| string |  会话保持策略  | N | CLIENT_IP_ONLY，根据client ip做会话保持 <br/>	CLIENT_ID_ONLY，根据请求中header做会话保持(默认值) <br>	CLIENT_ID_PERFERED，优先基于特定header，如果请求中没有对应header，则使用client ip| 
| hash_header| string |  指定CLIENT_ID使用的header | N |	当使用cookie作为会话保持的哈希key时，数据格式为Cookie:${key} | 


<a id="timeouts">表：超时设置</a>

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| timeout_conn_serv| string |  连接后端超时(ms)| Y |  | 
| timeout_response_header| string |  读后端响应头部超时(ms)| Y |  | 
| timeout_readbody_client| string |  读请求body超时(ms)| Y |  | 
| timeout_read_client_again| string |  与用户的长连接超时(ms) | Y |  | 
| timeout_write_client| string |  写响应超时(ms)| Y | - | 

<a id="passive_health_check">表: 被动健康检查</a>

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| failnum| int |  进入健康检查的失败次数阈值 | Y | 连续转发失败多次后，BFE进入健康检查状态，对下游RS发起探活 | 
| interval| int |  连续健康检查的时间间隔 | Y | 单位ms | 
| host| string |  健康检查请求的域名| Y | 域名后的部分 | 
| uri| string |  健康检查请求的URI  | Y |  | 
| statuscode| int |  期望的健康检查返回码 | Y | 如果需要忽略返回码，此处可以填0 | 

#### HTTP BODY中参数示例
```
{
    "name": "news_static",
    "description": "新闻静态页面集群",
    "basic": {
        "connection": {
            "max_idle_conn_per_rs": 0,
            "cancel_on_client_close": false
        },
        "retries": {
            "max_retry_in_subcluster": 2,
            "max_retry_cross_subcluster": 0
        },
        "buffers": {
            "req_write_buffer_size": 512
        },
        "timeouts": {
            "timeout_conn_serv": 50000,
            "timeout_response_header": 50000,
            "timeout_readbody_client": 30000,
            "timeout_read_client_again": 30000,
            "timeout_write_client": 60000
        }
    },
    "sticky_sessions": {
        "session_sticky_type": "INSTANCE",
        "hash_strategy": "CLIENT_ID_ONLY",
        "hash_header": "Cookie:USERID"
    },
	"sub_clusters": [
		"sub_cluster_1",
		"sub_cluster_2"
    ],
	"passive_health_check": {
		"interval": 1000,
		"failnum": 10,
		"host": "news.bfe-networks.com",
		"uri": "/index.html",
		"statuscode": 200,
	}
}
```


### 返回数据(Data内容)
同请求参数

## 2 集群列表
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点	 | /products/{product_name}/clusters || 
| 动作	 | GET ||
| 含义	 | 产品线的所有集群列表 | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | - |

### 返回数据(Data内容)
数组，单元素同创建接口


## 3 集群详情
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点 |	/products/{product_name}/clusters/{cluster_name} ||
| method |	GET  ||
| 含义 |	产品线的单个集群详情 | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字|  Y | - |

### 返回数据(Data内容)
同创建接口

## 4 更新集群基本配置
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	更新集群基本信息 |  可编辑描述信息, Basic配置段, sticky_sessions配置段, healthcheck配置段|
| 端点 |	/products/{product_name}/clusters/{cluster_name} | |
| method |	PATCH | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字 |  Y | - |


#### Body参数

可修改字段含义同创建接口。示例如下：

```
{ 
	"name": "news_static", 
	"description": "新闻静态页面集群", 
	"ready": false,
	"basic": {
		"connection": {
	 		"max_idle_conn_per_rs": 0,
			"cancel_on_client_close": false
		},
		"retries": {
			"max_retry_in_subcluster": 2,
			"max_retry_cross_subcluster": 0
		},
		"buffers": {
				"req_write_buffer_size": 512
		},
		"timeouts": {
			"timeout_conn_serv": 50000, 
			"timeout_response_header": 50000, 
			"timeout_readbody_client": 30000, 
			"timeout_read_client_again": 30000, 
			"timeout_write_client": 60000
		},
	},
	"sticky_sessions": {
		"session_sticky_type": "INSTANCE",
		"hash_strategy": "CLIENT_ID_ONLY",
		"hash_header": "Cookie:USERID"
	},
	"passive_health_check": {
		"interval": 1000,
		"failnum": 10,
		"host": "news.bfe-networks.com",
		"uri": "/index.html",
		"statuscode": 200,
	}
}
```


### 返回数据(Data内容)
同创建接口


## 5 更新集群的子集群

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点 |	/products/{product_name}/clusters/{cluster_name}/sub-clusters ||
| method |	PATCH ||
| 含义 |	在集群上，更新挂载的子集群列表 | - |


### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字|  Y | - |


#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| name | string | 集群名字|  Y ||
| sub_clusters | []string | 子集群名字列表 | Y | - |

- 新增子集群时，集群的流量调度默认执行以下处理：
	- 将子集群的流量分流比例设置为0
- 摘除一个子集群时，必须先调用流量调度接口，将子集群上的流量切走：
	- 将子集群的流量分流比例为0
- 挂载的子集群，必须满足以下条件
	- 子集群必须没有被其他集群挂载: cluster_name字段是空串
	- 子集群的必须已经就绪: ready字段是true

#### HTTP BODY中参数示例
```
{ 
	"name": "news_static", 
	"sub_clusters": [ 
		"sub_cluster_1",
		"sub_cluster_2"
	]
}
```

### 返回数据(Data内容)
同创建接口

## 6 删除集群
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 | 	删除产品线的集群 ||
| 端点 | 	/products/{product_name}/clusters/{cluster_name} || 
| method | 	DELETE | - | 

### 输入参数
#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字|  Y | - |

### 返回数据(Data内容)
同创建接口


## 7 集群就绪状态获取
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点 | 	/products/{product_name}/clusters/{cluster_name}/ready | |
| method | 	GET  ||
| 含义 | 	获取集群是否就绪的状态(可以承接线上流量)  | 当前，集群默认是就绪的 |

### 输入参数
#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字|  Y | - |

### 返回数据(Data内容)

```
{ 
	"name": "news_static", 
	"ready": false
}
```
