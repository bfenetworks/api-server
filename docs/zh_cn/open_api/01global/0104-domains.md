# 域名
## 1 添加域名

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义|	为产品线添加域名||
| 端点 |	/products/{product_name}/domains ||
| method |	POST| -|

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |

#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| name | string | 域名 | Y | 添加的域名 |

#### 请求参数示例

```
{
	"name": "static.bfe-networks.com"
}
```

### 返回数据(Data内容)
同请求参数

## 2 域名列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 | 	获取产品线的域名列表 ||
| 端点 | 	/products/{product_name}/domains ||
| method | 	GET | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |

### 返回数据(Data内容)
```
[ 
	"www.bfe-networks.com",
    "static.bfe-networks.com" 
]
```

## 3 删除域名

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	删除产品线的域名 ||
| 端点 |	/products/{product_name}/domains/{domain_name} ||
| method |	DELETE | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | |
|	domain_name | string | 待删除的域名 | Y | - |

### 返回数据(Data内容)
同创建接口

## 4 域名使用情况

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义|	查看产品线的域名 ||
| 端点|	/products/{product_name}/domains/{domain_name}/use-status ||
| 动作|	GET | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | |
| domain_name | string | 待删除的域名 | Y | - |

### 返回数据(Data内容)
| 参数名 | 类型 |参数含义 |  补充描述 |
| - | -  | - | - | 
| be_used | bool | 是否被使用 | |
| dep_type | string | 被依赖对象的类型 | 只有被使用时才返回，当前只有一种类型 ConditionExpression, 表示域名被某个条件表达式引用|
| dep_name | string | 被依赖对象的名字 | 使用当前域名的条件表达式名字 | 

说明：如果有多个被依赖的对象，只返回第一个。

#### 成功返回数据示例
```
{
    "be_used": true,
    "dep_type": "ConditionExpression",
    "dep_name": "rule_name_abc"
}
```