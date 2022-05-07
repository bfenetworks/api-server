# 子集群
## 1 创建子集群
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点	| /products/{product_name}/sub-clusters ||
| 动作	| POST  ||
| 含义	| 创建子集群 | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名字 | Y | - |


#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| name | string |   子集群名字 | Y | 一个产品线内子集群名字唯一 |
| instance_pool | string |  子集群关联的实例池 | Y | 必须是实例池的完整名字：{product_name}.{instance_pool_name} | 
| description | string |   子集群描述信息 | N | - |

#### 输入参数示例
```
{ 
	"name": "subcluster_demo", 
	"instance_pool": "pool_demo", 
	"description": "description message"
}
```

### 返回数据(Data内容)
同输入参数

#### 成功返回数据示例
```
{ 
	"name": "subcluster_demo", 
	"instance_pool": "pool_demo", 
	"ready": true, 
	"description": "description message", 
	"cluster_name":"cluster_demo"
}
```

- ready: 表示子集群是否就绪。只有已就绪的子集群可以被挂载到集群，准备接入流量。
- cluster_name: 在集群配置时，会选择子集群进行挂载
	- 如果当前子集群被挂载到某个集群上，那么cluster_name就是对应的集群名字
	- 如果当前子集群未挂载到某个集群上，那么cluster_name就是空串



## 2 子集群详情
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点	| /products/{product_name}/sub-clusters/{sub_cluster_name}  ||
| 动作	| GET  ||
| 含义	| 获取子集群详情 | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名字 | Y | |
| sub_cluster_name | string | 子集群名字 | Y | - |

### 返回数据(Data内容)
同创建接口

#### 成功返回数据示例
同创建的成功返回数据示例


## 3 子集群列表
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	获取产品线的子集群列表 | |
| 端点 |	/products/{product_name}/sub-clusters  | |
| method |	GET  | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名字 | Y | - |

### 返回数据(Data内容)
同创建接口

#### 成功返回数据示例
```
[ 
	{ 
		"name": "subcluster_demo", 
		"instance_pool": "pool_demo", 
		"ready": true, 
		"description": "description message", 
		"cluster_name":"cluster_demo"
	}
] 
```

## 4 更新子集群
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点 |	/products/{product_name}/sub-clusters/{sub_cluster_name} ||
| method |	PATCH ||
| 含义 |	更新子集群配置 | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名字 | Y | |
| sub_cluster_name | string | 子集群名字 | Y | - |

#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| name | string |   子集群名字 | N | 一个产品线内子集群名字唯一 |
| description | string |   子集群描述信息 | N | - |

#### 输入参数示例
```
{ 
	"name": "subcluster_demo", 
	"description": "description message"
}
```

### 返回数据(Data内容)
同创建接口



## 5 删除子集群
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	 |删除子集群 ||
| 端点	 |/products/{product_name}/sub-clusters/{sub_cluster_name}||
| 动作	 |DELETE| -|

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名字 | Y | |
| sub_cluster_name | string | 子集群名字 | Y | - |

### 返回数据(Data内容)
同创建接口


## 6 子集群列表（系统管理员模式）
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	获取所有的子集群列表 || 
| 端点 |	/sub-clusters ||
| method |	GET | - |

### 输入参数	
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| Query参数 | | | 
| pool_name | string | 实例池名字 | N | - |

### 返回数据(Data内容)
同创建接口

#### 成功返回数据示例
```
[ 
	{ 
		"name": "subcluster_demo", 
		"product_name": "baike", 
		"instance_pool": "pool_demo", 
		"ready": true, 
		"description": "description message", 
		"cluster_name":"cluster_demo"
	}
] 
```
