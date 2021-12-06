# 内网调度

## 1 调度配置获取

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	获取某集群的调度配置 | | 
| 端点 |	/products/{product_name}/clusters/{cluster_name}/scheduler | |
| method |	GET | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字|  Y | - |

### 返回数据(Data内容)
| 参数名 | 类型 |参数含义 | 补充描述 |
| - | -  | - |  - | 
| cluster | string |  集群名 |  | 
| scheduler | object | 调度的配置 |  详细说明见[说明2](#scheduler_explain) |

#### <a id="scheduler_explain">说明</a>

scheduler为流量调度的配置矩阵:
- 每个BFE集群是一个配置项
	- 例如：bfe-cluster1.sk
- 每个BFE集群分流到每个子集群的流量，通过0-100之间的数字指定
	- 例如："sub_cluster_1": 100 即100%分流到sub_cluster_1
	- 例如："sub_cluster_1": 40 和"sub_cluster_2": 60，表示分别分流40%和60%到sub_cluster_1和sub_cluster_2
- 每个BFE集群都可以强制分流到黑洞，即GSLB_BLACKHOLE分流项
- 每个BFE集群,分流目标的总比例之和必须是100


#### 数据示例

```
{ 
	"cluster": "cluster-demo",
	"scheduler": { 
		"bfe-cluster1.sk": { 
			"sub_cluster_1": 100, 
			"sub_cluster_2": 0, 
			"GSLB_BLACKHOLE": 0 
		}, 
		"bfe-cluster2.xl": { 
			"sub_cluster_1": 40, 
			"sub_cluster_2": 60, 
			"GSLB_BLACKHOLE": 0 
		}
	}
} 
```

## 2 设置调度参数
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	设置产品线的调度参数 || 
| 端点 |	/products/{product_name}/clusters/{cluster_name}/scheduler ||
| method |	PATCH | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | |
| cluster_name | string | 集群名字|  Y | - |

#### Body参数
含义见说明
请求参数示例：
```
{ 
	"bfe-cluster1.sk": { 
		"sub_cluster_1": 100, 
		"sub_cluster_2": 0, 
		"GSLB_BLACKHOLE": 0 
	}, 
	"bfe-cluster2.xl": { 
		"sub_cluster_1": 0, 
		"sub_cluster_2": 100, 
		"GSLB_BLACKHOLE": 0 
	} 
}
```

### 返回数据(Data内容)
同获取接口
