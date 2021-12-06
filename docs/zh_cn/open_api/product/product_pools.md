# 产品线实例池

## 1 创建实例池
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 创建产品线的实例池 ||
| 端点	| /products/{product_name}/instance-pools ||
| 动作	| POST  | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |

#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| name| string | 实例池的完整名字 | Y | 格式为 {product_name}.{instance_pool_name}<br/>产品线名字必须与URL中product_name相同，<br/>点"."是一个特殊字符，仅能用于分隔产品线名与实例池名 |
| instances| [] |  实例列表 | Y | |
| instances[].hostname| string | 实例所在的主机名 | Y |在没有主机名时，可以填写主机的IP地址|
| instances[].ip| string |  实例的IP地址 | Y | |
| instances[].weight| int | 实例的权重，数字范围[0,100] | Y | |
| instances[].ports| string | 实例上的端口 | Y |  每个端口有一个名字 <br> 每个实例至少有一个默认端口，名字是Default |
| instances[].tags| string | 实例上的标签 | N | 每个标签都是一个key/value对，value必须是字符串 |

#### HTTP BODY中参数示例
```
{
    "name": "product1.instance_pool2", 
    "instances": [ 
        { 
            "hostname": "hostname1", 
            "ip": "10.70.29.3", 
            "weight": 1, 
            "ports": {
                "Default": 80
            },
            "tags": {
                "tag1": "val1"
            }
        }
    ]
}
```

### 返回数据(Data内容)
同创建接口

#### 成功返回数据示例

```
{
    "name": "product1.instance_pool2", 
    "instances": [ 
        { 
            "hostname": "hostname1", 
            "ip": "10.70.29.3", 
            "weight": 1, 
            "ports": {
                "Default": 80
            },
            "tags": {
                "tag1": "val1"
            }
        }
    ]
}
```


## 2 实例池列表
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取产品线的实例池列表  | |
| 端点	| /products/{product_name}/instance-pools | |
| 动作	| GET  | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | - |

### 返回数据(Data内容)

为一个字符串数组，每个元素为实例池名。

```
[ 
    "product1.instance_pool1",
    "product1.instance_pool2" 
]
```


## 3 实例池详情
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 | 	获取产品线的实例池的详情 ||
| 端点 | 	/products/{product_name}/instance-pools/{instance_pool_name} ||
| method | 	GET | - | 

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | |
|	instance_pool_name | string | 实例池名字 | Y | - |

### 返回数据(Data内容)

同创建接口

## 4 更新实例池
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	更新产品线的实例池 | 该更新是全量更新，不支持仅添加部分数据 |
| 端点 |	/products/{product_name}/instance-pools/{instance_pool_name} ||
| method |	PATCH | - | 

### 输入参数
同创建接口


### 返回数据(Data内容)
同创建接口


## 5 删除实例池
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 删除产品线的实例池 ||
| 端点	| /products/{product_name}/instance-pools/{instance_pool_name} ||
| 动作	| DELETE | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | |
|	instance_pool_name | string | 实例池名字 | Y | 如果实例池被使用，将删除失败 |

### 返回数据(Data内容)

同创建接口